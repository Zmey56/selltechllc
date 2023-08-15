package updatehandler

import (
	"encoding/xml"
	"github.com/Zmey56/selltechllc/pkg"
	"github.com/Zmey56/selltechllc/repository"
	"io"
	"net/http"
)

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

func UpdateHandler(db repository.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pkg.UpdatingFlag {
			http.Error(w, `{"result": false, "info": "service unavailable", "code": 503}`, http.StatusServiceUnavailable)
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
		}
		defer resp.Body.Close()

		xmlData, err := io.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
			}
		}

		sdnList := SdnList{}

		if err := xml.Unmarshal(xmlData, &sdnList); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"result": false, "info": "service unavailable", "code": 503}`))
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
