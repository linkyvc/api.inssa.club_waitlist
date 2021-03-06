package test

import (
	"bytes"
	"encoding/json"
	"inssa_club_waitlist_backend/cmd/server/errors"
	"inssa_club_waitlist_backend/cmd/server/forms"
	"inssa_club_waitlist_backend/cmd/server/middlewares"
	"inssa_club_waitlist_backend/cmd/server/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"
	"gopkg.in/guregu/null.v4"
)

var engine *gin.Engine

func performRequest(h http.Handler, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func performRequestWithForm(h http.Handler, method string, path string, form interface{}) *httptest.ResponseRecorder {
	formJSON, _ := json.Marshal(form)
	body := strings.NewReader(string(formJSON))
	return performRequest(h, method, path, body)
}

func responseToMap(body *bytes.Buffer, resultMap *map[string]interface{}) {
	stringBody := body.String()
	bytesBody := []byte(stringBody)
	json.Unmarshal(bytesBody, resultMap)
}

// helper functions for easier request

func requestInterestWithoutEmailTest(t *testing.T) {
	var response map[string]interface{}

	form := &forms.AddInterestRequest{
		ClubhouseUserID: null.NewInt(123, true),
	}

	w := performRequestWithForm(engine, "POST", "/interest", form)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseToMap(w.Body, &response)
	assert.Equal(t, errors.ValidationError, response["error_type"])
}

func requestInterestWithoutProperEmailTest(t *testing.T) {
	var response map[string]interface{}

	form := &forms.AddInterestRequest{
		ClubhouseUserID: null.NewInt(123, true),
		Email:           "example",
	}

	w := performRequestWithForm(engine, "POST", "/interest", form)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseToMap(w.Body, &response)
	assert.Equal(t, errors.ValidationError, response["error_type"])
}

func requestInterest(t *testing.T) {
	form := &forms.AddInterestRequest{
		ClubhouseUserID: null.NewInt(123, true),
		Email:           "example@example.com",
	}

	w := performRequestWithForm(engine, "POST", "/interest", form)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func requestInterestWithDuplicateEmailTest(t *testing.T) {
	var response map[string]interface{}

	form := &forms.AddInterestRequest{
		ClubhouseUserID: null.NewInt(123, true),
		Email:           "example@example.com",
	}

	w := performRequestWithForm(engine, "POST", "/interest", form)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseToMap(w.Body, &response)
	assert.Equal(t, errors.DuplicateEmailError, response["error_type"])
}

func deleteInterest(t *testing.T) {
	form := &forms.AddInterestRequest{
		Email: "example@example.com",
	}

	w := performRequestWithForm(engine, "DELETE", "/interest", form)
	assert.Equal(t, http.StatusOK, w.Code)
}

// test interest

func TestEngineSetup(t *testing.T) {
	engine = gin.New()
	utils.InitDB()
	middlewares.Setup(engine)
	setupRoutes(engine)
}

func TestInterest(t *testing.T) {
	t.Run("Test interest feature", func(t *testing.T) {
		t.Run("Request interest without email should fail", requestInterestWithoutEmailTest)
		t.Run("Request interest without proper email should fail", requestInterestWithoutProperEmailTest)
		t.Run("Request with proper email, user id should success", requestInterest)
		t.Run("Request interest with duplicate email should fail", requestInterestWithDuplicateEmailTest)
		t.Run("Request delete an interest should success", deleteInterest)
	})
}
