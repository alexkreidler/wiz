package httpjson

import (
	"fmt"
	"github.com/alexkreidler/wiz/taguk"
	"github.com/alexkreidler/wiz/taguk/test"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServe(t *testing.T) {
	api := taguk.NewServer()
	api.AddResources(test.Branch{})
	s := NewServer()
	s.Configure(api.Resources)

	spew.Dump(s.Resources)
	spew.Dump(s.Router)

	r := s.Router

	req, err := http.NewRequest("GET", "/branch/12/get", nil)
	test.Ok(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our router satisfies http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		rr.Body.String(), expected)
	//}
	fmt.Println(rr.Body.String())
}
