package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/photobooth/backend/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidToken      = errors.New("invalid token")
)

// UserRepository handles user data operations
type UserRepository struct {
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

// NewUserRepository creates a new user repository
func NewUserRepository(collection *gocb.Collection, cluster *gocb.Cluster) *UserRepository {
	return &UserRepository{
		collection: collection,
		cluster:    cluster,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	// Generate user ID
	userID := fmt.Sprintf("user::%s", uuid.New().String()[:8])
	user.ID = userID
	user.DocType = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.PhotoCount = 0
	user.SubscriptionTier = "free"
	user.SubscriptionStatus = "active"
	user.RefreshTokens = []models.RefreshToken{}

	fmt.Printf("Inserting user with ID: %s, Email: %s\n", userID, user.Email)

	// Insert user
	_, err := r.collection.Insert(userID, user, nil)
	if err != nil {
		fmt.Printf("Error inserting user: %v\n", err)
		if errors.Is(err, gocb.ErrDocumentExists) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("User created successfully: %s\n", userID)
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(userID string) (*models.User, error) {
	var user models.User
	result, err := r.collection.Get(userID, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := result.Content(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	// Query using N1QL from the users collection
	query := `
		SELECT META().id as id, u.* 
		FROM photobooth_data._default.users AS u
		WHERE u.docType = "user" AND u.email = $1 
	`

	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{email},
	})
	if err != nil {
		fmt.Printf("Query error for email %s: %v\n", email, err)
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	defer rows.Close()

	// Check if we have results
	if !rows.Next() {
		fmt.Printf("No user found with email: %s\n", email)
		return nil, ErrUserNotFound
	}

	// Get the row data
	var rawResult map[string]interface{}
	if err := rows.Row(&rawResult); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	// Manually map the user fields
	var user models.User
	user.ID = rawResult["id"].(string)
	user.DocType = rawResult["docType"].(string)
	user.Email = rawResult["email"].(string)
	user.Name = rawResult["name"].(string)
	user.PasswordHash = rawResult["passwordHash"].(string)
	user.SubscriptionTier = rawResult["subscriptionTier"].(string)
	user.SubscriptionStatus = rawResult["subscriptionStatus"].(string)
	user.PhotoCount = int(rawResult["photoCount"].(float64))

	createdAt, _ := time.Parse(time.RFC3339, rawResult["createdAt"].(string))
	user.CreatedAt = createdAt
	updatedAt, _ := time.Parse(time.RFC3339, rawResult["updatedAt"].(string))
	user.UpdatedAt = updatedAt

	// Parse refresh tokens if present
	if tokens, ok := rawResult["refreshTokens"].([]interface{}); ok {
		user.RefreshTokens = make([]models.RefreshToken, len(tokens))
		for i, token := range tokens {
			tokenMap := token.(map[string]interface{})
			user.RefreshTokens[i] = models.RefreshToken{
				TokenID:     tokenMap["tokenId"].(string),
				HashedToken: tokenMap["hashedToken"].(string),
			}
			if expiresAt, ok := tokenMap["expiresAt"].(string); ok {
				expiry, _ := time.Parse(time.RFC3339, expiresAt)
				user.RefreshTokens[i].ExpiresAt = expiry
			}
			if createdAt, ok := tokenMap["createdAt"].(string); ok {
				created, _ := time.Parse(time.RFC3339, createdAt)
				user.RefreshTokens[i].CreatedAt = created
			}
		}
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.collection.Upsert(user.ID, user, nil)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// AddRefreshToken adds a refresh token to a user
func (r *UserRepository) AddRefreshToken(userID string, refreshToken models.RefreshToken) error {
	user, err := r.GetByID(userID)
	if err != nil {
		return err
	}

	// Add refresh token
	user.RefreshTokens = append(user.RefreshTokens, refreshToken)

	return r.Update(user)
}

// RemoveRefreshToken removes a refresh token from a user
func (r *UserRepository) RemoveRefreshToken(userID, tokenID string) error {
	user, err := r.GetByID(userID)
	if err != nil {
		return err
	}

	// Remove refresh token
	var newTokens []models.RefreshToken
	found := false
	for _, token := range user.RefreshTokens {
		if token.TokenID != tokenID {
			newTokens = append(newTokens, token)
		} else {
			found = true
		}
	}

	if !found {
		return ErrInvalidToken
	}

	user.RefreshTokens = newTokens

	return r.Update(user)
}

// GetRefreshToken retrieves a refresh token from a user
func (r *UserRepository) GetRefreshToken(userID, tokenID string) (*models.RefreshToken, error) {
	user, err := r.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Find refresh token
	for _, token := range user.RefreshTokens {
		if token.TokenID == tokenID {
			return &token, nil
		}
	}

	return nil, ErrInvalidToken
}

// IncrementPhotoCount increments the user's photo count
func (r *UserRepository) IncrementPhotoCount(userID string) error {
	user, err := r.GetByID(userID)
	if err != nil {
		return err
	}

	user.PhotoCount++

	return r.Update(user)
}
