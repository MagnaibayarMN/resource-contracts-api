package queries

import (
	"fmt"
	"iltodgeree/api/internal/structs"
	"strings"
)

type Highlight struct {
	FragmentSize      uint16 `json:"fragmant_size"`
	NumberOfFragments uint16 `json:"number_of_fragments"`
}

type HighlightFields struct {
	MetadataString    Highlight
	PdfTextString     Highlight
	AnnotationsString Highlight
}

type Query map[string]interface{}
type QueryFilters map[string]interface{}

type Bool struct {
	Should interface{} `json:"should"`
}

type SimpleQueryString struct {
	Fields          []string `json:"fields"`
	Query           string   `json:"query"`
	DefaultOperator string   `json:"default_operator"`
}

type FullTextQuery struct {
	Size      uint16                 `json:"size"`
	From      uint16                 `json:"from"`
	Query     Query                  `json:"query"`
	Highlight map[string]interface{} `json:"highlight"`
	Filter    QueryFilters           `json:"filter"`
	Sort      map[string]interface{} `json:"sort"`
}

type Annotation struct{}
type Group interface{}

type Result struct {
	ID                string       `json:"id"`
	OpenContractingID string       `json:"open_contracting_id"`
	Name              string       `json:"name"`
	YearSigned        string       `json:"year_signed"`
	ContractType      string       `json:"contract_type"`
	Resource          string       `json:"resource"`
	CountryCode       string       `json:"country_code"`
	Language          string       `json:"language"`
	Category          string       `json:"category"`
	IsOcrReviewed     bool         `json:"is_ocr_reviewed"`
	Metadata          string       `json:"metadata"`
	Annotations       []Annotation `json:"annotations"`
	Text              string       `json:"text"`
	Group             Group        `json:"group"`
}

