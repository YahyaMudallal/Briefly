# Briefly: Backend

This folder contains the server part of the application, built with Go using the native `net/http` package.

## Running the server

If you have [air](https://github.com/air-verse/air) configured, you can launch the server with hot reload using the following command:
```bash
air
```

Otherwise, you can launch the backend by running the `main.go` file at the root of the directory:
```Bash
go run main.go
```

## REST API Routes

The server exposes a REST API structured around versioned endpoints (/v1). Below is the complete list of available routes, their descriptions, and access constraints:

| Method & Route | Description | Authentication |
| :--- | :--- | :--- |
| **GET** `/v1/health` | Checks the server's availability and health status (returns `200 OK`). | None |
| **GET** `/v1/articles` | Fetches a paginated list of news articles using `page` and `limit` query parameters. | None |
| **GET** `/v1/articles/{id}` | Fetches the full details of a specific article by its unique ID. | None |
| **GET** `/v1/comments/{id}` | Fetches a specific comment by its ID in JSON format. | None |
| **GET** `/v1/comments/article/{id}` | Fetches all nested comments and replies associated with a specific article. | None |
| **GET** `/v1/users/{id}` | Fetches public profile information for a specific user. | None |
| **POST** `/v1/users` | Registers a new user account after validating the email format and password strength. Returns a JWT token. | None |
| **POST** `/v1/login` | Authenticates a user with email and password. Returns a JWT token upon success. | None |
| **POST** `/v1/comments` | Publishes a new comment or reply under an article. | **JWT Token** (Header) |
| **POST** `/v1/articles/{id}/upvote` | Casts a positive vote for the specified article. User identity is extracted from the token. | **JWT Token** (Header) |
| **POST** `/v1/articles/{id}/downvote`| Casts a negative vote for the specified article. User identity is extracted from the token. | **JWT Token** (Header) |
| **POST** `/v1/articles/{id}/summary` | Triggers an asynchronous call to the Gemini API to generate a summary for the article. | **JWT Token** (Header) |
| **POST** `/v1/sync-articles` | Triggers the daily synchronization with NewsData.io to fetch 10 new articles. | **Secret Key** (Header) |
| **DELETE** `/v1/comments/{id}` | Deletes a comment. Authorized only for the comment's author or an administrator. | **JWT Token** (Header) |
| **DELETE** `/v1/articles/{id}` | Deletes an article from the database. Restricted to administrators. | **JWT Token** (Header) |
| **DELETE** `/v1/users/{id}` | Deletes a user account from the platform. Restricted to administrators. | **JWT Token** (Header) |

## Design of the backend

The folder structure of the backend is the following:
```Bash
.
├── Dockerfile
├── README.md
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   └── api.go
│   ├── apperrors
│   │   └── errors.go
│   ├── auth
│   │   └── jwt.go
│   ├── clients
│   │   ├── geminiapi.go
│   │   └── newsapi.go
│   ├── database
│   │   └── db.go
│   ├── handlers
│   │   ├── articles.go
│   │   ├── comments.go
│   │   └── users.go
│   ├── jobs
│   │   └── scheduler.go
│   ├── middleware
│   │   ├── authMiddleware.go
│   │   ├── enableCORS.go
│   │   ├── logger.go
│   │   ├── optionalAuthMiddleware.go
│   │   └── recover.go
│   ├── models
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── user.go
│   │   └── vote.go
│   ├── repositories
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── interfaces.go
│   │   ├── user.go
│   │   └── vote.go
│   └── services
│       ├── article.go
│       ├── comment.go
│       └── user.go
├── main.go
└── tmp
    ├── build-errors.log
    └── main
```