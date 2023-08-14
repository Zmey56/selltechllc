package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	host     = "selltechllc-db-1"
	port     = 5432
	user     = "zmey56"
	password = "zmey56"
	dbname   = "postgres"
)

type DBImpl struct {
	db *sql.DB
}

type SDN struct {
	UID       int    `xml:"uid"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

func NewPostgreSQL() (*DBImpl, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	// check db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DBImpl{db: db}, nil
}

func (db *DBImpl) DBClose() error {
	return db.db.Close()

}

func (db *DBImpl) CreateTableSellTechLCC() error {
	_, err := db.db.Exec(`CREATE TABLE IF NOT EXIST selltechlcc(
								uid INT not NULL PRIMARY KEY,
								first_name VARCHAR ( 255 ),
								last_name VARCHAR ( 255 ),
								created_on TIMESTAMP NOT NULL
							)`)
	return err
}

func (db *DBImpl) InsertDataTable(uid, first_name, last_name string) error {
	sqlStatement := `INSERT INTO selltechlcc (uid, first_name, last_name, created_on) VALUES ($1, $2, $3, $4)
 						ON CONFLICT (uid) DO UPDATE SET first_name = $2, last_name = $3, created_on = $4;`

	_, err := db.db.Exec(sqlStatement, uid, first_name, last_name, time.Now())
	if err != nil {
		log.Println("Error inserting data:", err)
		return err
	}

	return nil
}

func (db *DBImpl) CountData() (int, error) {
	sqlStatmentCount := `SELECT COUNT(*) FROM selltechlcc`
	var count int
	row := db.db.QueryRow(sqlStatmentCount)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DBImpl) GetNameFromDB(names []string, nameType string) ([]SDN, error) {
	var rows *sql.Rows
	var err error

	if len(nameType) > 0 {
		nameType = strings.ToLower(nameType)
	}

	if nameType == "strong" {
		if len(names) > 1 {
			tmpName := names[0:(len(names) - 1)]
			name := strings.Join(tmpName, " ")
			rows, err = db.db.Query("SELECT uid, first_name, last_name FROM selltechlcc WHERE"+
				" first_name = '%' || $1 || '%' AND last_name = '%' || $2 || '%'", name, names[len(names)-1])
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = db.db.Query("SELECT uid, first_name, last_name FROM selltechlcc WHERE"+
				" first_name = $1 OR last_name = $1", names[0])
			if err != nil {
				return nil, err
			}
		}
	} else {
		for _, name := range names {
			rows, err = db.db.Query("SELECT uid, first_name, last_name FROM selltechlcc WHERE"+
				" first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%'", name)
			if err != nil {
				return nil, err
			}
		}
	}

	var results []SDN
	for rows.Next() {
		var sdn SDN
		err := rows.Scan(&sdn.UID, &sdn.FirstName, &sdn.LastName)
		if err != nil {
			if err != nil {
				return nil, err
			}
		}
		results = append(results, sdn)
	}

	return results, nil
}
