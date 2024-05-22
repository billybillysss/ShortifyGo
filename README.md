Sure, here is an updated version of the `README.md` with emojis:

---

# 🚀 ShortifyGo

ShortifyGo is a URL shortening service built with Go and Docker. It uses Redis for data storage and provides an API to shorten and retrieve URLs.

## 📚 Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Endpoints](#endpoints)

## 🛠️ Installation

To get started with ShortifyGo, you'll need to have Docker and Docker Compose installed on your machine.

1. **Clone the repository:**
   ```sh
   git clone https://github.com/billybillysss/shortifygo.git
   cd shortifygo
   ```

2. **Build and start the services:**
   ```sh
   docker-compose up --build
   ```

## 🚀 Usage

Once the services are up and running, you can access the API at `http://localhost:7001`.

### ✂️ Shorten a URL

To shorten a URL, send a POST request to `/api/v1` with the URL in the request body.

Example:
```sh
curl -X POST -d '{"url":"https://www.example.com"}' http://localhost:7001/api/v1
```

### 🔍 Retrieve a URL

To retrieve the original URL, send a GET request to `/{shortId}`.

Example:
```sh
curl http://localhost:7001/retrieve/abc123
```

## 📂 Project Structure

```plaintext
SHORTIFYGO/
├── .data/
├── api/
│   ├── database/
│   │   └── redis.go
│   ├── handlers/
│   │   ├── retriever.go
│   │   └── shortener.go
│   ├── utils/
│   ├── .env
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── db/
│   └── Dockerfile
└── docker-compose.yml
```

### 📄 Description of Files

- `.data/` - Directory for Redis data storage.
- `.env` - Environment variables file.
- `api/` - Contains the source code for the API.
  - `database/redis.go` - Redis client and database operations.
  - `handlers/retriever.go` - Handler for retrieving URLs.
  - `handlers/shortener.go` - Handler for shortening URLs.
  - `utils/` - Utility functions (if any).
  - `.env` - Environment variables for the API.
  - `Dockerfile` - Dockerfile for building the API service.
  - `go.mod` - Go module file.
  - `go.sum` - Go dependencies.
  - `main.go` - Entry point of the API service.
- `db/` - Contains the Dockerfile for the Redis service.
- `docker-compose.yml` - Docker Compose file to run the services.

## 📡 Endpoints

### `POST /api/v1`
- **Description:** Shortens a given URL.
- **Request Body:** JSON containing the URL to be shortened and an optional short code.
  ```json
  {
    "URL": "https://www.example.com",
    "Short": "example"
  }
  ```
- **Response Body:** JSON containing the status, original URL, short code, remaining requests, and limit reset time.
  ```json
  {
    "Status": "success",
    "URL": "https://www.example.com",
    "Short": "example",
    "RemainRequest": 5,
    "LimitRestTime": "2024-05-23T00:00:00Z"
  }
  ```

### `GET /{shortId}`
- **Description:** Retrieves the original URL for a given shortened code.
- **Path Parameters:** `shortId` - The code representing the shortened URL.
