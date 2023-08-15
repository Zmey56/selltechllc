package main

import (
	"fmt"
	"github.com/Zmey56/selltechllc/getnames"
	"github.com/Zmey56/selltechllc/repository"
	"github.com/Zmey56/selltechllc/statehandler"
	"github.com/Zmey56/selltechllc/updatehandler"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := repository.NewPostgreSQL(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Failed to connect to DB: %s\n", err)
		os.Exit(1)
	}
	defer db.DBClose()

	//err = db.CreateTableSellTechLCC()
	if err != nil {
		log.Fatal(err)
	}

	//create routers
	routers := mux.NewRouter()
	routers.HandleFunc("/update", updatehandler.UpdateHandler(db))
	routers.HandleFunc("/state", statehandler.StateHandler(db))
	routers.HandleFunc("/get_names", getnames.GetNames(db))

	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(routers)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
