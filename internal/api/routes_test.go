package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	reqBody := `{"login":"some@log.in", "password":"somepassword"}`
	req, _ := http.NewRequest("POST", "/api/user/register", strings.NewReader(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Authorization"))
}
