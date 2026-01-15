package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	DocType               string         `json:"docType"`
	ID                    string         `json:"id"`
	Email                 string         `json:"email"`
	Name                  string         `json:"name"`
	PasswordHash          string         `json:"passwordHash"`
	SubscriptionTier      string         `json:"subscriptionTier"`   // "free" | "premium"
	SubscriptionStatus    string         `json:"subscriptionStatus"` // "active" | "cancelled" | "expired"
	StripeCustomerID      string         `json:"stripeCustomerId,omitempty"`
	StripeSubscriptionID  string         `json:"stripeSubscriptionId,omitempty"`
	SubscriptionExpiresAt *time.Time     `json:"subscriptionExpiresAt,omitempty"`
	PhotoCount            int            `json:"photoCount"`
	RefreshTokens         []RefreshToken `json:"refreshTokens"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
}

// RefreshToken represents a refresh token stored with a user
type RefreshToken struct {
	TokenID     string    `json:"tokenId"`
	HashedToken string    `json:"hashedToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserResponse represents the user data returned to the client
type UserResponse struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	SubscriptionTier string    `json:"subscriptionTier"`
	PhotoCount       int       `json:"photoCount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ToUserResponse converts a User to a UserResponse (hides sensitive data)
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		Name:             u.Name,
		SubscriptionTier: u.SubscriptionTier,
		PhotoCount:       u.PhotoCount,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}
