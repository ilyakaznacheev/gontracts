package gontracts

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testCheckAuthResponse(loc string, t *testing.T, w *httptest.ResponseRecorder, respStatus int) {
	if w.Code != respStatus {
		t.Errorf("[%s]:\twrong StatusCode: got %d, expected %d",
			loc, w.Code, respStatus)
	}

	body, _ := ioutil.ReadAll(w.Result().Body)
	bodyStr := string(body)
	if len(bodyStr) == 0 {
		t.Errorf("[%s]:\tempty Response body", loc)
	}
}
func TestGetToken(t *testing.T) {
	a := NewAuthHandler([]byte("test"))

	url := "/get-token"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	testHandle("/get-token", w, req, a.GenerateToken)
	testCheckAuthResponse("GenerateToken:", t, w, http.StatusOK)

}
