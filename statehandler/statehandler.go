package statehandler

import (
	"encoding/json"
	"github.com/Zmey56/selltechllc/pkg"
	"github.com/Zmey56/selltechllc/repository"
	"net/http"
)

type State struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

func StateHandler(db repository.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			stateResponse(w, false, "updating")
			return
		}

		vol, err := db.CountData()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "DB unavailable", "code": 503}`))
		}

		if vol == 0 {
			stateResponse(w, false, "empty")
		} else {
			stateResponse(w, true, "ok")
		}
	}
}

func stateResponse(w http.ResponseWriter, result bool, info string) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := State{
		Result: result,
		Info:   info,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Write(jsonResponse)

	return nil
}
