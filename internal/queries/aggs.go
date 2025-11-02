// Package queries provides aggregation operations for statistical analysis.
package queries

import (
	"context"
	appcontext "iltodgeree/api/internal/app_context"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

// AggResult holds aggregation results with document count.
type AggResult struct {
	aggs  elastic.Aggregations // Aggregation buckets
	count int64              // Total document count
}

// YearFilterAggregations computes contract counts by year for a specific province.
// Returns data formatted for charting (x: year, y: count).
//
// Parameters:
//   - provinceID: The province ID to filter by
//
// Returns:
//   - *[]map[string]interface{}: Array of year/count pairs
//   - error: Error if aggregation fails
func YearFilterAggregations(provinceID int) (*[]map[string]interface{}, error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	aggSize := 10000

	provinceFilter := elastic.NewTermQuery("metadata.provinces.province", provinceID)
	year := elastic.NewTermsAggregation().Field("metadata.signature_year.keyword").Size(aggSize)
	filteredAggregation := elastic.NewFilterAggregation().Filter(provinceFilter).SubAggregation("filtered_year", year)

	index := "iltodgeree_v2.2"
	docType := "master"

	result, err := client.Search().
		Index(index).
		Type(docType).
		Size(0).
		Aggregation("year_summary", filteredAggregation).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	agg, found := result.Aggregations.Filter("year_summary")
	if found {
		yearSummary, found := agg.Aggregations.Terms("filtered_year")

		if found {
			for _, bucket := range yearSummary.Buckets {
				data = append(data, map[string]interface{}{
					"x": bucket.Key,
					"y": int(bucket.DocCount),
				})
			}
		}
	}

	return &data, nil
}

// Aggregations computes comprehensive statistics across all contracts.
// Includes aggregations by year, resource, type, country, province, company, etc.
// Also computes resource distribution by year for trend analysis.
//
// Returns:
//   - *map[string]interface{}: Map containing 'aggs' (aggregations) and 'count' (total)
//   - error: Error if query fails
func Aggregations() (*map[string]interface{}, error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	aggSize := 10000

	// total := elastic.NewValueCountAggregation().Field("_id.keyword")
	year := elastic.NewTermsAggregation().Field("metadata.signature_year.keyword").Size(aggSize)
	resource := elastic.NewTermsAggregation().Field("metadata.resource.keyword").Size(aggSize)
	document := elastic.NewTermsAggregation().Field("metadata.document_type.keyword").Size(aggSize)
	contractType := elastic.NewTermsAggregation().Field("metadata.contract_type.keyword").Size(aggSize)
	country := elastic.NewTermsAggregation().Field("metadata.country_code.keyword").Size(aggSize)
	company := elastic.NewTermsAggregation().Field("metadata.company_name.keyword").Size(aggSize)
	government := elastic.NewTermsAggregation().Field("metadata.government_entity.entity.keyword").Size(aggSize)
	provinces := elastic.NewTermsAggregation().Field("metadata.provinces.province.keyword").Size(aggSize)
	districts := elastic.NewTermsAggregation().Field("metadata.provinces.district.keyword").Size(aggSize)
	annotationCategories := elastic.NewTermsAggregation().Field("annotations_category.keyword").Size(aggSize)

	index := "iltodgeree_v2.2"
	docType := "master"

	countService := client.Count().
		Index(index). // Your index
		Type(docType) // Your document type (for ES5, types still exist)

	ctx := context.Background()
	count, err := countService.Do(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.Search().
		Index(index).
		Type(docType).
		Size(0).
		Aggregation("year_summary", year).
		Aggregation("resource_summary", resource).
		Aggregation("document_summary", document).
		Aggregation("country_summary", country).
		Aggregation("provinces_summary", provinces).
		Aggregation("districts_summary", districts).
		Aggregation("contract_type_summary", contractType).
		Aggregation("government_summary", government).
		Aggregation("company_summary", company).
		Aggregation("annotations_summary", annotationCategories).
		Aggregation("resource_by_years_summary", elastic.NewTermsAggregation().
			Field("metadata.resource.keyword").
			Size(aggSize).
			SubAggregation("signature_years", elastic.NewTermsAggregation().
				Field("metadata.signature_year.keyword").
				Size(aggSize),
			)).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return &map[string]interface{}{
		"aggs":  result.Aggregations,
		"count": count,
	}, nil
}

// ResourceByYearsAggregation computes resource distribution over years.
// Provides nested aggregations: resources -> years within each resource.
//
// Returns:
//   - *elastic.Aggregations: Nested aggregation results
//   - *error: Error if query fails
func ResourceByYearsAggregation() (*elastic.Aggregations, *error) {
	client, err := appcontext.ElasticInstance.GetV5()
	if err != nil {
		panic(err)
	}

	aggSize := 10000

	index := "iltodgeree_v2"
	docType := "metadata"
	result, err := client.
		Search().
		Index(index).
		Type(docType).
		Size(0).
		Aggregation("resource_summary", elastic.NewTermsAggregation().
			Field("metadata.resource.keyword").
			Size(aggSize).
			SubAggregation("signature_years", elastic.NewTermsAggregation().
				Field("metadata.signature_year.keyword").
				Size(aggSize),
			)).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Error fetching document: %s", err)
		return nil, &err
	}

	return &result.Aggregations, &err
}
