// Package correction provides utilities for bulk updating Elasticsearch documents.
// It handles corrections for resources, contract types, and document types using
// Elasticsearch's update-by-query API with Painless scripts.
package correction

import (
	"context"
	"fmt"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

// ResourcesCorrection performs bulk updates on resource fields in Elasticsearch documents.
// It finds all documents with a specific resource value and replaces it with a new value.
//
// Parameters:
//   - index: The Elasticsearch index name
//   - docType: The document type within the index
//   - key: The current resource value to find and replace
//   - value: The new resource value to set
//
// The function uses a Painless script to iterate through resource arrays and update matches.
func ResourcesCorrection(index string, docType string, key string, value string) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	query := elastic.NewTermQuery("metadata.resource.keyword", key)
	scr := `
			for (int i = 0; i < ctx._source.metadata.resource.length; i++) {
				if (ctx._source.metadata.resource[i] == '` + key + `') {
						ctx._source.metadata.resource[i] = '` + value + `';
				}
			}
    `
	script := elastic.NewScript(scr).Lang("painless")

	updateResult, err := client.UpdateByQuery().
		Index(index).
		Type(docType).
		Query(query).
		Script(script).
		Do(ctx)

	if err != nil {
		log.Fatalf("Error performing update: %s", err)
	}

	fmt.Printf("Updated %d documents\n", updateResult.Updated)
}

// ContractTypesCorrection performs bulk updates on contract type fields in Elasticsearch.
// It updates the contract_type_raw field for all matching documents.
//
// Parameters:
//   - index: The Elasticsearch index name
//   - docType: The document type within the index
//   - k: The current contract type value to find
//   - v: The new contract type value to set
func ContractTypesCorrection(index string, docType string, k string, v string) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	scr := "ctx._source.metadata.contract_type_raw = '" + v + "'"
	query := elastic.NewTermQuery("metadata.contract_type_raw.keyword", k)
	script := elastic.NewScript(scr).Lang("painless")

	updateResult, err := client.UpdateByQuery().
		Index(index).
		Type(docType).
		Query(query).
		Script(script).
		Do(ctx)

	if err != nil {
		log.Fatalf("Error performing update: %s", err)
	}

	fmt.Printf("Updated %d documents\n", updateResult.Updated)
}

// DocumentTypesCorrection performs bulk updates on document type fields in Elasticsearch.
// It updates the document_type field for all matching documents.
//
// Parameters:
//   - index: The Elasticsearch index name
//   - docType: The document type within the index
//   - k: The current document type value to find
//   - v: The new document type value to set
func DocumentTypesCorrection(index string, docType string, k string, v string) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	scr := "ctx._source.metadata.document_type = '" + v + "'"
	query := elastic.NewTermQuery("metadata.document_type.keyword", k)
	script := elastic.NewScript(scr).Lang("painless")

	updateResult, err := client.UpdateByQuery().
		Index(index).
		Type(docType).
		Query(query).
		Script(script).
		Do(ctx)

	if err != nil {
		log.Fatalf("Error performing update: %s", err)
	}

	fmt.Printf("Updated %d documents\n", updateResult.Updated)
}
