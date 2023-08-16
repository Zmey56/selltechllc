package handler

import (
	"encoding/json"
	"encoding/xml"
	"github.com/Zmey56/selltechllc/pkg"
	"github.com/Zmey56/selltechllc/repository"
	"io"
	"net/http"
	"strings"
)

type State struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

type SDNList struct {
	XMLName xml.Name `xml:"sdnList"`
	SDNs    []SDN    `xml:"sdnEntry"`
}

type SDN struct {
	UID          string `xml:"uid"`
	SDNFirstName string `xml:"firstName"`
	SDNLastName  string `xml:"lastName"`
	sdnType      string `xml:"sdnType"`
}

func GetNames(db repository.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			http.Error(w, `{"result": false, "info": "service unavailable", "code": 503}`, http.StatusServiceUnavailable)
			return
		}

		pkg.DBMutex.Lock()
		defer pkg.DBMutex.Unlock()

		nameStr := r.URL.Query().Get("name")
		nameType := r.URL.Query().Get("type")

		if len(strings.TrimSpace(nameStr)) < 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "GetNameFromDB unavailable", "code": 503}`))
		} else {
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
}

func StateHandler(db repository.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			stateResponse(w, false, "updating")
			return
		}

		vol, err := db.CountData()
		if err != nil {
			http.Error(w, `{"result": false, "info": "DB unavailable", "code": 503}`, http.StatusServiceUnavailable)
			return
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

func UpdateHandler(db repository.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
			return
		}

		pkg.DBMutex.Lock()
		defer pkg.DBMutex.Unlock()

		pkg.UpdatingFlag = true
		defer func() {
			pkg.UpdatingFlag = false
		}()

		resp, err := http.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
			return
		}
		defer resp.Body.Close()

		xmlData, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
			return
		}

		sdnList := SdnList{}

		if err := xml.Unmarshal(xmlData, &sdnList); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
			return
		}

		for _, j := range sdnList.SdnEntry {
			if j.SdnType == "Individual" {
				db.InsertDataTable(j.Uid, j.FirstName, j.LastName)
			}
		}

		// Respond with a success message if everything went well
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": true, "info": "update successful", "code": 200}`))

	}
}

//Ответ на 3 вопрос: Обновление внешних данных может быть асинхронным, чтобы не блокировать обработку запроса. Так же можно
//сохранять время последнего обновления данных.
