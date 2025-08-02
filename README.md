
# Chirpy

Chirpy is a simple Twitter-like web server and REST API for posting short messages ("chirps") and managing users. It is designed for learning, experimentation, and as a backend for small social applications.

## Why use Chirpy?
- Lightweight, easy-to-understand Go codebase
- Demonstrates user authentication, JWT, refresh tokens, and webhooks
- Great for learning Go web development and RESTful API design

## Installation & Running
1. **Clone the repo:**
   ```bash
   git clone https://github.com/nbvanting/chirpy.git
   cd chirpy
   ```
2. **Set up environment variables:**
   - `DB_URL`: PostgreSQL connection string
   - `PLATFORM`: Set to `dev` for development
   - `TOKEN_SECRET`: Secret for JWT signing
   - `POLKA_KEY`: Secret for webhook authentication
3. **Run the server:**
   ```bash
   go run .
   ```
   The server runs on port 8080 by default.

## API Endpoints

### Health & Metrics
- `GET /api/healthz` — Health check endpoint
- `GET /admin/metrics` — Admin metrics (number of file server hits)
- `POST /admin/reset` — Reset metrics and database (dev only)

### User Management
- `POST /api/users` — Register a new user (`email`, `password`)
- `PUT /api/users` — Update user email/password (JWT required)
- `POST /api/login` — Login, returns JWT and refresh token

### Chirps
- `POST /api/chirps` — Create a new chirp (JWT required)
- `GET /api/chirps` — List all chirps (optional `author_id`, `sort` query params)
- `GET /api/chirps/{chirpID}` — Get a single chirp by ID
- `DELETE /api/chirps/{chirpID}` — Delete a chirp (JWT required, must own chirp)

### Auth Tokens
- `POST /api/refresh` — Exchange refresh token for new JWT
- `POST /api/revoke` — Revoke a refresh token

### Webhooks
- `POST /api/polka/webhooks` — Handle Polka payment webhooks (upgrade user to Chirpy Red)

## Static Files
- `GET /app/*` — Serves static files from the project root (e.g., `index.html`)

---
Thank you to boot.dev for a great walkthrough of this project!
