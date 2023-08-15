package repository

type SDN struct {
	UID       int    `xml:"uid"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

type DB interface {
	DBClose() error
	CreateTableSellTechLCC() error
	InsertDataTable(uid, first_name, last_name string) error
	CountData() (int, error)
	GetNameFromDB(names []string, nameType string) ([]SDN, error)
}
