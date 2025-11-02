package queries

import (
	"encoding/json"
	"fmt"
	ctx "iltodgeree/api/internal/app_context"
	"iltodgeree/api/internal/structs"
	"strconv"
	"strings"
)

func GetStatesCount(id int64) (int, error) {
	documentType := "state_contracts"

	arguments := map[string]interface{}{
		"provinces.province": strconv.FormatInt(id, 10),
	}
	query, err := ctx.ElasticInstance.BuildQueryBoolMust(arguments)
	if err != nil {
		fmt.Println("Error marshaling to JSON query:", err)
		return 0, err
	}

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithDocumentType(documentType),
		client.Search.WithSize(defaultSize),
		client.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return 0, fmt.Errorf("error executing search: %v", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error in Elasticsearch response: %s", res.String())
	}

	var result structs.ResultWithHits

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	return len(result.Hits.Hits), nil
}
