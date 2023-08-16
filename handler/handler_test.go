package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zmey56/selltechllc/pkg"
	"github.com/Zmey56/selltechllc/repository"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

type FakeDB struct {
	ExpectedCountDataResult int
	CountDataError          error
	InsertCalled            bool
}

var httpGet = http.Get

func (f FakeDB) DBClose() error {
	//TODO implement me
	panic("implement me")
}

func (f FakeDB) CreateTableSellTechLCC() error {
	//TODO implement me
	panic("implement me")
}

func (f *FakeDB) InsertDataTable(uid, firstName, lastName string) error {
	f.InsertCalled = true
	return nil
}

func (fdb *FakeDB) SetCountDataResponse(result int, err error) {
	fdb.ExpectedCountDataResult = result
	fdb.CountDataError = err
}

func (fdb *FakeDB) CountData() (int, error) {
	return fdb.ExpectedCountDataResult, fdb.CountDataError
}

func (f *FakeDB) GetNameFromDB(names []string, nameType string) ([]repository.SDN, error) {
	fakeData := []repository.SDN{
		{UID: 1, FirstName: "Alice", LastName: "Johnson"},
		{UID: 2, FirstName: "Bob", LastName: "Smith"},
	}
	return fakeData, nil
}

func NewFakeDB() *FakeDB {
	return &FakeDB{}
}

func TestGetNamesHandler(t *testing.T) {
	fakeDB := NewFakeDB()
	handler := GetNames(fakeDB)

	tests := []struct {
		nameParam   string
		typeParam   string
		expectCode  int
		expectError bool
	}{
		{"Alice Bob", "full_name", http.StatusOK, false},
		{"Alice", "first_name", http.StatusOK, false},
		{"", "last_name", http.StatusServiceUnavailable, true},
		{"Charlie", "unknown_type", http.StatusOK, false},
		{"", "", http.StatusServiceUnavailable, true},
		{" ", "", http.StatusServiceUnavailable, true},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/?name="+url.QueryEscape(test.nameParam)+"&type="+test.typeParam, nil)
		rr := httptest.NewRecorder()

		pkg.UpdatingFlag = false

		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectCode {
			t.Errorf("For name='%s', type='%s', expected status code %d, but got %d", test.nameParam, test.typeParam, test.expectCode, rr.Code)
		}

		if test.expectError {
			expectedBody := `{"result": false, "info": "GetNameFromDB unavailable", "code": 503}`
			if rr.Body.String() != expectedBody {
				t.Errorf("For name='%s', type='%s', expected body %s, but got %s", test.nameParam, test.typeParam, expectedBody, rr.Body.String())
			}
		} else {
			var responseData []repository.SDN
			err := json.Unmarshal(rr.Body.Bytes(), &responseData)
			if err != nil {
				t.Errorf("For name='%s', type='%s', error decoding response body: %v", test.nameParam, test.typeParam, err)
			}
		}
	}
}

func TestStateHandler(t *testing.T) {
	fakeDB := NewFakeDB()
	handler := StateHandler(fakeDB)

	tests := []struct {
		UpdatingFlag bool
		Volume       int
		ExpectResult bool
		ExpectInfo   string
		ExpectCode   int
	}{
		{false, 0, false, "empty", http.StatusOK},
		{false, 5, true, "ok", http.StatusOK},
		{true, 0, false, "updating", http.StatusOK},
	}

	for _, test := range tests {
		pkg.UpdatingFlag = test.UpdatingFlag
		fakeDB.SetCountDataResponse(test.Volume, nil)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != test.ExpectCode {
			t.Errorf("Expected status code %d, but got %d", test.ExpectCode, rr.Code)
		}

		var responseData State
		err := json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}

		if responseData.Result != test.ExpectResult {
			t.Errorf("Expected result %v, but got %v for %v", test.ExpectResult, responseData.Result, test.ExpectInfo)
		}

		if responseData.Info != test.ExpectInfo {
			t.Errorf("Expected info '%s', but got '%s'", test.ExpectInfo, responseData.Info)
		}
	}
}

