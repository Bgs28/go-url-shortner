# URL Shortener

A URL Shortener REST API built with native Go and MySQL, featuring short code generation, HTTP redirect, and click tracking.

---

## What I Learned

- HTTP 301 Redirect with `http.Redirect()`
- Random short code generation with `math/rand`
- URL validation with `url.ParseRequestURI()`
- Click counter with atomic SQL update (`clicks = clicks + 1`)
- URL path parameter parsing with `strings.TrimPrefix()`
- Clean Architecture (Handler → Service → Repository)

---

## Tech Stack

| Category   | Tools                     |
|------------|---------------------------|
| Language   | Go (Golang)               |
| Database   | MySQL, database/sql       |
| Config     | godotenv                  |
| Driver     | go-sql-driver/mysql       |

---

## Project Structure

```
08-url-shortener/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/
│   │   └── url_handler.go
│   ├── service/
│   │   └── url_service.go
│   ├── repository/
│   │   └── url_repository.go
│   └── model/
│       └── url.go
├── .env.example
├── .gitignore
└── README.md
```

---

## Endpoints

```
POST   /shorten        → Generate short URL
GET    /:code          → Redirect to original URL
GET    /stats/:code    → Get click statistics
```

---

## How It Works

```
POST /shorten
      |
      v
Validate URL format (must be http/https)
      |
      v
Generate 6-character random short code
      |
      v
Save to MySQL
      |
      v
Return { short_code, short_url, original_url }


GET /:code
      |
      v
Lookup short code in database
      |
      v
Increment click counter (clicks = clicks + 1)
      |
      v
301 Redirect → original URL
```

---

## Database Schema

```sql
CREATE TABLE urls (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code   VARCHAR(10) UNIQUE NOT NULL,
    clicks       INT DEFAULT 0,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## How to Run

**1. Clone and navigate into the project:**
```bash
cd 08-url-shortener
```

**2. Copy and fill in the environment file:**
```bash
cp .env.example .env
```

**3. Fill in `.env`:**
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=url_shortener
APP_PORT=8080
BASE_URL=http://localhost:8080
```

**4. Create the database and table:**
```sql
CREATE DATABASE url_shortener;
USE url_shortener;

CREATE TABLE urls (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code   VARCHAR(10) UNIQUE NOT NULL,
    clicks       INT DEFAULT 0,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**5. Run the project:**
```bash
go run cmd/main.go
```

---

## API Usage

**Shorten a URL:**
```bash
POST http://localhost:8080/shorten
Content-Type: application/json

{
    "original_url": "https://www.google.com"
}
```

Response:
```json
{
    "short_code": "aB3xYz",
    "short_url": "http://localhost:8080/aB3xYz",
    "original_url": "https://www.google.com"
}
```

**Redirect:**
```
Open in browser → http://localhost:8080/aB3xYz
Automatically redirected to → https://www.google.com
```

**Get Stats:**
```bash
GET http://localhost:8080/stats/aB3xYz
```

Response:
```json
{
    "short_code": "aB3xYz",
    "original_url": "https://www.google.com",
    "clicks": 3,
    "created_at": "2026-05-08T10:00:00Z"
}
```

---

## Author

**Katon Bagas** — Information Systems Student | Backend Enthusiast | Future Web3 Architect

> *"Built with consistency and discipline. Learning never stops."*
