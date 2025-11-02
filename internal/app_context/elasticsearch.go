// Package appcontext provides application-level context and shared resources.
// It manages Elasticsearch client connections and query building utilities.
package appcontext

import (
	"encoding/json"
	"fmt"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v5"
	"gopkg.in/olivere/elastic.v5"
)

// Elastic manages Elasticsearch client instances for different API versions.
// It provides lazy initialization of clients and query building utilities.
type Elastic struct {
	client   *elasticsearch.Client // Elasticsearch v7+ client
	clientV5 *elastic.Client       // Elasticsearch v5 client
}

var config = elasticsearch.Config{
	Username: os.Getenv("ELASTICSEARCH_USERNAME"),
	Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
}

// GetV5 returns an Elasticsearch v5 client instance.
// It creates the client on first call and reuses it for subsequent calls.
func (e Elastic) GetV5() (clientV5 *elastic.Client, err error) {
	if e.clientV5 == nil {
		e.clientV5, err = elastic.NewClient(
			elastic.SetURL(os.Getenv("ELASTICSEARCH_HOST")),
			elastic.SetSniff(false),
		)
	}

	return e.clientV5, err
}

// Get returns an Elasticsearch v7+ client instance.
// It creates the client on first call and reuses it for subsequent calls.
func (e Elastic) Get() (client *elasticsearch.Client, err error) {
	if e.client == nil {
		e.client, err = elasticsearch.NewClient(config)
	}

	return e.client, err
}

// BuildQueryBoolMust constructs an Elasticsearch bool query with a must clause.
// It takes a map of field-value pairs and generates a JSON query string.
//
// Parameters:
//   - arguments: A map of search terms to build the query from
//
// Returns:
//   - query: JSON string containing the Elasticsearch query
//   - err: Error if JSON marshaling fails
func (e Elastic) BuildQueryBoolMust(arguments interface{}) (query string, err error) {
	jsonBytes, err := json.Marshal(arguments)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return "", err
	}

	query = `{
		"query": {
			"bool": {
				"must": {
					"term": ` + string(jsonBytes) + `
				}
			}
		}
	}`

	return query, err
}

// ElasticInstance is the global Elasticsearch client instance used throughout the application.
var ElasticInstance Elastic