func TestStateHandler_DBUnavailable(t *testing.T) {
	fakeDB := NewFakeDB()
	handler := StateHandler(fakeDB)

	// Установите ожидаемый результат и ошибку для метода CountData
	fakeDB.SetCountDataResponse(0, errors.New("fake DB error"))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusServiceUnavailable
	if rr.Code != expectedStatus {
		t.Errorf("Expected status code %d, but got %d", expectedStatus, rr.Code)
	}

	expectedBody := `{"result": false, "info": "DB unavailable", "code": 503}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, rr.Body.String())
	}
}

func TestUpdateHandler(t *testing.T) {
	db := &FakeDB{}
	pkg.DBMutex = sync.Mutex{}

	ts := httptest.NewServer(http.HandlerFunc(UpdateHandler(db)))
	defer ts.Close()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	pkg.UpdatingFlag = false

	UpdateHandler(db)(resp, req)

	// Validate response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Header().Get("Content-Type"), "application/json") // Change this line
	assert.Contains(t, resp.Body.String(), `"result": true`)
	assert.Contains(t, resp.Body.String(), `"info": "update successful"`)
	assert.Contains(t, resp.Body.String(), `"code": 200`)

	// Validate DB insert
	assert.True(t, db.InsertCalled)

	// Test concurrent access
	concurrency := 10
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", ts.URL, nil)
			assert.NoError(t, err)

			UpdateHandler(db)(resp, req)
		}()
	}

	wg.Wait()

	// Additional testing for UpdatingFlag
	pkg.UpdatingFlag = true

	resp = httptest.NewRecorder()
	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	UpdateHandler(db)(resp, req)

	// Validate response when UpdatingFlag is true
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, resp.Body.String(), `"result": false`)
	assert.Contains(t, resp.Body.String(), `"info": "service unavailable"`)
	assert.Contains(t, resp.Body.String(), `"code": 503`)

}

func TestUpdateHandlerHTTPGetError(t *testing.T) {
	db := &FakeDB{}
	pkg.DBMutex = sync.Mutex{} // Make sure pkg.DBMutex is initialized

	ts := httptest.NewServer(http.HandlerFunc(UpdateHandler(db)))
	defer ts.Close()

	pkg.UpdatingFlag = false

	// Mock http.Get to return an error
	httpGet = func(url string) (*http.Response, error) {
		return nil, io.EOF // Simulate an error
	}
	defer func() { httpGet = http.Get }() // Reset the mock

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	UpdateHandler(db)(resp, req) // Use resp and req here

	// Validate response when http.Get returns an error
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, resp.Body.String(), `"result": false`)
	assert.Contains(t, resp.Body.String(), `"info": "service unavailable"`)
	assert.Contains(t, resp.Body.String(), `"code": 503`)
}

func TestUpdateHandlerXMLUnmarshalError(t *testing.T) {
	db := &FakeDB{}
	pkg.DBMutex = sync.Mutex{} // Make sure pkg.DBMutex is initialized

	ts := httptest.NewServer(http.HandlerFunc(UpdateHandler(db)))
	defer ts.Close()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	pkg.UpdatingFlag = false

	// Mock http.Get to return a valid response with invalid XML data
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer([]byte("invalid-xml"))),
		}, nil
	}
	defer func() { httpGet = http.Get }() // Reset the mock

	UpdateHandler(db)(resp, req)

	// Validate response when XML unmarshal fails
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, resp.Body.String(), `"result": false`)
	assert.Contains(t, resp.Body.String(), `"info": "service unavailable"`)
	assert.Contains(t, resp.Body.String(), `"code": 503`)
}

func TestUpdateHandlerConcurrentAccess(t *testing.T) {
	db := &FakeDB{}
	pkg.DBMutex = sync.Mutex{} // Make sure pkg.DBMutex is initialized

	ts := httptest.NewServer(http.HandlerFunc(UpdateHandler(db)))
	defer ts.Close()

	pkg.UpdatingFlag = false

	// Mock http.Get to return valid XML data
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer([]byte(sampleXML))),
		}, nil
	}
	defer func() { httpGet = http.Get }() // Reset the mock

	concurrency := 10
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", ts.URL, nil)
			assert.NoError(t, err)

			UpdateHandler(db)(resp, req)
		}()
	}

	wg.Wait()

	// Validate DB insert for each goroutine
	assert.True(t, db.InsertCalled)
}

var sampleXML = `
<?xml version="1.0" encoding="UTF-8"?>
<sdnList>
	<sdnEntry>
		<uid>1</uid>
		<firstName>John</firstName>
		<lastName>Doe</lastName>
		<sdnType>Individual</sdnType>
	</sdnEntry>
	<sdnEntry>
		<uid>2</uid>
		<firstName>Jane</firstName>
		<lastName>Smith</lastName>
		<sdnType>Individual</sdnType>
	</sdnEntry>
</sdnList>
`
