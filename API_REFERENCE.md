# Iltodgeree API Reference

## Base URL
```
http://localhost:8080
```

## Table of Contents
1. [Search Operations](#search-operations)
2. [Contract Operations](#contract-operations)
3. [Annotation Operations](#annotation-operations)
4. [Aggregation Operations](#aggregation-operations)
5. [Administrative Operations](#administrative-operations)
6. [Export Operations](#export-operations)
7. [Data Correction Operations](#data-correction-operations)

---

## Search Operations

### Search Contracts

**Endpoint:** `GET /api/search`

**Description:** Performs a comprehensive search across contracts with multiple filter options, full-text search, highlighting, and pagination.

**Query Parameters:**

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `q` | string | No | Full-text search query | `gold mining` |
| `year` | string | No | Comma-separated years | `2020,2021,2022` |
| `contract_type` | string | No | Contract type (English) | `Concession Agreement` |
| `resource` | string | No | Comma-separated resource types | `41,30` (IDs) |
| `company` | string | No | Company name | `Mining Corp` |
| `government` | string | No | Government entity | `Ministry` |
| `document_type` | string | No | Document type (English) | `Contract` |
| `province` | string | No | Province ID | `1` |
| `district` | string | No | Comma-separated district IDs | `101,102` |
| `annotation_category` | string | No | Annotation category | `Environmental` |
| `annotated` | boolean | No | Only annotated contracts | `true` |
| `size` | integer | No | Results per page | `20` |
| `from` | integer | No | Pagination offset | `0` |
| `sort_by` | string | No | Sort field | `year`, `country`, `contract_name`, `resource`, `contract_type` |
| `is_asc` | boolean | No | Sort ascending | `true` |
| `download` | string | No | Export flag | `true` |
| `type` | string | No | Export format | `docx`, `tsv` |

**Response Example:**

```json
{
  "took": 25,
  "hits": {
    "total": 150,
    "hits": [
      {
        "_id": "12345",
        "_score": 2.5,
        "_source": {
          "metadata": {
            "contract_name": "Gold Mining Agreement",
            "signature_year": 2021,
            "signature_date": "2021-05-15",
            "contract_type": "Concession Agreement",
            "resource": ["алт"],
            "company_name": ["ABC Mining Corp"],
            "government_entity": [
              {"entity": "Ministry of Mining"}
            ],
            "provinces": [
              {"province": "1", "district": "101"}
            ]
          }
        },
        "highlight": {
          "pdf_text_string": [
            "The <strong>gold mining</strong> operations shall..."
          ]
        }
      }
    ]
  }
}
```

---

## Contract Operations

### Get Contract Metadata

**Endpoint:** `GET /api/contracts/:id`

**Description:** Retrieves contract metadata without full text content.

**URL Parameters:**
- `id` - Contract unique identifier

**Response Example:**

```json
{
  "_id": "12345",
  "_source": {
    "metadata": {
      "contract_name": "Gold Mining Agreement",
      "signature_year": 2021,
      "open_contracting_id": "MN-GOV-12345",
      "contract_type": "Concession Agreement",
      "resource": ["алт"],
      "company_name": ["ABC Mining Corp"],
      "file_url": "/storage/12345/contract.pdf"
    }
  }
}
```

### Get Contract Full Text

**Endpoint:** `GET /api/contracts/:id/text`

**Description:** Retrieves the full text content of a contract.

**URL Parameters:**
- `id` - Contract unique identifier

**Response Example:**

```json
{
  "text": "MINING CONCESSION AGREEMENT\n\nThis agreement is made on...\n\n[Full contract text]"
}
```

### Get Contract Metadata for SEO

**Endpoint:** `GET /api/metadata/:id`

**Description:** Retrieves simplified metadata with sanitized text for previews and SEO.

**URL Parameters:**
- `id` - Contract unique identifier

**Response Example:**

```json
{
  "title": "Gold Mining Agreement - ABC Mining Corp",
  "description": "This agreement establishes the terms for gold mining operations in..."
}
```

### Get Latest Contracts

**Endpoint:** `GET /api/contracts-latest`

**Description:** Retrieves the 20 most recently created contracts.

**Response Example:**

```json
{
  "hits": {
    "total": 1234,
    "hits": [
      {
        "_id": "12345",
        "_source": {
          "metadata": {
            "contract_name": "Recent Agreement",
            "created_at": "2024-01-15T10:30:00Z"
          }
        }
      }
    ]
  }
}
```

---

## Annotation Operations

### Get Contract Annotations

**Endpoint:** `GET /api/contracts/:id/annotations`

**Description:** Retrieves all annotations for a specific contract, sorted by ID.

**URL Parameters:**
- `id` - Contract ID

**Response Example:**

```json
{
  "took": 5,
  "hits": {
    "total": 25,
    "hits": [
      {
        "_source": {
          "id": "ann-001",
          "contract_id": "12345",
          "annotation_id": 1,
          "quote": "environmental impact assessment",
          "text": "EIA requirements must be met",
          "category": "Environmental",
          "category_key": "env_impact",
          "page_no": 5,
          "ranges": "[{\"start\":100,\"end\":150}]"
        }
      }
    ]
  }
}
```

---

## Aggregation Operations

### Get Summary Statistics

**Endpoint:** `GET /api/summary`

**Description:** Retrieves aggregated statistics across all contracts.

**Response Example:**

```json
{
  "count": 1234,
  "aggs": {
    "year_summary": {
      "buckets": [
        {"key": "2021", "doc_count": 150},
        {"key": "2020", "doc_count": 120}
      ]
    },
    "resource_summary": {
      "buckets": [
        {"key": "алт", "doc_count": 250},
        {"key": "зэс", "doc_count": 180}
      ]
    },
    "contract_type_summary": {
      "buckets": [
        {"key": "Concession Agreement", "doc_count": 300}
      ]
    },
    "provinces_summary": {
      "buckets": [
        {"key": "1", "doc_count": 200},
        {"key": "2", "doc_count": 150}
      ]
    }
  }
}
```

### Get Year Statistics by Province

**Endpoint:** `GET /api/summary/year/province/:id`

**Description:** Retrieves year-based contract counts for a specific province.

**URL Parameters:**
- `id` - Province ID (integer)

**Response Example:**

```json
[
  {"x": "2021", "y": 45},
  {"x": "2020", "y": 38},
  {"x": "2019", "y": 52}
]
```

---

## Administrative Operations

### Get Provinces

**Endpoint:** `GET /api/provinces`

**Description:** Retrieves provinces or districts based on parameters.

**Query Parameters:**
- `province_id` - Optional. If provided, returns districts within that province.

**Response Examples:**

**All Provinces:**
```json
[
  {
    "id": 1,
    "name": "Улаанбаатар",
    "type": 1,
    "parentId": 0
  },
  {
    "id": 2,
    "name": "Архангай",
    "type": 1,
    "parentId": 0
  }
]
```

**Districts in Province:**
```json
[
  {
    "id": 101,
    "name": "Баянгол",
    "type": 2,
    "parentId": 1
  },
  {
    "id": 102,
    "name": "Сүхбаатар",
    "type": 2,
    "parentId": 1
  }
]
```

### Get All Administrative Units

**Endpoint:** `GET /api/provinces/all-units`

**Description:** Retrieves all provinces and districts as a flat map.

**Response Example:**

```json
{
  "1": "Улаанбаатар",
  "101": "Баянгол",
  "102": "Сүхбаатар",
  "2": "Архангай",
  "201": "Эрдэнэбулган"
}
```

### Get Page Content

**Endpoint:** `GET /api/page/:id`

**Description:** Retrieves static page content.

**URL Parameters:**
- `id` - Page identifier

**Query Parameters:**
- `locale` - Language code (e.g., `en`, `mn`)

**Response Example:**

```json
{
  "id": "about",
  "title": "About Us",
  "content": "<p>Page content...</p>",
  "locale": "en"
}
```

### Get Law Content

**Endpoint:** `GET /api/law/:id`

**Description:** Retrieves law/regulation content.

**URL Parameters:**
- `id` - Law identifier

**Query Parameters:**
- `locale` - Language code

**Response Example:**

```json
{
  "id": "mining-law-2006",
  "title": "Law on Minerals",
  "content": "Full law text...",
  "locale": "en"
}
```

---

## Export Operations

### Download Contract File

**Endpoint:** `GET /api/contracts/download/:id/:type`

**Description:** Downloads contract file in specified format.

**URL Parameters:**
- `id` - Contract ID
- `type` - File type (`pdf` or `docx`)

**Response:** File download (binary)

**Example:**
```
GET /api/contracts/download/12345/pdf
→ Downloads the PDF file
```

### Export Search Results

**Endpoint:** `GET /api/search?download=true&type=docx`

**Description:** Exports search results to DOCX or TSV format.

**Query Parameters:**
- All search parameters (see Search Operations)
- `download=true` - Enable export
- `type` - Export format (`docx` or `tsv`)

**Response:** File download

**DOCX Export:**
- Combines multiple contracts into a single Word document
- Each contract is numbered and formatted
- Includes contract title and full text

**TSV Export Columns:**
1. # (Number)
2. Гэрээний нэр (Contract Name)
3. Эрдсийн төрөл (Resource Type)
4. Гэрээний төрөл (Contract Type)
5. Гэрээ байгуулсан огноо (Signature Date)
6. Баримт бичгийн төрөл (Document Type)
7. Аймаг / Сум (Province / District)
8. Гэрээ байгуулсан төрийн байгууллага (Government Entity)
9. Компанийн нэр (Company Name)
10. Төслийн нэр (Project Name)
11. Гэрээний файл (Contract File URL)
12. OCID (Open Contracting ID)
13. Аннотацийн текст (Annotation Text)
14. Метадата текст (Metadata Text)

---

## Data Correction Operations

### Update Resources

**Endpoint:** `POST /api/correction/resources`

**Description:** Performs bulk updates on resource field values.

**Request Body:**

```json
[
  {
    "key": "алт (old spelling)",
    "value": "алт (new spelling)"
  },
  {
    "key": "old_resource_name",
    "value": "new_resource_name"
  }
]
```

**Response:**

```json
null
```

**Note:** Updates are applied to all matching documents in Elasticsearch.

### Update Contract Types

**Endpoint:** `POST /api/correction/contract_types`

**Description:** Performs bulk updates on contract type field values.

**Request Body:**

```json
[
  {
    "key": "Old Contract Type",
    "value": "New Contract Type"
  }
]
```

### Update Document Types

**Endpoint:** `POST /api/correction/document_types`

**Description:** Performs bulk updates on document type field values.

**Request Body:**

```json
[
  {
    "key": "Old Document Type",
    "value": "New Document Type"
  }
]
```

---

## Static File Serving

### Serve Storage Files

**Endpoint:** `GET /storage/*filepath`

**Description:** Serves static files from the storage directory.

**Example:**
```
GET /storage/12345/contract.pdf
→ Returns the PDF file
```

**Headers:**
- `Access-Control-Allow-Origin`: Configured frontend URL
- `Access-Control-Allow-Methods`: GET, OPTIONS

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description",
  "panic": "Detailed error information (if panic occurred)"
}
```

**Common HTTP Status Codes:**
- `200 OK` - Successful request
- `400 Bad Request` - Invalid parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Rate Limiting

Currently not implemented. Consider adding rate limiting for production use.

## Authentication

Currently not implemented. All endpoints are publicly accessible.

## CORS Configuration

CORS is configured to allow requests from the frontend URL specified in `FRONT_END_URL` environment variable.
