# PhotoBooth Backend - Architecture Documentation

> **Last Updated:** January 2026  
> **Version:** 1.0.0  
> **Status:** MVP Development

---

## 📋 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Why Golang](#why-golang)
- [System Architecture](#system-architecture)
- [Database Design](#database-design)
- [API Endpoints](#api-endpoints)
- [Authentication Flow](#authentication-flow)
- [File Upload Flow](#file-upload-flow)
- [Payment & Subscription Flow](#payment--subscription-flow)
- [Folder Structure](#folder-structure)
- [Security Considerations](#security-considerations)

---

## 🎯 Overview

**PhotoBooth** is a web-based application that recreates the classic photobooth experience. Users capture 3 photos in sequence, which are automatically combined into a collage that can be downloaded, shared, and stored in their personal gallery.

### Key Business Logic

- **Free Tier**: 5 collages max, watermark enabled
- **Premium Tier**: Unlimited collages, no watermark
- **Storage**: Month-wise organization (e.g., January 2026, February 2026)
- **Sharing**: Public shareable links with optional expiry

---

## 🛠️ Tech Stack

| Component            | Technology               | Reasoning                                                      |
| -------------------- | ------------------------ | -------------------------------------------------------------- |
| **Backend Language** | Golang (Go 1.21+)        | High performance, excellent concurrency, fast image processing |
| **HTTP Framework**   | Gin / Fiber              | Fast, lightweight, production-ready with middleware support    |
| **Database**         | Couchbase                | NoSQL flexibility, excellent performance for JSON documents    |
| **Object Storage**   | Cloudflare R2            | **FREE 10GB permanently**, no egress fees, S3-compatible       |
| **Payment Gateway**  | Stripe                   | Industry standard, official Go SDK, webhook support            |
| **Authentication**   | JWT (golang-jwt)         | Stateless, scalable, short-lived access + refresh tokens       |
| **Image Processing** | imaging / bimg (libvips) | Native Go libraries, extremely fast, low memory usage          |
| **Validation**       | go-playground/validator  | Standard Go validation library                                 |

### Key Go Libraries

- **gin-gonic/gin** - HTTP web framework
- **aws/aws-sdk-go-v2** - S3-compatible SDK for R2
- **golang-jwt/jwt** - JWT implementation
- **stripe/stripe-go** - Official Stripe SDK
- **disintegration/imaging** - Image manipulation
- **couchbase/gocb** - Official Couchbase Go SDK
- **go-playground/validator** - Request validation

---

## 🚀 Why Golang

### Performance Advantages

- **Compiled binary** - No runtime overhead, direct machine code execution
- **Fast image processing** - 10-50x faster than Node.js for image operations
- **Low memory footprint** - Critical when processing multiple large images
- **Native concurrency** - Goroutines handle simultaneous uploads efficiently
- **Quick startup time** - Instant server start (no framework initialization)

### Development Advantages

- **Single binary deployment** - No dependencies, no node_modules
- **Strong typing** - Compile-time error catching
- **Standard library** - Batteries included (HTTP, JSON, crypto, image)
- **Simple dependency management** - Go modules (go.mod)
- **Cross-compilation** - Build for any platform from any platform

### Cost & Scalability

- **Lower server costs** - Less CPU/memory usage = smaller instances
- **Horizontal scaling** - Stateless design scales easily
- **Resource efficient** - Can handle more requests per server
- **Built for microservices** - Easy to split into services later

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────────-┐
│                         FRONTEND (Next.js)                       │
│  - Camera capture UI                                             │
│  - Collage generation (Canvas)                                   │
│  - User dashboard                                                │
│  - Sharing interface                                             │
└────────────────┬────────────────────────────────────────────────-┘
                 │
                 │ HTTP/REST APIs (JSON)
                 │
┌────────────────▼────────────────────────────────────────────────-┐
│                    BACKEND (Golang + Gin)                        │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │   Auth       │  │   Users      │  │   Photos     │            │
│  │   Handler    │  │   Handler    │  │   Handler    │            │
│  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐                              │
│  │   Sharing    │  │   Payments   │                              │
│  │   Handler    │  │   Handler    │                              │
│  └──────────────┘  └──────────────┘                              │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Middleware Layer                            │    │
│  │  - JWT Auth    - Rate Limit    - Subscription Check      │    │
│  │  - CORS        - Logging       - Error Recovery          │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Service Layer                               │    │
│  │  - Auth Service      - Photo Service                     │    │
│  │  - User Service      - Share Service                     │    │
│  │  - Payment Service   - Image Processing Service          │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │              Repository Layer                            │    │
│  │  - User Repository    - Photo Repository                 │    │
│  │  - Share Repository   - Subscription Repository          │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
└───────┬──────────────────────────┬──────────────────────────────-┘
        │                          │
        │                          │
        ▼                          ▼
┌───────────────────┐    ┌─────────────────────┐
│   Couchbase DB    │    │  Cloudflare R2      │
│                   │    │  (Object Storage)   │
│  - users          │    │                     │
│  - photos         │    │  - Collage images   │
│  - shares         │    │  - Compressed JPGs  │
│  - subscriptions  │    │  - Watermarked      │
└───────────────────┘    └─────────────────────┘

        │
        ▼
┌───────────────────┐
│   Stripe API      │
│                   │
│  - Subscriptions  │
│  - Webhooks       │
│  - Payment Intent │
└───────────────────┘
```

### Architecture Layers

**1. Handler Layer (API Controllers)**

- Handles HTTP requests and responses
- Request validation and parsing
- Response formatting
- Delegates business logic to services

**2. Middleware Layer**

- JWT authentication verification
- Rate limiting enforcement
- Subscription tier checking
- CORS handling
- Request logging
- Panic recovery

**3. Service Layer (Business Logic)**

- Core business rules
- Data transformation
- Image processing
- Watermark application
- Token generation
- Payment processing logic

**4. Repository Layer (Data Access)**

- Database operations (CRUD)
- Query building
- Data mapping
- Transaction handling
- Connection pooling

---

## 🗄️ Database Design

### Couchbase Structure

**Single Bucket:** `photobooth_data`

**Collections:**

1. `users` - User profiles and subscription info
2. `photos` - Photo metadata (actual images in R2)
3. `shares` - Share tokens and access tracking
4. `subscriptions` - Subscription details and history

---

### 1. Users Collection

**Document Key Pattern:** `user::{userId}`

```json
{
  "docType": "user",
  "id": "user::123abc",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "passwordHash": "$2b$10$...", // bcrypt hash
  "subscriptionTier": "free", // "free" | "premium"
  "subscriptionStatus": "active", // "active" | "cancelled" | "expired"
  "stripeCustomerId": "cus_...", // Stripe customer ID
  "stripeSubscriptionId": null, // Stripe subscription ID (null for free users)
  "subscriptionExpiresAt": null, // ISO timestamp or null
  "photoCount": 3, // Total collages created
  "refreshTokens": [
    {
      "tokenId": "refresh::xyz789",
      "hashedToken": "$2b$10$...",
      "expiresAt": "2026-01-17T12:00:00Z",
      "createdAt": "2026-01-10T12:00:00Z"
    }
  ],
  "createdAt": "2026-01-01T10:00:00Z",
  "updatedAt": "2026-01-10T14:30:00Z"
}
```

**Indexes:**

- `idx_email` - ON users(email) WHERE docType = "user"
- `idx_stripeCustomerId` - ON users(stripeCustomerId)

---

### 2. Photos Collection

**Document Key Pattern:** `photo::{photoId}`

```json
{
  "docType": "photo",
  "id": "photo::xyz789",
  "userId": "user::123abc",
  "collageUrl": "https://r2.cloudflare.com/photobooth/xyz789.jpg",
  "thumbnailUrl": "https://r2.cloudflare.com/photobooth/xyz789_thumb.jpg",
  "month": "2026-01", // For month-wise grouping
  "year": 2026,
  "timestamp": "2026-01-10T14:30:00Z",
  "hasWatermark": true, // true for free users
  "metadata": {
    "width": 400,
    "height": 1200,
    "fileSize": 245000, // bytes
    "format": "jpg",
    "layout": "vertical_strip" // "vertical_strip" (MVP)
  },
  "shareCount": 0, // How many times shared
  "viewCount": 0, // How many times viewed (via share link)
  "isDeleted": false, // Soft delete flag
  "createdAt": "2026-01-10T14:30:00Z",
  "updatedAt": "2026-01-10T14:30:00Z"
}
```

**Indexes:**

- `idx_userId_month` - ON photos(userId, month) WHERE docType = "photo" AND isDeleted = false
- `idx_userId_timestamp` - ON photos(userId, timestamp) WHERE docType = "photo" AND isDeleted = false
- `idx_createdAt` - ON photos(createdAt) WHERE docType = "photo"

**Query Examples:**

```sql
-- Get all photos for a user in January 2026
SELECT * FROM photobooth_data
WHERE docType = "photo"
  AND userId = "user::123abc"
  AND month = "2026-01"
  AND isDeleted = false
ORDER BY timestamp DESC;

-- Get user's total photo count
SELECT COUNT(*) as total FROM photobooth_data
WHERE docType = "photo"
  AND userId = "user::123abc"
  AND isDeleted = false;
```

---

### 3. Shares Collection

**Document Key Pattern:** `share::{token}`

```json
{
  "docType": "share",
  "id": "share::abc123token",
  "token": "abc123token", // Unique 12-char token
  "photoId": "photo::xyz789",
  "userId": "user::123abc", // Owner of the photo
  "expiresAt": "2026-02-10T14:30:00Z", // Optional expiry
  "isExpired": false, // Computed flag
  "viewCount": 42, // Track views
  "lastViewedAt": "2026-01-15T10:20:00Z",
  "createdAt": "2026-01-10T14:35:00Z"
}
```

**Indexes:**

- `idx_token` - ON shares(token) WHERE docType = "share"
- `idx_photoId` - ON shares(photoId) WHERE docType = "share"
- `idx_expiresAt` - ON shares(expiresAt) WHERE isExpired = false

---

### 4. Subscriptions Collection

**Document Key Pattern:** `subscription::{userId}::{timestamp}`

```json
{
  "docType": "subscription",
  "id": "subscription::user123::1673352000",
  "userId": "user::123abc",
  "stripeSubscriptionId": "sub_...",
  "stripePriceId": "price_...",
  "plan": "premium_monthly", // "premium_monthly" | "premium_yearly"
  "status": "active", // "active" | "cancelled" | "past_due"
  "currentPeriodStart": "2026-01-10T00:00:00Z",
  "currentPeriodEnd": "2026-02-10T00:00:00Z",
  "cancelAtPeriodEnd": false,
  "createdAt": "2026-01-10T12:00:00Z",
  "updatedAt": "2026-01-10T12:00:00Z"
}
```

**Indexes:**

- `idx_userId_status` - ON subscriptions(userId, status) WHERE docType = "subscription"
- `idx_stripeSubscriptionId` - ON subscriptions(stripeSubscriptionId)

---

## 🔌 API Endpoints

For detailed API documentation including all endpoints, request/response examples, error codes, and authentication details, see:

**[📖 API.md](./API.md)** - Complete API Documentation

---

## 🔐 Authentication Flow

### 1. Signup/Login Flow

```
Frontend                Backend                 Couchbase
   |                       |                        |
   |--1. POST /signup----->|                        |
   |   (email, password)   |                        |
   |                       |--2. Hash password      |
   |                       |--3. Create user------->|
   |                       |                        |
   |                       |--4. Generate tokens    |
   |                       |--5. Store refresh----->|
   |                       |                        |
   |<--6. Return tokens----|                        |
   |   (access + refresh)  |                        |
```

### 2. Token Structure

**Access Token (expires in 15 minutes):**

```json
{
  "sub": "user::123abc", // userId
  "email": "user@example.com",
  "tier": "free", // "free" | "premium"
  "iat": 1673352000,
  "exp": 1673352900 // 15 min later
}
```

**Refresh Token (expires in 7 days):**

```json
{
  "sub": "user::123abc",
  "tokenId": "refresh::xyz789",
  "iat": 1673352000,
  "exp": 1673956800 // 7 days later
}
```

### 3. Token Refresh Flow

```
Frontend                Backend                 Couchbase
   |                       |                        |
   |--1. POST /refresh---->|                        |
   |   (refreshToken)      |                        |
   |                       |--2. Verify token       |
   |                       |--3. Check DB---------->|
   |                       |<--4. Token valid-------|
   |                       |                        |
   |                       |--5. Generate new       |
   |                       |    access token        |
   |                       |                        |
   |<--6. Return tokens----|                        |
   |   (new access token)  |                        |
```

### 4. Protected Route Authentication

**Authentication Middleware**

- All protected endpoints use JWT authentication middleware
- Middleware validates access token
- Extracts user information (userId, email, tier)
- Attaches user context to request
- Returns 401 Unauthorized if token invalid/expired

---

## 📤 File Upload Flow

### Proxy Approach (Backend handles upload to R2)

```
Frontend                Backend                 R2 Storage         Couchbase
   |                       |                        |                 |
   |--1. POST /photos----->|                        |                 |
   |   (image blob)        |                        |                 |
   |                       |--2. Check photo limit->|                 |
   |                       |   (free: max 5)        |                 |
   |                       |                        |                 |
   |                       |--3. Process image      |                 |
   |                       |   - Compress           |                 |
   |                       |   - Add watermark?     |                 |
   |                       |   - Generate thumbnail |                 |
   |                       |                        |                 |
   |                       |--4. Upload to R2------>|                 |
   |                       |                        |                 |
   |                       |<--5. Get URLs----------|                 |
   |                       |                        |                 |
   |                       |--6. Save metadata-------------------->   |
   |                       |                        |                 |
   |                       |--7. Increment count------------------>   |
   |                       |                        |                 |
   |<--8. Response---------|                        |                 |
   |   (photoId, URLs)     |                        |                 |
```

### Image Processing Pipeline

**Processing Steps:**

1. **Validation**

   - Check subscription tier and photo count
   - Throw paywall error if free user exceeded 5 photos
   - Validate image format and size (max 10MB)

2. **Image Compression**

   - Convert to JPEG format
   - Apply 85% quality (balance size/quality)
   - Reduce file size for faster uploads/downloads

3. **Watermark Application**

   - Add PhotoBooth watermark for free users
   - Skip watermark for premium users
   - Position: bottom-right corner with transparency

4. **Thumbnail Generation**

   - Resize to 200x600 pixels (scaled proportion)
   - Used for gallery preview
   - Significant file size reduction

5. **Upload to R2**

   - Generate unique photo ID
   - Upload full collage image
   - Upload thumbnail image
   - Use concurrent goroutines for parallel upload

6. **Save Metadata**

   - Store photo metadata in Couchbase
   - Include URLs, timestamp, month, watermark status
   - Link to user document

7. **Update User Stats**
   - Increment user's photo count
   - Used for subscription limit enforcement

---

## 💳 Payment & Subscription Flow

### 1. Checkout Flow

```
Frontend                Backend                 Stripe
   |                       |                        |
   |--1. Click upgrade---->|                        |
   |                       |                        |
   |                       |--2. Create checkout--->|
   |                       |<--3. Session URL-------|
   |                       |                        |
   |<--4. Redirect URL-----|                        |
   |                       |                        |
   |--5. User completes payment on Stripe---------->|
   |                       |                        |
   |                       |<--6. Webhook-----------|
   |                       |   (payment success)    |
   |                       |                        |
   |                       |--7. Update user        |
   |                       |    subscription        |
   |                       |                        |
   |<--8. Redirect---------|                        |
   |   to success page     |                        |
```

### 2. Webhook Handling (Critical)

**Webhook Security & Processing:**

1. **Signature Verification**

   - Verify webhook signature using Stripe webhook secret
   - Prevents unauthorized/fake webhook calls
   - Reject requests with invalid signatures

2. **Event Type Handling**

   - `checkout.session.completed` - Activate new subscription
   - `customer.subscription.updated` - Update subscription status
   - `customer.subscription.deleted` - Cancel subscription
   - `invoice.payment_failed` - Handle payment failure

3. **Idempotency**

   - Track processed webhook event IDs
   - Prevent duplicate processing
   - Handle Stripe webhook retries gracefully

4. **Database Updates**
   - Update user subscription tier
   - Update subscription expiry date
   - Update Stripe customer/subscription IDs
   - Send confirmation email

### 3. Subscription Validation Middleware

**Paywall Enforcement:**

1. **Check User Tier**

   - Extract user info from JWT token
   - Check if user is on free or premium tier

2. **Validate Photo Count**

   - If free tier: check if photoCount >= 5
   - If premium tier: allow unlimited

3. **Return Paywall Error**

   - HTTP 403 Forbidden
   - Error code: `PAYWALL_LIMIT`
   - Include upgrade URL in response
   - Frontend displays upgrade prompt

4. **Applied To Endpoints**
   - POST /api/photos/upload
   - Any endpoint that creates photos

---

## 📁 Folder Structure

```
photobooth-backend/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point, server initialization
│
├── internal/                          # Private application code (not importable)
│   ├── api/
│   │   ├── handlers/                  # HTTP handlers (controllers)
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── photo_handler.go
│   │   │   ├── share_handler.go
│   │   │   └── payment_handler.go
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go                # JWT authentication middleware
│   │   │   ├── subscription.go        # Subscription limit checking
│   │   │   ├── ratelimit.go           # Rate limiting
│   │   │   ├── cors.go                # CORS configuration
│   │   │   ├── logger.go              # Request logging
│   │   │   └── recovery.go            # Panic recovery
│   │   │
│   │   ├── routes/
│   │   │   └── router.go              # Route definitions and middleware setup
│   │   │
│   │   └── dto/                       # Request/response DTOs
│   │       ├── auth_dto.go
│   │       ├── photo_dto.go
│   │       ├── share_dto.go
│   │       └── payment_dto.go
│   │
│   ├── services/                      # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── photo_service.go
│   │   ├── share_service.go
│   │   ├── payment_service.go
│   │   ├── image_service.go           # Image processing (compression, watermark)
│   │   └── email_service.go           # Email notifications
│   │
│   ├── repositories/                  # Database access layer
│   │   ├── user_repository.go
│   │   ├── photo_repository.go
│   │   ├── share_repository.go
│   │   └── subscription_repository.go
│   │
│   ├── models/                        # Domain models (structs)
│   │   ├── user.go
│   │   ├── photo.go
│   │   ├── share.go
│   │   └── subscription.go
│   │
│   └── storage/                       # External storage integrations
│       └── r2_client.go               # Cloudflare R2 client
│
├── pkg/                               # Public packages (reusable across projects)
│   ├── database/
│   │   └── couchbase.go               # Couchbase connection and setup
│   │
│   ├── jwt/
│   │   ├── token.go                   # JWT generation and validation
│   │   └── claims.go                  # Custom JWT claims
│   │
│   ├── validator/
│   │   └── validator.go               # Input validation utilities
│   │
│   └── utils/
│       ├── hash.go                    # Password hashing (bcrypt)
│       ├── random.go                  # Random token generation
│       └── time.go                    # Time/date utilities
│
├── config/
│   ├── config.go                      # Configuration loader
│   └── config.yaml                    # Configuration file (optional)
│
├── scripts/
│   ├── seed.go                        # Database seeding
│   └── migrate.go                     # Database migrations
│
├── docs/
│   └── api.md                         # API documentation
│
├── .env.example                       # Environment variables template
├── .gitignore
├── go.mod                             # Go module dependencies
├── go.sum                             # Dependency checksums
├── Dockerfile                         # Docker container definition
├── docker-compose.yml                 # Docker Compose for local development
├── Makefile                           # Build and run commands
├── ARCHITECTURE.md                    # This file
└── README.md                          # Project documentation
```

### Folder Structure Rationale

**cmd/** - Entry points

- Contains main applications
- `cmd/server/main.go` - HTTP server
- Minimal logic, delegates to internal packages

**internal/** - Private application code

- Not importable by other projects
- Contains all application-specific logic
- Organized by layer (handlers → services → repositories)

**pkg/** - Public/reusable packages

- Can be imported by other projects
- Database, JWT, validation utilities
- Generic, non-business-specific code

**config/** - Configuration management

- Environment-based configuration
- Struct-based config loading
- Validation of required configs

**Standard Go Project Layout**

- Follows Go community conventions
- Clear separation of concerns
- Easy to navigate and understand
- Scalable for future growth

---

## 🔒 Security Considerations

### 1. Authentication Security

- ✅ Passwords hashed with bcrypt (10 rounds)
- ✅ Refresh tokens stored as hashed values in DB
- ✅ Access tokens are short-lived (15 minutes)
- ✅ Refresh tokens invalidated on logout
- ✅ JWT secrets stored in environment variables

### 2. API Security

- ✅ Rate limiting (100 requests/minute per IP)
- ✅ CORS configured for frontend domain only
- ✅ Helmet.js for security headers
- ✅ Input validation with class-validator
- ✅ SQL injection prevention (parameterized queries)

### 3. File Upload Security

- ✅ File type validation (only images)
- ✅ File size limits (max 10MB)
- ✅ Images processed server-side (prevent malicious files)
- ✅ Unique file names (prevent overwrites)
- ✅ R2 credentials never exposed to frontend

### 4. Payment Security

- ✅ Webhook signature verification
- ✅ Idempotent webhook handling
- ✅ Stripe keys in environment variables
- ✅ Customer IDs stored securely

### 5. Data Privacy

- ✅ User photos only accessible by owner
- ✅ Share tokens are random and unpredictable
- ✅ Soft deletes (can be recovered)
- ✅ Optional share link expiry

---

**Last Updated:** January 10, 2026  
**Maintained by:** Backend Team  
**Language:** Golang 1.21+  
**Version:** 1.0.0 (MVP)
