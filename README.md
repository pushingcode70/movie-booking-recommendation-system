A movie booking platform designed to provide a complete cinema booking experience and movie recommendations.

## Features

- User registration and login
- Email OTP verification and password reset
- Movie browsing and TMDB search
- Razorpay payments
- Email ticket generation
- Wishlist and watched movies
- Semantic movie recommendations
- Admin dashboard

## Tech Stack

| Technology | Purpose |
|---|---|
| Go 1.22+ | Backend |
| Gin | HTTP framework |
| GORM | ORM and migrations |
| PostgreSQL 15+ | Database |
| pgvector | Vector search |
| Python 3.10+ | Embedding service |
| FastAPI | Embedding API |
| sentence-transformers | Embeddings |
| BAAI/bge-base-en-v1.5 | Embedding model |
| TMDB API | Movie data |
| Razorpay | Payments |
| SMTP | Email |
| HTML/CSS/JavaScript | Frontend |

## Project Structure

```text
.
├── cmd/seed/              # Movie catalogue seeding
├── config/                # Application configuration
├── database/              # Database connection
├── dto/                   # Request/response structures
├── handlers/              # HTTP handlers
├── middleware/            # Authentication and authorization
├── models/                # GORM models
├── repositories/          # Database access
├── routes/                # API routes
├── seed/                  # TMDB export processing
├── services/              # Business logic
├── templates/             # Email templates
├── utils/                 # Utility functions
├── embedding-service/     # Python embedding service
├── frontend/              # Frontend
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Git
- Go 1.22+
- Python 3.10+
- PostgreSQL 15+
- PostgreSQL pgvector extension

You also need:

- TMDB API key
- Razorpay credentials
- SMTP credentials

## PostgreSQL Setup

Create the database:

```sql
CREATE DATABASE movie_booking;
```

Connect to it:

```sql
\c movie_booking
```

Enable pgvector:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

Make sure PostgreSQL is running before starting the backend.

GORM AutoMigrate creates missing tables and applies compatible schema changes when the backend starts.

## Clone and Configure

Clone the repository:

```bash
git clone https://github.com/pushingcode70/movie-booking-recommendation-system.git
cd movie-booking-recommendation-system
```

Create the environment file:

```bash
cp .env.example .env
```

Update `.env` with:

- PostgreSQL credentials
- TMDB API key
- Razorpay credentials
- SMTP credentials
- JWT secret
- Embedding service URL

## Embedding Service

The recommendation system uses:

```text
BAAI/bge-base-en-v1.5
```

Setup:

```bash
cd embedding-service
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python3 main.py
```

Embedding service:

```text
http://localhost:8001
```

## Movie Catalogue Seeding

Movies are seeded from the TMDB daily movie export.

Export format:

```text
https://files.tmdb.org/p/exports/movie_ids_MM_DD_YYYY.json.gz
```

Download the required export and run:

```bash
go run ./cmd/seed /path/to/movie_ids_MM_DD_YYYY.json.gz
```

The seed process:

```text
TMDB daily export
       ↓
Read movie IDs
       ↓
Select top 20,000 movies by popularity
       ↓
Fetch movie details and credits
       ↓
Store movies and genres
       ↓
Build recommendation text
       ↓
Generate embeddings
       ↓
Store embeddings in pgvector
```

Existing movies and embeddings are skipped.

## Backend

From the project root:

```bash
go mod download
go run main.go
```

Backend:

```text
http://localhost:8000
```

GORM automatically creates and migrates the required tables when the backend starts.

## Frontend

From the project root:

```bash
python3 -m http.server 3000 --directory frontend
```

Frontend:

```text
http://localhost:3000
```

## Admin Access

Register a user and assign the admin role:

```sql
UPDATE users
SET role = 'admin'
WHERE email = 'your_email@example.com';
```

Admin dashboard:

```text
http://localhost:3000/#/admin
```

## Running the Project

### Terminal 1 — PostgreSQL

Make sure PostgreSQL is running.

### Terminal 2 — Embedding Service

```bash
cd embedding-service
source venv/bin/activate
python3 main.py
```

### Terminal 3 — Backend

```bash
go run main.go
```

### Terminal 4 — Frontend

```bash
python3 -m http.server 3000 --directory frontend
```

## Local URLs

| Service | URL |
|---|---|
| Frontend | http://localhost:3000 |
| Backend | http://localhost:8000 |
| Embedding Service | http://localhost:8001 |
