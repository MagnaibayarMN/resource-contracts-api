// Package structs defines data structures for contracts, annotations, and search results.
package structs

// Annotation represents a single annotation entry on a contract document.
// Annotations mark specific text regions and associate them with categories.
type Annotation struct {
	ContractID       string `json:"contract_id"`
	OpenContractID   string `json:"open_contract_id"`
	ID               int    `json:"id"`
	AnnotationID     int    `json:"annotation_id"`
	Quote            string `json:"quote"`
	Text             string `json:"string"`
	Category         string `json:"category"`
	CategoryKey      string `json:"category_key"`
	ArticleReference string `json:"article_reference"`
	PageNo           int    `json:"page_no"`
	Ranges           string `json:"ranges"`
	Cluster          string `json:"cluster"`
}

// AnnotationResponse wraps a list of annotations with a total count.
type AnnotationResponse struct {
	Total  int32       // Total number of annotations
	Result []Annotation // List of annotation objects
}

// ResultItem represents a single search result from Elasticsearch.
type ResultItem struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	Score  float32                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// Term represents an Elasticsearch term query parameter.
type Term map[string]interface{}

// Result represents a basic Elasticsearch search response with hit count.
type Result struct {
	Hits struct {
		Total int `json:"total"`
	} `json:"hits"`
}

// ResultWithHits represents an Elasticsearch search response including result documents.
type ResultWithHits struct {
	Hits struct {
		Total int          `json:"total"`
		Hits  []ResultItem `json:"hits"`
	} `json:"hits"`
}

// Shape defines the geometric bounds of an annotation on a document page.
type Shape struct {
	Type     string `json:"type"`
	Geometry struct {
		X      float64 `json:"x"`
		Y      float64 `json:"y"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"geometry"`
}

// Range represents a text range within a document.
type Range map[string]interface{}

// Page represents annotation data for a specific page of a contract document.
type Page struct {
	ID               int    `json:"id"`
	PageNo           int    `json:"page_no"`
	Quote            string `json:"quote"`
	ArticleReference string `json:"article_reference"`
	Shapes           []Shape
	Ranges           []Range
}

// AnnotationGroup groups related annotations across multiple pages.
// Annotations with the same category and text are grouped together.
type AnnotationGroup struct {
	ID             string `json:"id"`
	ContractID     string `json:"contract_id"`
	OpenContractID string `json:"open_contracting_id"`
	Text           string `json:"text"`
	CategoryKey    string `json:"category_key"`
	Category       string `json:"category"`
	Pages          []Page `json:"pages"`
	Cluster        string `json:"cluster"`
}
