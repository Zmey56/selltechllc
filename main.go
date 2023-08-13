package main

import (
	"fmt"
	"github.com/Zmey56/selltechllc/getnames"
	"github.com/Zmey56/selltechllc/repository"
	"github.com/Zmey56/selltechllc/statehandler"
	"github.com/Zmey56/selltechllc/updatehandler"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	db, err := repository.NewPostgreSQL()
	if err != nil {
		fmt.Printf("Failed to connect to DB: %s\n", err)
		os.Exit(1)
	}
	defer db.DBClose()

	http.HandleFunc("/update", updatehandler.UpdateHandler(db))
	http.HandleFunc("/state", statehandler.StateHandler(db))
	http.HandleFunc("/get_names", getnames.GetNames(db))
	http.ListenAndServe("localhost:8080", nil)
}
