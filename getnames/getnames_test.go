package getnames

import (
	"encoding/json"
	"errors"
	"github.com/Zmey56/selltechllc/repository/dbrepo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDBImpl struct {
	dbrepo.DBImpl
	GetNameFromDBFunc func(names []string, nameType string) ([]string, error)
}

func (m *MockDBImpl) GetNameFromDB(names []string, nameType string) ([]string, error) {
	if m.GetNameFromDBFunc != nil {
		return m.GetNameFromDBFunc(names, nameType)
	}
	return nil, nil // Вернуть пустой результат в случае отсутствия мока
}

func TestGetNamesHandler(t *testing.T) {
	t.Run("Successful request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/?name=John%20Doe&type=person", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		mockDB := &MockDBImpl{}
		mockDB.GetNameFromDBFunc = func(names []string, nameType string) ([]string, error) {
			return []string{"John Doe"}, nil
		}

		// Вызов обработчика
		handler := GetNames(mockDB)
		handler.ServeHTTP(recorder, req)

		// Проверка статуса и ответа
		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, []string{"John Doe"}, response)
	})

	t.Run("Failed database request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/?name=Jane%20Doe&type=person", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		mockDB := &MockDBImpl{}
		mockDB.GetNameFromDBFunc = func(names []string, nameType string) ([]string, error) {
			return nil, errors.New("database error")
		}

		handler := GetNames(mockDB)
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, false, response["result"])
		assert.Equal(t, "GetNameFromDB unavailable", response["info"])
		assert.Equal(t, float64(503), response["code"])
	})

	t.Run("JSON marshaling error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/?name=Jane%20Doe&type=person", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		mockDB := &MockDBImpl{}
		mockDB.GetNameFromDBFunc = func(names []string, nameType string) ([]string, error) {
			return []string{"Jane Doe"}, nil
		}

		// Эмуляция ошибки при маршалинге JSON
		badJSONMarshal := make(chan interface{})
		close(badJSONMarshal)

		handler := GetNames(mockDB)
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, false, response["result"])
		assert.Equal(t, "problem with Marshal", response["info"])
		assert.Equal(t, float64(503), response["code"])
	})
}
