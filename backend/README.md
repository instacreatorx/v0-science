# Medium-like Blog Backend

Go/Gin backend with GORM, PostgreSQL, SOLID layering, and explicit auth/article state machines.

## Architecture

```
handlers → services → repositories → GORM/PostgreSQL
                ↳ services/state (FSM guards)
```

## Setup

1. Go 1.21+
2. Environment variables:

```
DATABASE_URL=postgresql://...
JWT_SECRET=your_secret_key
PORT=8080
CORS_ORIGINS=http://localhost:3000
GIN_MODE=release
```

3. Run (when online, run `go mod tidy` once):

```bash
cd backend
go run main.go
```

## Auth state machine

- OTP → access JWT (15 min) + refresh token (30 days, hashed in DB)
- `POST /api/auth/logout` revokes refresh token
- `POST /api/auth/refresh` rotates refresh token

## Article state machine

States: `draft` → `published` ↔ `draft`, `archived` → `draft`

- `POST /api/articles` creates a **draft**
- `POST /api/articles/:id/publish|unpublish|archive` — invalid transitions return `409`

## API Endpoints

### Auth
- `POST /api/auth/send-otp`
- `POST /api/auth/verify-otp` — returns `{ token, refresh_token, user }`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`

### Users
- `GET /api/users/me`
- `PUT /api/users/me`
- `GET /api/users/:id` — optional auth adds `is_following`
- `GET /api/users/:id/articles` — published only
- `POST /api/users/:id/follow`

### Articles
- `GET /api/articles` — published only
- `GET /api/articles/me?status=draft|published|archived`
- `GET /api/articles/slug/:slug`
- `GET /api/articles/:id`
- `GET /api/articles/trending`
- `GET /api/articles/search?q=`
- `POST /api/articles` — creates draft
- `PUT /api/articles/:id`
- `DELETE /api/articles/:id`
- `POST /api/articles/:id/publish|unpublish|archive`
- `POST /api/articles/:id/like|bookmark`

### Social
- `GET /api/feed` — followed authors (auth required)
- `GET /api/bookmarks`
- Comments: `GET|POST|DELETE /api/articles/:id/comments`

### Teams (Phase 2)
- `POST /api/teams`
- `GET /api/teams/mine`
- `GET /api/teams/slug/:slug`
- `PUT /api/teams/:id`
- `GET|POST|DELETE /api/teams/:id/members`
- `POST /api/teams/:id/verify-request`

### Admin
- `GET /api/admin/verification-requests`
- `POST /api/admin/verification-requests/:id/approve`
- `POST /api/admin/verification-requests/:id/reject`

Set a user's `role` to `super_admin` in the database for admin access.

## Tests

```bash
go test ./services/state/...
```

## Docker

```bash
docker build -t blog-backend .
docker run -p 8080:8080 -e DATABASE_URL=... -e JWT_SECRET=... blog-backend
```
