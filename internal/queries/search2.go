// Package queries provides search and query functionality for mining contracts in Elasticsearch.
// It handles full-text search, filtering, aggregations, and result retrieval.
package queries

import (
	"context"
	"encoding/json"
	"fmt"
	appcontext "iltodgeree/api/internal/app_context"
	"iltodgeree/api/internal/correction"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

var _defaultSize = 10
var _defaultFrom = 0

// SearchParams encapsulates all possible search parameters for querying contracts.
// It supports full-text search, multiple filters, pagination, and sorting options.
type SearchParams struct {
	q                  string
	years              []interface{}
	resources          []interface{}
	companies          string
	governments        string
	contractTypes      []interface{}
	documentTypes      []interface{}
	annotationCategory []interface{}
	annotated          bool
	province           string
	district           []interface{}
	size               int
	from               int
	sortBy             *string
	order              *bool
}

// NewSearchParams creates a new SearchParams instance with the provided filter values.
// It parses comma-separated strings into arrays for multi-value filters.
//
// Parameters:
//   - q: Full-text search query string
//   - years: Comma-separated list of years
//   - contractTypes: Comma-separated list of contract types
//   - resources: Comma-separated list of resource types
//   - companies: Company name filter
//   - governments: Government entity filter
//   - documentTypes: Comma-separated list of document types
//
// Returns:
//   - *SearchParams: Configured search parameters object
func NewSearchParams(
	q string,
	years string,
	contractTypes string,
	resources string,
	companies string,
	governments string,
	documentTypes string,
) *SearchParams {

	var _years []interface{}
	var _contractTypes []interface{}
	var _resources []interface{}
	var _companies string
	var _governments string
	var _documentTypes []interface{}

	if years != "" {
		for _, part := range strings.Split(years, ",") {
			num, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			_years = append(_years, num)
		}
	}

	if contractTypes != "" {
		for _, part := range strings.Split(contractTypes, ",") {
			_contractTypes = append(_contractTypes, part)
		}
	}

	if resources != "" {
		for _, part := range strings.Split(resources, ",") {
			res := part
			_resources = append(_resources, res)
		}
	}

	if companies != "" {
		_companies = companies
	}

	if governments != "" {
		_governments = governments
	}

	if documentTypes != "" {
		for _, part := range strings.Split(documentTypes, ",") {
			_documentTypes = append(_documentTypes, part)
		}
	}

	return &SearchParams{
		q:             q,
		resources:     _resources,
		years:         _years,
		contractTypes: _contractTypes,
		companies:     _companies,
		governments:   _governments,
		documentTypes: _documentTypes,
		sortBy:        new(string),
		order:         new(bool),
	}
}

func (s *SearchParams) SetProvince(province string) {
	s.province = province
}

func (s *SearchParams) SetDistrict(district string) {
	var _districts []interface{}
	if district != "" {
		for _, part := range strings.Split(district, ",") {
			num, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			_districts = append(_districts, num)
		}
	}
	s.district = _districts
}

func (s *SearchParams) SetAnnotationCategories(categories string) {
	if categories != "" {
		for _, category := range strings.Split(categories, ",") {
			s.annotationCategory = append(s.annotationCategory, category)
		}
	}
}

func (s *SearchParams) SetAnnotated(annotated bool) {
	s.annotated = annotated
}

func (s *SearchParams) SetSortBy(sortBy string) {
	if sortBy != "" {
		if sortBy == "country" {
			*s.sortBy = "metadata.country_name.keyword"
		} else if sortBy == "year" {
			// *s.sortBy = "metadata.signature_year.keyword"
			*s.sortBy = "metadata.signature_date"
		} else if sortBy == "contract_name" {
			*s.sortBy = "metadata.contract_name.keyword"
		} else if sortBy == "resource" {
			*s.sortBy = "metadata.resource_raw.keyword"
		} else if sortBy == "contract_type" {
			*s.sortBy = "metadata.contract_type.keyword"
		}
	} else {
		*s.sortBy = "metadata.signature_date"
		*s.order = false
	}
}

func (s *SearchParams) SetOrder(order string) {
	if order != "" {
		_order, err := strconv.ParseBool(order)
		if err != nil {
			log.Println("boolean утгыг хөрвүүлж чадсангүй.")
		}
		*s.order = _order
	}
}

func (s *SearchParams) SetSize(size string) {
	if size == "" {
		s.size = _defaultSize
	} else {
		_size, err := strconv.Atoi(size)
		if err != nil {
		}
		s.size = _size
	}
}

func (s *SearchParams) SetFrom(from string) {
	if from == "" {
		s.from = _defaultFrom
	} else {
		_from, err := strconv.Atoi(from)
		if err != nil {
		}
		s.from = _from
	}
}

// SearchV2 executes a comprehensive search query against Elasticsearch.
// It builds a bool query with filters, performs full-text search if specified,
// applies highlights, sorting, and pagination.
//
// Parameters:
//   - params: SearchParams object containing all search criteria
//
// Returns:
//   - *elastic.SearchResult: Search results from Elasticsearch
//   - *error: Error if the search fails
func SearchV2(params *SearchParams) (*elastic.SearchResult, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	boolQuery := elastic.NewBoolQuery()

	filters := []elastic.Query{}
	phrases := []elastic.Query{}

	if len(params.years) != 0 {
		filters = append(filters, elastic.NewTermsQuery("metadata.signature_year", params.years...))
	}

	if len(params.resources) != 0 {
		filters = append(filters, elastic.NewTermsQuery("metadata.resource", params.resources...))
	}

	if params.province != "" {
		filters = append(filters, elastic.NewTermsQuery("metadata.provinces.province", params.province))
	}

	if len(params.district) > 0 {
		filters = append(filters, elastic.NewTermsQuery("metadata.provinces.district", params.district...))
	}

	if len(params.documentTypes) > 0 {
		for _, documentType := range params.documentTypes {
			filters = append(filters, elastic.NewTermsQuery("metadata.document_type.keyword", correction.DocumentTypesReverse[documentType.(string)]))
		}
	}

	if len(params.contractTypes) > 0 {
		for _, t := range params.contractTypes {
			filters = append(filters, elastic.NewTermsQuery("metadata.contract_type.keyword", correction.ContractTypesReverse[t.(string)]))
		}
	}

	if len(params.annotationCategory) > 0 {
		for _, t := range params.annotationCategory {
			phrases = append(phrases, elastic.NewMatchPhraseQuery("annotations_category", t))
		}
	}

	if len(params.governments) > 0 {
		filters = append(filters, elastic.NewTermsQuery("metadata.government_entity.entity.keyword", params.governments))
		// phrases = append(phrases, elastic.NewMatchQuery("metadata.government_entity.entity", params.governments).Operator("and"))
	}

	if len(params.companies) > 0 {
		filters = append(filters, elastic.NewTermsQuery("metadata.company_name.keyword", params.companies))
		// phrases = append(phrases, elastic.NewMatchQuery("metadata.company_name", params.companies).Operator("and"))
	}

	fields := []string{
		"metadata.contract_name",
		"metadata.project_title",
		"metadata.open_contracting_id",
		"metadata.country_code",
		"metadata.country_name",
		"metadata.resource",
		"metadata.resource_raw",
		"metadata.language",
		"metadata.company_name",
		"metadata.type_of_contract",
		// "metadata.corporate_grouping",
		"metadata.show_pdf_text",
		"metadata.category",
		"metadata_string",
		"pdf_text_string",
	}

	highlights := []string{
		"pdf_text_string",
		"metadata_string",
	}

	var ftsQuery *elastic.SimpleQueryStringQuery
	highlight := elastic.NewHighlight().PreTags("<strong>").PostTags("</strong>")

	if params.q != "" {
		ftsQuery = elastic.NewSimpleQueryStringQuery(params.q)
		for _, field := range fields {
			ftsQuery = ftsQuery.Field(field)
		}

		ftsQuery = ftsQuery.DefaultOperator("AND")

		for _, h := range highlights {
			highlight = highlight.Field(h).FragmentSize(50).NumOfFragments(2)
		}
	}

	// considered as unnessasary
	// boolQuery = boolQuery.MustNot(elastic.NewMatchQuery("pdf_text_string", "NotEnoughCredits"))

	if len(phrases) > 0 {
		boolQuery = boolQuery.Should(phrases...)
	}

	if len(filters) > 0 {
		boolQuery = boolQuery.Filter(filters...)
	}

	src, err := boolQuery.Source()
	if err != nil {
		log.Fatalf("Error getting query source: %v", err)
	}

	jsonData, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling query: %v", err)
	}

	fmt.Println("Generated Query:")
	fmt.Println(string(jsonData))

	index := os.Getenv("ELASTICSEARCH_SECONDARY")
	docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")

	q := client.Search().
		Index(index).
		Type(docType)

	if ftsQuery != nil {
		boolQuery = boolQuery.Must(ftsQuery)
	}

	// if ftsQuery != nil {
	// 	q = q.Query(ftsQuery)
	// } else {
	// 	q = q.Query(boolQuery)
	// }

	q = q.Query(boolQuery)

	if highlight != nil {
		q = q.Highlight(highlight)
	}

	if params.sortBy != nil {
		q = q.Sort(*params.sortBy, *params.order)
	}

	q = q.From(params.from)
	q = q.Size(params.size)

	result, err := q.Pretty(true).Do(context.Background())

	if err != nil {
		if elasticErr, ok := err.(*elastic.Error); ok {
			fmt.Printf("Elastic error details: %v\n", elasticErr.Details)
		}
		log.Fatalf("Error executing search: %s", err)
	}
	return result, &err
}