func SearchInMaster(args structs.FullTextSearchArguments) {
	// var documentType = "master"
	var filters []map[string]interface{}
	var phrases []map[string]interface{}

	var types []string
	splitGroups := strings.Split(*args.Group, ",")
	for _, g := range splitGroups {
		types = append(types, strings.TrimSpace(g))
	}

	// filters
	if args.Year != nil {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.signature_year": strings.Split(*args.Year, ","),
			},
		})
	}
	if args.CountryCode != nil {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.country_code": strings.Split(*args.CountryCode, ","),
			},
		})
	}
	if args.Resource != nil {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.resource_raw": strings.Split(*args.Resource, ","),
			},
		})
	}
	if args.Category != nil {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.category": strings.Split(*args.Category, ","),
			},
		})
	}
	if args.CorporateGroup != nil {
		for _, lang := range strings.Split(*args.CorporateGroup, ",") {
			item := map[string]interface{}{
				"terms": map[string]interface{}{
					"metadata.corporate_grouping": lang,
				},
			}
			filters = append(filters, item)
		}
	}
	if args.Annotated != nil {
		annotated := map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": map[string]interface{}{
					"missing": map[string]interface{}{
						"field":     "annotations_string",
						"existence": true,
					},
				},
			},
		}
		filters = append(filters, annotated)
	}
	if args.Province != nil {
		item := map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.provinces.province": strings.Split(*args.Province, ","),
			},
		}
		filters = append(filters, item)
	}
	if args.District != nil {
		item := map[string]interface{}{
			"terms": map[string]interface{}{
				"metadata.district_id": strings.Split(*args.Province, ","),
			},
		}
		filters = append(filters, item)
	}

	// phrases
	if args.ContractType != nil {
		contractTypes := strings.Split(*args.ContractType, ",")

		for _, contractType := range contractTypes {
			item := map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"metadata.contract_type": contractType,
				},
			}
			phrases = append(phrases, item)
		}
	}

	if args.DocumentType != nil {
		documentTypes := strings.Split(*args.DocumentType, ",")

		for _, documentType := range documentTypes {
			item := map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"metadata.document_type": documentType,
				},
			}
			phrases = append(phrases, item)
		}
	}

	if args.Language != nil {
		for _, lang := range strings.Split(*args.Language, ",") {
			item := map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"metadata.language": lang,
				},
			}
			phrases = append(phrases, item)
		}
	}

	if args.CompanyName != nil {
		for _, lang := range strings.Split(*args.CompanyName, ",") {
			item := map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"metadata.company_name": lang,
				},
			}
			phrases = append(phrases, item)
		}
	}

	if args.AnnotationCategory != nil {
		for _, lang := range strings.Split(*args.AnnotationCategory, ",") {
			item := map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"annotations_category": lang,
				},
			}
			phrases = append(phrases, item)
		}
	}

	if args.Project != nil {
		item := map[string]interface{}{
			"match": map[string]interface{}{
				"metadata.project_title": map[string]interface{}{
					"query":    *args.Project,
					"operator": "AND",
				},
			},
		}
		phrases = append(phrases, item)
	}

	if args.Government != nil {
		item := map[string]interface{}{
			"match": map[string]interface{}{
				"metadata.government_entity.entity": map[string]interface{}{
					"query":    *args.Government,
					"operator": "AND",
				},
			},
		}
		phrases = append(phrases, item)
	}

	fields := []string{
		"metadata.contract_name",
		"metadata.signature_year",
		"metadata.open_contracting_id",
		"metadata.signature_date",
		"metadata.file_size",
		"metadata.country_code",
		"metadata.country_name",
		"metadata.resource",
		"metadata.language",
		"metadata.file_size",
		"metadata.company_name",
		"metadata.contract_type",
		"metadata.corporate_grouping",
		"metadata.show_pdf_text",
		"metadata.category",
		"metadata.district_id",
		"metadata.provinces.province",
		"metadata_string",
		"pdf_text_string",
		"annotations_string",
	}

	sorts := make(map[string]interface{})
	order := *args.Order

	if *args.SortBy == "country" {
		sorts["metadata.country_name"].(map[string]interface{})["order"] = order
	}
	if *args.SortBy == "year" {
		sorts["metadata.year"].(map[string]interface{})["order"] = order
	}
	if *args.SortBy == "contract_name" {
		sorts["metadata.contract_name"].(map[string]interface{})["order"] = order
	}
	if *args.SortBy == "resource" {
		sorts["metadata.resource"].(map[string]interface{})["order"] = order
	}
	if *args.SortBy == "contract_type" {
		sorts["metadata.contract_type"].(map[string]interface{})["order"] = order
	}

	highlights := HighlightFields{
		MetadataString: Highlight{
			FragmentSize:      200,
			NumberOfFragments: 1,
		},
		PdfTextString: Highlight{
			FragmentSize:      200,
			NumberOfFragments: 50,
		},
		AnnotationsString: Highlight{
			FragmentSize:      50,
			NumberOfFragments: 1,
		},
	}

	bodyHighlight := map[string]interface{}{
		"pre_tags":  []string{"<strong>"},
		"post_tags": []string{"</strong>"},
		"fields":    highlights,
	}

	var query map[string]interface{}
	var queryFilters map[string]interface{}

	if len(phrases) > 0 && (*args.Q == "" || args.Q == nil) {
		query = map[string]interface{}{
			"bool": Bool{
				Should: phrases,
			},
		}
	}

	if *args.Q != "" && args.Q != nil {
		query = map[string]interface{}{
			"simple_query_string": SimpleQueryString{
				Fields:          fields,
				Query:           *args.Q,
				DefaultOperator: "AND",
			},
		}
	}

	if len(filters) > 0 {
		queryFilters = map[string]interface{}{
			"filter": map[string]interface{}{
				"and": map[string]interface{}{
					"filters": filters,
				},
			},
		}
	}

	body := FullTextQuery{
		Size:      *args.PerPage,
		From:      *args.From,
		Query:     query,
		Filter:    queryFilters,
		Highlight: bodyHighlight,
		Sort:      sorts,
	}

	fmt.Println(body)
}

func Search() {

}
