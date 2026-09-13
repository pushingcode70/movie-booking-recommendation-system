# CinePass — Movie Booking System

CinePass is a movie booking application built with Go, PostgreSQL, Python, and a vanilla JavaScript frontend.

## Features

- User registration and login
- Email OTP verification and password reset
- Movie browsing and TMDB search
- Theatre and screen management
- Automatic seat generation
- Show scheduling and seat availability
- Razorpay payments
- Email ticket generation
- Wishlist and watched movies
- Semantic movie recommendations
- Admin dashboard
- Postman API collection

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
├── postman/               # Postman collection
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md

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

CREATE DATABASE movie_booking;

Connect to it:

\c movie_booking

Enable pgvector:

CREATE EXTENSION IF NOT EXISTS vector;

Make sure PostgreSQL is running before starting the backend.

GORM handles the database migrations when the backend starts.

## Clone and Configure

Clone the repository:

git clone <repository-url>
cd <project-directory>

Create the environment file:

cp .env.example .env

Update .env with:

- PostgreSQL credentials
- TMDB API key
- Razorpay credentials
- SMTP credentials
- JWT secret
- Embedding service URL

## Embedding Service

The recommendation system uses:

BAAI/bge-base-en-v1.5

Setup:

cd embedding-service

python3 -m venv venv
source venv/bin/activate

pip install -r requirements.txt
python3 main.py

Embedding service:

http://localhost:8001

## Movie Catalogue Seeding

Movies are seeded from the TMDB daily movie export.

Export format:

https://files.tmdb.org/p/exports/movie_ids_MM_DD_YYYY.json.gz

Download the required export and run:

go run ./cmd/seed /path/to/movie_ids_MM_DD_YYYY.json.gz

The seed process:

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

Existing movies and embeddings are skipped.

## Backend

From the project root:

go mod download
go run main.go

Backend:

http://localhost:8000

GORM automatically creates/migrates the required tables when the backend starts.

## Frontend

From the project root:

python3 -m http.server 3000 --directory frontend

Frontend:

http://localhost:3000

Backend:

http://localhost:8000

## Admin Access

Register a user and assign the admin role:

UPDATE users
SET role = 'admin'
WHERE email = 'your_email@example.com';

Admin dashboard:

http://localhost:3000/#/admin

## Running the Project

Terminal 1 — PostgreSQL

Make sure PostgreSQL is running.

Terminal 2 — Embedding Service

cd embedding-service
source venv/bin/activate
python3 main.py

Terminal 3 — Backend

go run main.go

Terminal 4 — Frontend

python3 -m http.server 3000 --directory frontend

## Local URLs

Frontend    http://localhost:3000
Backend     http://localhost:8000
Embedding   http://localhost:8001
