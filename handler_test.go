package main

import "testing"
import "net/http"
import "net/http/httptest"
import "errors"
import "time"
import "io/ioutil"
import "encoding/json"

type MockSunsetFinder struct {
	queryFunc func(location string) (sunsetResult, error)
}

func (msf *MockSunsetFinder) Query(location string) (sunsetResult, error) {
	return msf.queryFunc(location)
}

func Test_LocationNotSet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?", nil)

	w := httptest.NewRecorder()

	h := getHandler(nil)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("code should be %v", http.StatusBadRequest)
	}
}

func Test_ApiReturnsError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?location=Jacksonville,+FL", nil)

	w := httptest.NewRecorder()

	msf := &MockSunsetFinder{}
	msf.queryFunc = func(location string) (sunsetResult, error) {
		return sunsetResult{}, errors.New("api's broke man")
	}

	h := getHandler(msf)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("code should be %v", http.StatusBadGateway)
	}
}

func Test_LocationNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?location=Notown,+AA", nil)

	w := httptest.NewRecorder()

	msf := &MockSunsetFinder{}
	msf.queryFunc = func(location string) (sunsetResult, error) {
		return sunsetResult{Sunset: "foobar"}, errNotFound
	}

	h := getHandler(msf)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("code should be %v", http.StatusNotFound)
	}
}

func Test_OnlyGetSupported(t *testing.T) {

	for _, method := range [...]string{http.MethodDelete,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut} {

		req, _ := http.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()

		h := getHandler(nil)

		h.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("code should be %v", http.StatusMethodNotAllowed)
		}
	}
}

func Test_LocationLookupWorks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?location=Jacksonville,+FL", nil)

	w := httptest.NewRecorder()

	tm := time.Now()
	msf := &MockSunsetFinder{}
	msf.queryFunc = func(location string) (sunsetResult, error) {
		return sunsetResult{Sunset: "5:30 PM", Timestamp: tm}, nil
	}

	h := getHandler(msf)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code should be %v", http.StatusOK)
	}

	body, _ := ioutil.ReadAll(w.Body)

	var result sunsetResult
	json.Unmarshal(body, &result)

	if result.Sunset != "5:30 PM" {
		t.Error("Didn't write result properly")
	}
}
