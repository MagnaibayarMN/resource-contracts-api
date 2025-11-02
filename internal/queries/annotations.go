package queries

import (
	"encoding/json"
	"fmt"
	ctx "iltodgeree/api/internal/app_context"
	"iltodgeree/api/internal/structs"
	"strings"
)

var documentType = "annotations"

// GetAnnotationsCount
func GetAnnotationsCount() (int, error) {
	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithDocumentType(documentType),
	)

	if err != nil {
		return 0, fmt.Errorf("error executing search: %v", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error in Elasticsearch response: %s", res.String())
	}

	var result structs.Result

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Hits.Total, nil
}

func parseArguments(contractId interface{}, page *int) structs.Term {
	term := make(structs.Term)
	if value, ok := contractId.(string); ok {
		term["open_contracting_id"] = map[string]interface{}{
			"value": value,
		}
	}

	if value, ok := contractId.(int); ok {
		term["contract_id"] = map[string]interface{}{
			"value": value,
		}
	}

	if page != nil {
		term["page"] = map[string]interface{}{
			"value": page,
		}
	}

	return term
}

func GetAnnotationPages(contractId interface{}, page *int) (response *[]structs.Annotation, err error) {
	arguments := parseArguments(contractId, page)
	query, err := ctx.ElasticInstance.BuildQueryBoolMust(arguments)
	if err != nil {
		fmt.Println("Error marshaling to JSON query:", err)
		return nil, err
	}

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithDocumentType(documentType),
		client.Search.WithSize(defaultSize),
		client.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	var result structs.ResultWithHits

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error while decoding JSON: %v", err)
	}

	// TODO: implement https://gitlab.com/iltodgeree/elasticsearch-api/-/blob/master/app/Services/APIServices.php?ref_type=heads#L919
	for _, item := range result.Hits.Hits {
		*response = append(*response, structs.Annotation{
			ContractID:       item.Source["contract_id"].(string),
			OpenContractID:   item.Source["open_contracting_id"].(string),
			ID:               item.Source["id"].(int),
			AnnotationID:     item.Source["annotation_id"].(int),
			Quote:            item.Source["quote"].(string),
			Text:             item.Source["text"].(string),
			Category:         item.Source["cateogry"].(string),
			CategoryKey:      item.Source["category_key"].(string),
			ArticleReference: item.Source["article_reference"].(string),
			PageNo:           item.Source["page_no"].(int),
			Ranges:           item.Source["ranges"].(string),
			Cluster:          item.Source["cluster"].(string),
		})
	}

	defer res.Body.Close()

	return response, nil
}

// todo: group pages by annotation id
func GetAnnotationGroup(contractId int, page *int) (response *[]structs.AnnotationGroup, err error) {
	arguments := parseArguments(contractId, page)
	query, err := ctx.ElasticInstance.BuildQueryBoolMust(arguments)
	if err != nil {
		fmt.Println("Error marshaling to JSON query:", err)
		return nil, err
	}

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithDocumentType(documentType),
		client.Search.WithSize(defaultSize),
		client.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	var result structs.ResultWithHits

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error while decoding JSON: %v", err)
	}

	// TODO: implement https://gitlab.com/iltodgeree/elasticsearch-api/-/blob/master/app/Services/APIServices.php?ref_type=heads#L919
	for _, item := range result.Hits.Hits {

		var pages []structs.Page

		*response = append(*response, structs.AnnotationGroup{
			ID:             item.Source["id"].(string),
			ContractID:     item.Source["contract_id"].(string),
			OpenContractID: item.Source["open_contracting_id"].(string),
			Text:           item.Source["text"].(string),
			CategoryKey:    item.Source["category_key"].(string),
			Category:       item.Source["category"].(string),
			Cluster:        item.Source["cluster"].(string),
			Pages:          pages,
		})
	}

	defer res.Body.Close()

	return response, nil
}
