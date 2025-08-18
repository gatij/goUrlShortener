# URL Shortener

A simple, scalable URL shortener service written in Go.

## Features (Phase 1)

- Shorten long URLs to concise, easy-to-share links
- Redirect shortened URLs to their original destinations
- Basic metrics tracking (ex - Top 3 most hot domains)
- In-memory storage for URLs and metrics
- HTTPS enforcement and validation
- Infinite loop prevention

## Getting Started

### Prerequisites
- Go 1.22 or higher
- Git

### Installation

1. Clone the repository
```bash
git clone https://github.com/gatij/goUrlShortener.git
cd goUrlShortener
```

2. Install dependencies
```bash
go mod download
```

3. Create a `.env` file in the root directory
```
PORT=3000
BASE_URL=http://localhost:3000
CODE_LENGTH=6
```

4. Run the application
```bash
go run cmd/server/main.go
```

### Docker

You can also run the application using Docker:

1. Build and start the container
```bash
docker-compose up -d
```

2. Stop the container
```bash
docker-compose down
```

Alternatively, you can build and run the Docker image directly:

```bash
# Build the Docker image
docker build -t urlshortener .

# Run the container
docker run -p 3000:3000 -e PORT=3000 -e BASE_URL=http://localhost:3000 -e CODE_LENGTH=6 urlshortener
```

## API Documentation

### Shorten a URL
```
POST /api/v1/urls
```

Request body:
```json
{
  "url": "https://example.com/very/long/url/that/needs/shortening"
}
```

Response:
```json
{
  "short_code": "ab12cd",
  "short_url": "http://localhost:3000/ab12cd",
  "original_url": "https://example.com/very/long/url/that/needs/shortening"
}
```

### Get Top Domains
```
GET /api/v1/metrics/domains?limit=3
```

Response:
```json
{
  "top_domains": [
    {
      "domain": "example.com",
      "shorten_count": 42
    },
    {
      "domain": "github.com",
      "shorten_count": 18
    },
    {
      "domain": "google.com",
      "shorten_count": 7
    }
  ],
  "limit": 3
}
```

### Redirect to Original URL
```
GET /{shortCode}
```
Redirects to the original URL associated with the provided short code.

## Project Structure

```
goUrlShortener/
├── cmd/
│   └── server/
│       └── main.go                # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── shortener.go       # URL shortening endpoint
│   │   │   ├── redirect.go        # Redirect endpoint
│   │   │   └── metrics.go         # Metrics endpoint
│   │   ├── middleware/
│   │   │   └── logging.go         # Basic logging middleware
│   │   └── router.go              # Route setup
│   ├── service/
│   │   ├── shortener.go           # URL shortening logic
│   │   └── metrics.go             # Domain metrics logic
│   ├── storage/
│   │   ├── url/
│   │   │   ├── interface.go       # URLStorage interface definition
│   │   │   └── memory.go          # In-memory implementation
│   │   ├── metrics/
│   │   │   ├── interface.go       # MetricsStorage interface
│   │   │   └── memory.go          # In-memory implementation
│   │   └── factory/               # Factory to create storage based on config
│   └── model/
│       ├── url.go                 # URL data structure
│       └── domainMetrics.go       # Metrics data structure
├── pkg/
│   └── utils/
│       ├── validator.go           # URL validation utilities
│       └── generator.go           # Short URL generation algorithm
├── config/
│   └── config.go                  # Configuration with storage selection
├── scripts/                       # Utility scripts
└── ArchitectureUrlShortner.png    # Architecture diagram
```

## Development Roadmap

### Phase 2: Persistence
- File-based storage implementation
- Data persistence across application restarts
- URL expiration support
- Bulk import/export functionality

### Phase 3: Scalability
- Redis-based storage for improved performance
- Distributed architecture support
- Horizontal scaling capabilities
- Caching layer for high-traffic URLs

## Future Enhancements
- Custom URL support (user-defined short codes)
- Rate limiting for API access
- Visit metrics and analytics
- Authentication and user management
- QR code generation for shortened URLs
- API key authentication
