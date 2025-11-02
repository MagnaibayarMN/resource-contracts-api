# Iltodgeree Mining Contracts API Documentation

## Overview

The Iltodgeree Mining Contracts API is a Go-based RESTful service that provides access to a comprehensive database of mining contracts, their metadata, annotations, and related documents. The system is built on Elasticsearch for full-text search capabilities and PostgreSQL for administrative data.

## Architecture

### Technology Stack

- **Language**: Go 1.23.4
- **Web Framework**: Gin (HTTP router and middleware)
- **Search Engine**: Elasticsearch v5
- **Database**: PostgreSQL (via pgx driver)
- **Document Processing**: UniOffice/UniPDF for DOCX generation

### Project Structure

```
api-v2/
├── cmd/
│   ├── front-service/          # Main API service
│   └── indexing-service/       # Document indexing service (placeholder)
├── internal/
│   ├── app_context/            # Elasticsearch client management
│   ├── common/                 # Common utilities
│   ├── correction/             # Data correction and translation mappings
│   ├── document/               # Document generation (DOCX, TSV)
│   ├── indexing/               # Contract state management
│   ├── queries/                # Elasticsearch query operations
│   ├── sql/                    # PostgreSQL operations
│   └── structs/                # Data structures
├── templates/                  # DOCX templates
└── result/                     # Temporary export files
```

## Core Modules

### 1. Front Service (`cmd/front-service/main.go`)

The main HTTP service that exposes REST endpoints for:
- Contract search and retrieval
- Metadata access
- Annotation queries
- Document export (DOCX, TSV)
- Data correction operations
- Administrative unit lookups

### 2. Application Context (`internal/app_context`)

Manages Elasticsearch client connections with lazy initialization:
- **v5 Client**: For legacy Elasticsearch 5.x compatibility
- **v7+ Client**: For newer Elasticsearch versions
- Query builder utilities

### 3. Queries Module (`internal/queries`)

Provides Elasticsearch operations:
- **search2.go**: Advanced search with filters, highlighting, pagination
- **contracts.go**: Contract retrieval by ID
- **annotations.go**: Annotation queries
- **aggs.go**: Aggregations for statistics and faceted search
- **download.go**: File download handlers

### 4. Document Module (`internal/document`)

Document generation and export:
- **process.go**: DOCX and TSV generation from search results
- **docx.go**: Word document XML utilities
- **zip.go**: ZIP archive operations for DOCX files

### 5. Correction Module (`internal/correction`)

Data correction and bilingual mappings:
- **metadata.go**: Bulk update operations via Elasticsearch
- **types.go**: Translation maps for resources, contract types, document types

### 6. SQL Module (`internal/sql`)

PostgreSQL operations for:
- **provinces.go**: Province and district data
- **page.go**: Static page content

## API Endpoints

### Search & Retrieval

#### `GET /api/search`
Full-text search with multiple filters.

**Query Parameters:**
- `q` - Search query string
- `year` - Comma-separated years
- `contract_type` - Contract type filter
- `resource` - Resource type filter
- `company` - Company name filter
- `government` - Government entity filter
- `document_type` - Document type filter
- `province` - Province ID
- `district` - District ID
- `annotation_category` - Annotation category
- `annotated` - Boolean, has annotations
- `size` - Results per page (default: 10)
- `from` - Pagination offset (default: 0)
- `sort_by` - Sort field (country, year, contract_name, resource, contract_type)
- `is_asc` - Sort order (true/false)
- `download` - Export results
- `type` - Export format (docx, tsv)

**Response:** JSON with search results, highlights, and metadata

#### `GET /api/contracts/:id`
Get contract metadata by ID.

**Response:** Contract metadata object

#### `GET /api/contracts/:id/text`
Get contract full text content.

**Response:** 
```json
{
  "text": "Full contract text..."
}
```

#### `GET /api/metadata/:id`
Get contract title and description for SEO/preview.

**Response:**
```json
{
  "title": "Contract Name",
  "description": "Sanitized description text..."
}
```

#### `GET /api/contracts-latest`
Get 20 most recently created contracts.

**Response:** Array of contract objects

### Annotations

#### `GET /api/contracts/:id/annotations`
Get all annotations for a specific contract.

**Response:** Array of annotation objects with categories, quotes, and page references

### Aggregations & Statistics

#### `GET /api/summary`
Get aggregated statistics across all contracts.

