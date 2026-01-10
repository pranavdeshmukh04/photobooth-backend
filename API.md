# PhotoBooth API Documentation

> **Version:** v1 | **Base URL:** `/api/v1` | **Last Updated:** January 2026

---

## 📋 Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
- [Error Codes](#error-codes)
- [Rate Limiting & Pagination](#rate-limiting--pagination)

---

## 🎯 Overview

RESTful API for PhotoBooth application with JWT authentication, file uploads, sharing, and subscription management.

**Features:** JWT auth • Multipart uploads • Public sharing • Stripe integration • Month-wise organization

**Content Types:** `application/json` (default) • `multipart/form-data` (uploads)

---

## 🔐 Authentication

### Using Authentication

**Header Format:**

```
Authorization: Bearer <access_token>
```

**Token Types:**

- **Access Token:** Short-lived (15 min), used for API requests
- **Refresh Token:** Long-lived (7 days), used to get new access tokens

**Access Token Claims:**

```json
{
  "sub": "user::123abc",
  "email": "user@example.com",
  "tier": "free",
  "iat": 1673352000,
  "exp": 1673352900
}
```

---

## 📡 Endpoints

### Authentication

#### POST /api/auth/signup

Register new user.

**Request:**

```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe"
}
```

**Response (201):**

```json
{
  "user": {
    "id": "user::123abc",
    "email": "user@example.com",
    "name": "John Doe",
    "subscriptionTier": "free",
    "photoCount": 0
  },
  "tokens": {
    "accessToken": "eyJhbGci...",
    "refreshToken": "refresh::xyz..."
  }
}
```

**Validation:** Email (valid format, unique) • Password (min 8 chars, 1 upper, 1 lower, 1 number) • Name (2-50 chars)

---

#### POST /api/auth/login

Authenticate user.

**Request:** `{ "email": "...", "password": "..." }`

**Response (200):** Same as signup response

**Errors:** `401` Invalid credentials • `429` Rate limit exceeded

---

#### POST /api/auth/refresh

Get new access token.

**Request:** `{ "refreshToken": "refresh::xyz..." }`

**Response (200):** `{ "accessToken": "...", "refreshToken": "..." }`

**Note:** Old refresh token is invalidated

---

#### POST /api/auth/logout

**Auth Required**

Logout and invalidate refresh token.

**Request:** `{ "refreshToken": "refresh::xyz..." }`

**Response (200):** `{ "message": "Logged out successfully" }`

---

#### GET /api/auth/me

**Auth Required**

Get current user info.

**Response (200):**

```json
{
  "user": {
    "id": "user::123abc",
    "email": "user@example.com",
    "name": "John Doe",
    "subscriptionTier": "free",
    "photoCount": 3,
    "createdAt": "2026-01-01T10:00:00Z"
  }
}
```

---

### Users

#### GET /api/users/profile

**Auth Required**

Get user profile.

**Response (200):** User object with subscription details

---

#### PATCH /api/users/profile

**Auth Required**

Update profile.

**Request:** `{ "name": "Updated Name" }`

**Updatable:** Name only (email/password changes in future)

---

#### GET /api/users/stats

**Auth Required**

Get user statistics.

**Response (200):**

```json
{
  "stats": {
    "totalPhotos": 5,
    "totalShares": 12,
    "subscriptionTier": "free",
    "photosRemaining": 0,
    "storageUsed": "1.2 MB"
  }
}
```

---

### Photos

#### POST /api/photos/upload

**Auth Required • Subscription Check**

Upload collage image.

**Content-Type:** `multipart/form-data`

**Request:**

```
file: [binary image]
metadata: { "timestamp": "2026-01-10T14:30:00Z" }
```

**Response (201):**

```json
{
  "photo": {
    "id": "photo::xyz789",
    "collageUrl": "https://r2.cloudflare.com/.../xyz789.jpg",
    "thumbnailUrl": "https://r2.cloudflare.com/.../xyz789_thumb.jpg",
    "month": "2026-01",
    "timestamp": "2026-01-10T14:30:00Z",
    "hasWatermark": true
  }
}
```

**Validation:** Max 10MB • JPEG/PNG only • Free tier: 5 photos max

**Errors:** `403` Paywall (see [Paywall Error](#paywall-error)) • `413` File too large • `415` Invalid file type

---

#### GET /api/photos

**Auth Required**

Get user's photo gallery (paginated).

**Query Params:** `page` (default: 1) • `limit` (default: 10, max: 50) • `month` (YYYY-MM)

**Response (200):**

```json
{
  "photos": [
    {
      "id": "photo::xyz789",
      "collageUrl": "...",
      "thumbnailUrl": "...",
      "month": "2026-01",
      "timestamp": "2026-01-10T14:30:00Z",
      "hasWatermark": true,
      "shareCount": 5
    }
  ],
  "pagination": {
    "total": 12,
    "page": 1,
    "limit": 10,
    "hasMore": true
  }
}
```

---

#### GET /api/photos/:photoId

**Auth Required**

Get single photo details.

**Response (200):** Photo object with metadata and stats

**Errors:** `404` Not found or not owned by user

---

#### DELETE /api/photos/:photoId

**Auth Required**

Delete photo (soft delete).

**Response (200):** `{ "message": "Photo deleted successfully" }`

**Note:** Photo count decremented, share links deactivated

---

#### GET /api/photos/months

**Auth Required**

Get list of months with photos.

**Response (200):**

```json
{
  "months": [
    { "month": "2026-01", "count": 12, "label": "January 2026" },
    { "month": "2025-12", "count": 8, "label": "December 2025" }
  ]
}
```

---

### Sharing

#### POST /api/shares/:photoId

**Auth Required**

Generate share link for photo.

**Request:** `{ "expiresInDays": 30 }` (optional, null for no expiry)

**Response (201):**

```json
{
  "share": {
    "token": "abc123token",
    "photoId": "photo::xyz789",
    "shareUrl": "https://photobooth.app/share/abc123token",
    "qrCodeData": "https://photobooth.app/share/abc123token",
    "expiresAt": "2026-02-10T14:30:00Z"
  }
}
```

**Note:** Multiple share links allowed per photo • QR code generated by frontend

---

#### GET /api/shares/:token

**Public (No Auth)**

Get shared photo.

**Response (200):**

```json
{
  "photo": {
    "collageUrl": "...",
    "timestamp": "2026-01-10T14:30:00Z",
    "metadata": { "width": 400, "height": 1200 }
  },
  "shareInfo": {
    "viewCount": 43,
    "expiresAt": "2026-02-10T14:30:00Z"
  }
}
```

**Errors:** `404` Token not found • `410` Expired/revoked

**Note:** View count incremented on access • No user info exposed

---

#### DELETE /api/shares/:token

**Auth Required**

Revoke share link (owner only).

**Response (200):** `{ "message": "Share link revoked successfully" }`

---

### Payments

#### GET /api/payments/plans

**Public**

Get subscription plans.

**Response (200):**

```json
{
  "plans": [
    {
      "id": "free",
      "name": "Free",
      "price": 0,
      "features": ["5 collages", "Watermark included"]
    },
    {
      "id": "premium_monthly",
      "name": "Premium Monthly",
      "price": 4.99,
      "interval": "month",
      "stripePriceId": "price_...",
      "features": ["Unlimited collages", "No watermark", "Priority processing"]
    },
    {
      "id": "premium_yearly",
      "name": "Premium Yearly",
      "price": 49.99,
      "interval": "year",
      "stripePriceId": "price_...",
      "features": ["Unlimited", "No watermark", "Priority", "Save 17%"]
    }
  ]
}
```

---

#### POST /api/payments/checkout

**Auth Required**

Create Stripe checkout session.

**Request:** `{ "priceId": "price_monthly_premium" }`

**Response (200):** `{ "sessionId": "cs_test_...", "checkoutUrl": "https://checkout.stripe.com/..." }`

**Flow:** Call endpoint → Redirect to checkoutUrl → Payment on Stripe → Webhook → Success redirect

---

#### POST /api/payments/webhooks/stripe

**Public (Stripe Signature Verified)**

Stripe webhook endpoint.

**Headers:** `Stripe-Signature: t=...,v1=...`

**Response (200):** `{ "received": true }`

**Events:** `checkout.session.completed` • `customer.subscription.updated` • `customer.subscription.deleted` • `invoice.payment_failed`

---

#### GET /api/payments/subscription

**Auth Required**

Get user's subscription details.

**Response (200):**

```json
{
  "subscription": {
    "tier": "premium",
    "status": "active",
    "plan": "premium_monthly",
    "currentPeriodEnd": "2026-02-10T00:00:00Z",
    "cancelAtPeriodEnd": false
  }
}
```

**Free User Response:** `{ "subscription": { "tier": "free", "status": "active", "photosUsed": 5, "photosLimit": 5 } }`

---

#### POST /api/payments/subscription/cancel

**Auth Required**

Cancel subscription (at period end).

**Response (200):** Subscription object with `cancelAtPeriodEnd: true`

**Note:** Premium access retained until period end → Auto-downgrade to free

---

#### POST /api/payments/portal

**Auth Required**

Get Stripe customer portal link.

**Response (200):** `{ "portalUrl": "https://billing.stripe.com/session/..." }`

**Usage:** Manage payment methods, invoices, billing

---

### Health Checks

#### GET /health

**Public**

Basic health check.

**Response (200):** `{ "status": "ok", "timestamp": "2026-01-10T14:30:00Z" }`

---

#### GET /health/ready

**Public**

Readiness check (includes database).

**Response (200):** `{ "status": "ready", "checks": { "database": "ok", "storage": "ok" }, "timestamp": "..." }`

**Response (503):** `{ "status": "not_ready", "checks": { "database": "error", "storage": "ok" }, ... }`

---

#### GET /health/live

**Public**

Liveness check.

**Response (200):** `{ "status": "live", "timestamp": "..." }`

---

## ❌ Error Codes

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { "field": "additional context" }
  },
  "timestamp": "2026-01-10T14:30:00Z"
}
```

---

### Common Error Codes

| Code                     | Status | Description                           |
| ------------------------ | ------ | ------------------------------------- |
| `INVALID_CREDENTIALS`    | 401    | Wrong email or password               |
| `TOKEN_EXPIRED`          | 401    | Access token expired                  |
| `TOKEN_INVALID`          | 401    | Token malformed or invalid            |
| `UNAUTHORIZED`           | 401    | Authentication required               |
| `FORBIDDEN`              | 403    | Insufficient permissions              |
| `VALIDATION_ERROR`       | 400    | Request validation failed             |
| `INVALID_EMAIL`          | 400    | Email format invalid                  |
| `WEAK_PASSWORD`          | 400    | Password doesn't meet requirements    |
| `EMAIL_TAKEN`            | 409    | Email already registered              |
| `PHOTO_NOT_FOUND`        | 404    | Photo doesn't exist                   |
| `FILE_TOO_LARGE`         | 413    | File exceeds 10MB                     |
| `INVALID_FILE_TYPE`      | 415    | File is not an image                  |
| `SHARE_NOT_FOUND`        | 404    | Share token doesn't exist             |
| `SHARE_EXPIRED`          | 410    | Share link expired                    |
| `PAYMENT_FAILED`         | 402    | Payment processing failed             |
| `INVALID_PRICE_ID`       | 400    | Stripe price ID not found             |
| `WEBHOOK_INVALID`        | 400    | Webhook signature verification failed |
| `SUBSCRIPTION_NOT_FOUND` | 404    | No active subscription                |
| `RATE_LIMIT_EXCEEDED`    | 429    | Too many requests                     |
| `INTERNAL_ERROR`         | 500    | Unexpected server error               |
| `DATABASE_ERROR`         | 500    | Database operation failed             |
| `STORAGE_ERROR`          | 500    | Object storage operation failed       |

---

### Paywall Error

**Most Important Error - Free Tier Limit:**

```json
{
  "error": {
    "code": "PAYWALL_LIMIT",
    "message": "Free tier limit reached. Upgrade to premium for unlimited photos.",
    "details": {
      "currentCount": 5,
      "maxAllowed": 5,
      "upgradeUrl": "/upgrade"
    }
  },
  "timestamp": "2026-01-10T14:30:00Z"
}
```

**Status:** `403 Forbidden`

**Triggered:** Free user tries to upload 6th photo

**Frontend Action:** Display upgrade modal → Show pricing → Redirect to checkout

---

## ⚡ Rate Limiting & Pagination

### Rate Limiting

**Limits:**

- Global: 100 requests/min per IP
- Auth endpoints: 5 requests/min per IP
- Photo uploads: 10 uploads/min per user

**Headers:**

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1673352960
```

**Rate Limit Error (429):**

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later.",
    "details": { "retryAfter": 60 }
  }
}
```

---

### Pagination

**Query Parameters:**

- `page` - Page number (default: 1, min: 1)
- `limit` - Items per page (default: 10, min: 1, max: 50)

**Example:** `GET /api/photos?page=2&limit=20`

**Response Format:**

```json
{
  "photos": [...],
  "pagination": {
    "total": 45,
    "page": 2,
    "limit": 20,
    "totalPages": 3,
    "hasMore": true,
    "hasPrevious": true
  }
}
```

---

## 🔧 Additional Info

### CORS

- **Allowed Origins:** Frontend URL (configured), localhost (dev)
- **Methods:** GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Headers:** Authorization, Content-Type

### Timestamps

- **Format:** ISO 8601 (UTC) - `2026-01-10T14:30:00Z`
- **Fields:** `createdAt`, `updatedAt`, `timestamp`

### Content Types

- **Request:** `application/json` (default) • `multipart/form-data` (uploads)
- **Response:** `application/json` (all endpoints)

---

## 📞 Support

**For API Questions:**

1. Check this API documentation
2. Review [ARCHITECTURE.md](./ARCHITECTURE.md) for system design
3. Check code comments for implementation
4. Create GitHub issue for bugs/features

---

**API Version:** v1 • **Backend:** Golang + Gin • **Last Updated:** January 10, 2026
