package queries

import (
	ctx "iltodgeree/api/internal/app_context"
	"os"
)

var client, _ = ctx.ElasticInstance.Get()
var indexName = os.Getenv("ELASTICSEARCH_INDEX")
var defaultSize = 10000
