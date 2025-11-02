package document

type Metadata struct {
	ContractName string
}

type ContractPage struct {
	ContractId int64
	Id         int64
	PageNo     int64
	Text       string
	CreatedAt  string
	UpdatedAt  string
}

// `gorm:"type:json"`
type Contract struct {
	Id              int64
	Metadata        string
	File            string
	Filehas         string
	UserID          int64
	CreatedDateTime string
	ParentID        int64
	ProvinceID      int64
	DistrictID      int64
	ContractPages   []ContractPage
}
