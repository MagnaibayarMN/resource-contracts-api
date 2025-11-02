// Package queries provides database query operations for mining contracts.
package queries

import (
	"context"
	"encoding/json"
	"fmt"
	appcontext "iltodgeree/api/internal/app_context"
	"iltodgeree/api/internal/structs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

// GetLatestContracts retrieves the most recently created contracts from Elasticsearch.
// Results are sorted by creation date in descending order.
//
// Parameters:
//   - size: Maximum number of contracts to retrieve
//
// Returns:
//   - *structs.ResultWithHits: Search results containing contract metadata
//   - error: Error if the query fails
func GetLatestContracts(size int) (*structs.ResultWithHits, error) {
	documentType := "metadata"

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithDocumentType(documentType),
		client.Search.WithSize(size),
		client.Search.WithBody(strings.NewReader(`{
			"sort": [
					"_score",
					{
							"created_at": {
									"order": "desc",
									"unmapped_type": "date"
							}
					}
			],
			"query": {
					"match_all": {}
			}
		}`)),
	)

	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error in Elasticsearch response: %s", res.String())
	}

	var result structs.ResultWithHits

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// GetContractMaster retrieves a complete contract document including full text.
// This fetches from the 'master' document type which contains all contract data.
//
// Parameters:
//   - id: The unique contract identifier
//
// Returns:
//   - *elastic.GetResult: Complete contract document
//   - *error: Error if the document is not found
func GetContractMaster(id string) (*elastic.GetResult, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}
	index := os.Getenv("ELASTICSEARCH_SECONDARY")
	docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")

	result, err := client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Error fetching document: %s", err)
	}

	if !result.Found {
		err := fmt.Errorf("error executing search: %v", err)
		return nil, &err
	}

	return result, &err
}

// GetMetadata retrieves contract metadata and sanitized text for preview/SEO.
// Returns a simplified view with title and cleaned description text.
//
// Parameters:
//   - id: The unique contract identifier
//
// Returns:
//   - map[string]interface{}: Map containing 'title' and 'description' fields
//   - *error: Error if the document is not found
func GetMetadata(id string) (map[string]interface{}, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	// Define index, type, and document ID
	index := os.Getenv("ELASTICSEARCH_SECONDARY")
	docType := os.Getenv("ELASTICSEARCH_DOC_MASTER")
	// Get the document by ID
	result, err := client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(context.Background())

	if err != nil {
		err := fmt.Errorf("error fetching document: %v", err)
		return nil, &err
	}

	if !result.Found {
		err := fmt.Errorf("error executing search: %v", err)
		return nil, &err
	}

	var contract map[string]interface{}
	err = json.Unmarshal(*result.Source, &contract)

	if err != nil {
		log.Fatalf("Error fetching document: %s", err)
	}

	metadata := contract["metadata"].(map[string]interface{})
	title := metadata["contract_name"].(string)
	description := contract["pdf_text_string"].(string)

	re := regexp.MustCompile(`\s+`) // \s matches any whitespace character

	return map[string]interface{}{"title": title, "description": re.ReplaceAllString(description, " ")}, &err
}

// GetContract retrieves contract metadata without full text.
// This fetches from the 'metadata' document type for faster queries.
//
// Parameters:
//   - id: The unique contract identifier
//
// Returns:
//   - *elastic.GetResult: Contract metadata document
//   - *error: Error if the document is not found
func GetContract(id string) (*elastic.GetResult, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	// Define index, type, and document ID
	index := os.Getenv("ELASTICSEARCH_SECONDARY")
	docType := os.Getenv("ELASTICSEARCH_DOC_METADATA")
	// Get the document by ID
	result, err := client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(context.Background())

	if err != nil {
		return nil, &err
	}

	// Print the document
	if !result.Found {
		err := fmt.Errorf("error executing search: %v", err)
		return nil, &err
	}

	return result, &err
}

// GetAnnotationByContract retrieves all annotations associated with a contract.
// Annotations are sorted by ID in ascending order.
//
// Parameters:
//   - id: The contract ID to retrieve annotations for
//
// Returns:
//   - *elastic.SearchResult: Search results containing annotations
//   - *error: Error if the query fails
func GetAnnotationByContract(id string) (*elastic.SearchResult, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	query := elastic.NewTermQuery("contract_id", ID)

	index := os.Getenv("ELASTICSEARCH_SECONDARY")
	docType := "annotations"
	result, err := client.Search().
		Index(index).
		Type(docType).
		Query(query).
		Size(10000).
		From(0).
		Sort("id.keyword", true).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Error fetching document: %s", err)
	}

	if err != nil {
		if elasticErr, ok := err.(*elastic.Error); ok {
			fmt.Printf("Elastic error details: %v\n", elasticErr.Details)
		}
		log.Fatalf("Error executing search: %s", err)
	}

	return result, &err
}
