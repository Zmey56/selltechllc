package repository

type DatabaseRepo interface {
	Close() error
	CreateTableSellTechLCC() error
	InsertDataTable(uid, first_name, last_name string) error
	CountData() (int, error)
	GetNameFromDB(names []string, nameType string) ([]SDN, error)
}
