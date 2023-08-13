package getnames

import (
	"encoding/json"
	"github.com/Zmey56/selltechllc/pkg"
	"github.com/Zmey56/selltechllc/repository"
	"net/http"
	"strings"
)

func GetNames(db *repository.DBImpl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			http.Error(w, `{"result": false, "info": "service unavailable", "code": 503}`, http.StatusServiceUnavailable)
			return
		}

		pkg.DBMutex.Lock()
		defer pkg.DBMutex.Unlock()

		nameStr := r.URL.Query().Get("name")
		nameType := r.URL.Query().Get("type")

		names := strings.Split(nameStr, " ")
		result, err := db.GetNameFromDB(names, nameType)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "GetNameFromDB unavailable", "code": 503}`))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(result)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "problem with Marshal", "code": 503}`))
		}

		w.Write(jsonResponse)
	}
}
