# Iltodgeree Mining Contracts API

A RESTful API service for managing and searching mining contracts in Mongolia. Built with Go, Elasticsearch, and PostgreSQL.

## Features

- **Full-Text Search**: Advanced search capabilities across contract text and metadata
- **Multi-Filter Support**: Filter by year, type, resource, location, company, government entity
- **Document Export**: Export search results to DOCX or TSV formats
- **Annotation System**: Retrieve and search contract annotations
- **Aggregations**: Statistical summaries and faceted search
- **Bilingual Support**: Mongolian and English translations
- **Geographic Filtering**: Province and district-based filtering

## Quick Start

### Prerequisites

- Go 1.23.4+
- Elasticsearch 5.x
- PostgreSQL 12+
- Docker (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/iltodgeree-api.git
cd iltodgeree-api

# Install dependencies
go mod download

# Copy environment configuration
cp .env.example .env
# Edit .env with your configuration

# Run the service
go run cmd/front-service/main.go
```

The API will be available at `http://localhost:8080`

### Docker Setup

```bash
# Build the image
docker build -t iltodgeree-api .

# Run the container
docker run -p 8080:8080 --env-file .env iltodgeree-api
```

## Configuration

Create a `.env` file in the project root with the following variables:

```bash
# Elasticsearch Configuration
ELASTICSEARCH_HOST=http://localhost:9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=changeme
ELASTICSEARCH_SECONDARY=iltodgeree_v2.2
ELASTICSEARCH_DOC_MASTER=master
ELASTICSEARCH_DOC_METADATA=metadata

# PostgreSQL Configuration
PGSQL_URL=postgresql://user:password@localhost:5432/iltodgeree

# File System Paths
DOCUMENT_PATH=/var/iltodgeree/documents
TEMPLATE_PATH=/var/iltodgeree/templates
STORAGE_PATH=/var/iltodgeree/storage
PUBLIC_URL=https://api.iltodgeree.mn

# Frontend Configuration
FRONT_END_URL=http://localhost:3000

# Application Mode
GIN_MODE=release
```

## API Usage

### Search Contracts

```bash
# Simple text search
curl "http://localhost:8080/api/search?q=gold+mining"

# Filter by year and resource
curl "http://localhost:8080/api/search?year=2021,2022&resource=41"

# Paginated search with sorting
curl "http://localhost:8080/api/search?q=copper&size=20&from=0&sort_by=year&is_asc=false"
```

### Get Contract Details

```bash
# Get contract metadata
curl "http://localhost:8080/api/contracts/12345"

# Get full contract text
curl "http://localhost:8080/api/contracts/12345/text"

# Get contract annotations
curl "http://localhost:8080/api/contracts/12345/annotations"
```

### Export Results

```bash
# Export to DOCX
curl "http://localhost:8080/api/search?q=mining&download=true&type=docx" -o results.docx

# Export to TSV
curl "http://localhost:8080/api/search?year=2021&download=true&type=tsv" -o results.tsv
```

### Get Statistics

```bash
# Overall statistics
curl "http://localhost:8080/api/summary"

# Province-specific statistics
curl "http://localhost:8080/api/summary/year/province/1"
```

## Documentation

- **[Full Documentation](./DOCUMENTATION.md)** - Comprehensive guide covering architecture, modules, and data flow
- **[API Reference](./API_REFERENCE.md)** - Complete API endpoint documentation with examples
- **[Code Comments]** - Inline documentation throughout the codebase

## Project Structure

```
api-v2/
├── cmd/
│   └── front-service/          # Main API service entry point
├── internal/
│   ├── app_context/            # Elasticsearch client management
│   ├── correction/             # Data correction and translations
│   ├── document/               # Document generation (DOCX, TSV)
│   ├── queries/                # Elasticsearch query operations
│   ├── sql/                    # PostgreSQL operations
│   └── structs/                # Data structures and types
├── templates/                  # DOCX template files
├── DOCUMENTATION.md            # Comprehensive documentation
├── API_REFERENCE.md            # API endpoint reference
└── README.md                   # This file
```

## Key Technologies

- **[Gin](https://github.com/gin-gonic/gin)** - HTTP web framework
- **[Elasticsearch Go Client](https://github.com/elastic/go-elasticsearch)** - Search engine integration
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver
- **[UniOffice](https://github.com/unidoc/unioffice)** - DOCX generation
- **[godotenv](https://github.com/joho/godotenv)** - Environment configuration

## Development

### Building

```bash
# Build for current platform
go build -o build/front-service cmd/front-service/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o build/front-service-linux cmd/front-service/main.go
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test ./internal/queries -v
```

### Code Generation

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/search` | Search contracts with filters |
| GET | `/api/contracts/:id` | Get contract metadata |
| GET | `/api/contracts/:id/text` | Get contract full text |
| GET | `/api/contracts/:id/annotations` | Get contract annotations |
| GET | `/api/contracts-latest` | Get recent contracts |
| GET | `/api/metadata/:id` | Get contract preview data |
| GET | `/api/summary` | Get aggregated statistics |
| GET | `/api/summary/year/province/:id` | Get province statistics |
| GET | `/api/provinces` | Get provinces/districts |
| GET | `/api/provinces/all-units` | Get all administrative units |
| GET | `/api/page/:id` | Get static page content |
| GET | `/api/law/:id` | Get law content |
| GET | `/api/contracts/download/:id/:type` | Download contract file |
| GET | `/storage/*filepath` | Serve static files |
| POST | `/api/correction/resources` | Update resource values |
| POST | `/api/correction/contract_types` | Update contract types |
| POST | `/api/correction/document_types` | Update document types |

## Data Models

### Contract Metadata
- Contract name, type, date
- Company information
- Government entities
- Resources (minerals)
- Geographic location (province, district)
- Open Contracting ID (OCID)

### Annotations
- Text quotes and positions
- Categories and classifications
- Page numbers
- Geometric bounds on pages

### Search Results
- Matching documents
- Highlighted text snippets
- Relevance scores
- Aggregation facets

## Translation Mappings

The system includes comprehensive bilingual mappings:

- **119 Resource Types** (Minerals and materials)
- **30+ Contract Types** (Concession, Joint Venture, etc.)
- **20+ Document Types** (Contracts, EIA, Plans, etc.)

All mappings support bidirectional translation between Mongolian and English.

## Performance Considerations

- Connection pooling for Elasticsearch and PostgreSQL
- Lazy client initialization
- Temporary file cleanup after exports
- Configurable pagination limits
- Aggregation bucket size limits

## Security

⚠️ **Important Security Notes:**

- Currently no authentication implemented
- No rate limiting in place
- CORS configured for specific frontend URL
- Consider adding API keys for production
- File access should be restricted

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[Specify your license here]

## Support

For questions or issues:
- Create an issue on GitHub
- Contact: [your-email@example.com]

## Acknowledgments

- Iltodgeree Project Team
- Open Contracting Partnership
- MongoDB for initial data storage design inspiration