**Response:**
```json
{
  "aggs": {
    "year_summary": {...},
    "resource_summary": {...},
    "contract_type_summary": {...},
    "country_summary": {...},
    "provinces_summary": {...}
  },
  "count": 1234
}
```

#### `GET /api/summary/year/province/:id`
Get year-based statistics filtered by province.

**Parameters:**
- `:id` - Province ID

**Response:** Array of year/count pairs

### Administrative Data

#### `GET /api/provinces`
Get all provinces or districts within a province.

**Query Parameters:**
- `province_id` - Optional, returns districts if provided

**Response:** Array of province/district objects

#### `GET /api/provinces/all-units`
Get all provinces and districts as a flat map.

**Response:** Map of ID to name

### Static Content

#### `GET /api/page/:id`
Get static page content.

**Query Parameters:**
- `locale` - Language code

**Response:** Page content object

#### `GET /api/law/:id`
Get law/regulation content.

**Query Parameters:**
- `locale` - Language code

**Response:** Law content object

### Document Export

#### `GET /api/contracts/download/:id/:type`
Download contract file.

**Parameters:**
- `:id` - Contract ID
- `:type` - File type (pdf, docx)

**Response:** File download

### Data Correction (Admin)

#### `POST /api/correction/resources`
Bulk update resource values.

**Request Body:**
```json
[
  {
    "key": "current_value",
    "value": "new_value"
  }
]
```

#### `POST /api/correction/contract_types`
Bulk update contract type values.

#### `POST /api/correction/document_types`
Bulk update document type values.

### Static Files

#### `GET /storage/*filepath`
Serve static files from storage directory.

## Data Structures

### SearchParams
Encapsulates search parameters with support for:
- Full-text queries
- Multiple filters (years, types, resources)
- Geographic filters (province, district)
- Pagination and sorting
- Annotation filters

### Annotation
Represents document annotations with:
- Contract references
- Text quotes and positions
- Categories and classifications
- Page numbers and geometric bounds

### Province/District
Administrative units with:
- ID and name
- Parent relationships
- Geographic data

## Configuration

Required environment variables:

```bash
# Elasticsearch
ELASTICSEARCH_HOST=http://localhost:9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=password
ELASTICSEARCH_SECONDARY=iltodgeree_v2.2
ELASTICSEARCH_DOC_MASTER=master
ELASTICSEARCH_DOC_METADATA=metadata

# PostgreSQL
PGSQL_URL=postgresql://user:pass@localhost:5432/iltodgeree

# File Paths
DOCUMENT_PATH=/path/to/documents
TEMPLATE_PATH=/path/to/templates
STORAGE_PATH=/path/to/storage
PUBLIC_URL=https://api.example.com

# Frontend
FRONT_END_URL=http://localhost:3000
```

## Development

### Running the Service

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/front-service/main.go
```

### Building

```bash
# Build for production
go build -o build/front-service cmd/front-service/main.go
```

### Using Docker

```bash
# Build image
docker build -t iltodgeree-api .

# Run container
docker run -p 8080:8080 --env-file .env iltodgeree-api
```

## Data Flow

1. **Search Request** → Search parameters parsed → Elasticsearch query built → Results returned with highlights
2. **Contract Retrieval** → ID lookup → Elasticsearch GET → Document returned
3. **Export** → Search executed → Results processed → DOCX/TSV generated → File download
4. **Annotation Query** → Contract ID → Elasticsearch filter → Annotations returned

## Translation & Localization

The system supports bilingual (Mongolian/English) operations:
- Contract types, document types, and resources have translation maps
- Static content supports multiple locales
- API responses include localized field names

## Performance Considerations

- Elasticsearch connection pooling with lazy initialization
- PostgreSQL connection pool (min 2 connections)
- Temporary files cleaned up after export
- Pagination for large result sets
- Aggregation size limits (10,000 buckets)

## Error Handling

- Custom panic recovery middleware captures errors
- Structured error responses with HTTP status codes
- Detailed logging for debugging

## Security Notes

- CORS configured for specific frontend URLs
- File serving restricted to storage directory
- SQL injection protection via parameterized queries
- Elasticsearch authentication via basic auth

## Future Enhancements

- Implement indexing service for document ingestion
- Add authentication/authorization layer
- Implement rate limiting
- Add caching layer (Redis)
- Enhanced monitoring and metrics
- GraphQL API option
