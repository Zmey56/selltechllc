package driver

import "database/sql"

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func ConnectSQL(dsn string) (*DB, error) {
	// open database
	d, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// check db
	err = d.Ping()
	if err != nil {
		return nil, err
	}

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
