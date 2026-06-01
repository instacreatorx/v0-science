# Medium-like Blog Backend

Go/Gin backend for the Medium-like blog platform.

## Setup

1. Make sure you have Go 1.21+ installed
2. Set environment variables:
   ```
   DATABASE_URL=your_neon_connection_string
   JWT_SECRET=your_secret_key
   PORT=8080
   ```

3. Run the server:
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Users
- `GET /api/users/me` - Get current user (auth required)
- `PUT /api/users/me` - Update current user (auth required)
- `GET /api/users/:id` - Get user by ID
- `GET /api/users/:id/articles` - Get user's articles
- `POST /api/users/:id/follow` - Follow/unfollow user (auth required)

### Articles
- `GET /api/articles` - List articles (paginated)
- `GET /api/articles/trending` - Get trending articles
- `GET /api/articles/search?q=query` - Search articles
- `GET /api/articles/:id` - Get single article
- `POST /api/articles` - Create article (auth required)
- `PUT /api/articles/:id` - Update article (auth required)
- `DELETE /api/articles/:id` - Delete article (auth required)
- `POST /api/articles/:id/like` - Like/unlike article (auth required)
- `POST /api/articles/:id/bookmark` - Bookmark/unbookmark article (auth required)

### Comments
- `GET /api/articles/:id/comments` - Get article comments
- `POST /api/articles/:id/comments` - Add comment (auth required)
- `DELETE /api/articles/:id/comments/:commentId` - Delete comment (auth required)

### Bookmarks
- `GET /api/bookmarks` - Get user's bookmarked articles (auth required)

## Docker

Build and run with Docker:

```bash
docker build -t blog-backend .
docker run -p 8080:8080 -e DATABASE_URL=... -e JWT_SECRET=... blog-backend
```
