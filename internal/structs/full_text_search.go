// Package structs defines data structures for API requests and responses.
package structs

// FullTextSearchArguments contains all possible search parameters for contract queries.
// All fields are optional pointers to distinguish between unset and zero values.
type FullTextSearchArguments struct {
	Year               *string `json:"year"`
	CountryCode        *string `json:"country_code"`
	Resource           *string `json:"resource"`
	Category           *string `json:"category"`
	ContractType       *string `json:"contract_type"`
	DocumentType       *string `json:"document_type"`
	Language           *string `json:"language"`
	CompanyName        *string `json:"company_name"`
	CorporateGroup     *string `json:"corporate_group"`
	AnnotationCategory *string `json:"annotation_category"`
	Annotated          *bool   `json:"annotated"`
	Province           *string `json:"province"`
	District           *string `json:"district"`
	Project            *string `json:"project"`
	Government         *string `json:"government"`
	Q                  *string `json:"q"`
	SortBy             *string `json:"sort_by"`
	Order              *string `json:"order"`
	Download           *bool   `json:"download"`
	Group              *string `json:"group"`
	From               *uint16 `json:"from"`
	PerPage            *uint16 `json:"per_page"`
}
