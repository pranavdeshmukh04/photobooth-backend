# 📸 PhotoBooth Backend

> A web-based photobooth application that captures 3 photos and creates shareable collages

---

## 🎯 What is PhotoBooth?

PhotoBooth recreates the classic photobooth experience on the web. Users capture 3 photos in sequence, which are automatically combined into a collage that can be downloaded, shared, and stored in their personal gallery.

**Key Features:**

- 📷 3-photo capture flow with countdown
- 🎨 Automatic vertical strip collage generation
- 💾 Month-wise photo gallery (January 2026, February 2026, etc.)
- 🔗 Shareable links with optional expiry + QR codes
- 💳 Free tier (5 photos) and Premium (unlimited)
- 🔐 JWT authentication with refresh tokens
- ☁️ Cloud storage on Cloudflare R2

---

## 🛠️ Tech Stack

| Component          | Technology     | Why                                   |
| ------------------ | -------------- | ------------------------------------- |
| **Language**       | Golang 1.21+   | Fast, concurrent, efficient           |
| **Framework**      | Gin            | Lightweight HTTP framework            |
| **Database**       | Couchbase      | NoSQL, flexible JSON storage          |
| **Object Storage** | Cloudflare R2  | Free 10GB, S3-compatible, zero egress |
| **Payment**        | Stripe         | Subscription management               |
| **Auth**           | JWT            | Stateless authentication              |
| **Image Proc**     | imaging / bimg | Fast compression and watermarking     |

---

## 📁 Project Structure

```
photobooth-backend/
├── cmd/server/              # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers, middleware, routes
│   ├── services/            # Business logic
│   ├── repositories/        # Database operations
│   ├── models/              # Data structures
│   └── storage/             # R2 client
├── pkg/                     # Reusable packages (DB, JWT, utils)
├── config/                  # Configuration
├── ARCHITECTURE.md          # 📘 Detailed system design
├── API.md                   # 📗 Complete API reference
├── go.mod
└── README.md                # This file
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Couchbase Server
- Cloudflare R2 account
- Stripe account (test mode)

### Installation

```bash
# 1. Clone repository
git clone <repository-url>
cd photobooth-backend

# 2. Install dependencies
go mod download

# 3. Set up environment
cp .env.example .env
# Edit .env with your credentials

# 4. Start Couchbase (Docker)
docker-compose up -d couchbase

# 5. Run application
go run cmd/server/main.go
# Server starts on http://localhost:8080
```

---

## 📚 Documentation

### 📘 [ARCHITECTURE.md](./ARCHITECTURE.md)

**Complete system design and technical details**

- System architecture and component design
- Database schemas (Couchbase collections)
- Authentication flow (JWT tokens)
- File upload flow (image processing pipeline)
- Payment & subscription flow (Stripe webhooks)
- Security considerations
- Folder structure explained

👉 **Read this to understand how the system works**

---

### 📗 [API.md](./API.md)

**Complete API reference for developers**

- All 26 API endpoints documented
- Authentication endpoints (signup, login, refresh, logout)
- Photos endpoints (upload, gallery, sharing)
- Payment endpoints (checkout, webhooks, subscriptions)
- Request/response examples
- Error codes and handling
- Rate limiting and pagination

👉 **Read this to integrate with the API**

---

## 💳 Subscription Tiers

| Tier        | Price  | Photos    | Watermark |
| ----------- | ------ | --------- | --------- |
| **Free**    | $0     | 5 max     | Yes       |
| **Premium** | $4.99  | Unlimited | No        |
| **Yearly**  | $49.99 | Unlimited | No        |

**Note:** Yearly plan saves 17% (equivalent to 2 months free)

---

## 🤝 Contributing

1. Create feature branch from `develop`
2. Make changes with tests
3. Run `make lint` and `make test`
4. Submit pull request

**Branch Convention:**

- `feature/` - New feature
- `fix/` - Bug fix
- `setup/` - Project setup/infrastructure
- `chore/` - Maintenance tasks
- `docs/` - Documentation
- `refactor/` - Code refactoring

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file

---

## 🔗 Links

- **Frontend Repository:** [photobooth-frontend](https://github.com/saumyasrv/photobooth-frontend) (Next.js)

---

## 📞 Support

- **Issues:** Create a GitHub issue
- **Questions:** Check docs first, then open discussion
- **Security:** Email security concerns (don't open public issues)

---
