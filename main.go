package main

import (
	"fmt"
	"github.com/Zmey56/selltechllc/driver"
	"github.com/Zmey56/selltechllc/getnames"
	"github.com/Zmey56/selltechllc/repository"
	"github.com/Zmey56/selltechllc/repository/dbrepo"
	"github.com/Zmey56/selltechllc/statehandler"
	"github.com/Zmey56/selltechllc/updatehandler"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Repository struct {
	DB repository.DB
}

// NewRepo creates a new repository
func NewRepo(db *driver.DB) *Repository {
	return &Repository{
		DB: dbrepo.NewPostgresRepo(db.SQL),
	}
}

func main() {
	db, err := driver.ConnectSQL(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Failed to connect to DB: %s\n", err)
		os.Exit(1)
	}
	repo := NewRepo(db)

	defer db.SQL.Close()

	//err = db.CreateTableSellTechLCC()
	if err != nil {
		log.Fatal(err)
	}

	//create routers
	routers := mux.NewRouter()
	routers.HandleFunc("/update", updatehandler.UpdateHandler(repo.DB))
	routers.HandleFunc("/state", statehandler.StateHandler(repo.DB))
	routers.HandleFunc("/get_names", getnames.GetNames(repo.DB))

	log.Fatal(http.ListenAndServe(":8080", jsonContentTypeMiddleware(routers)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
