package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/repository"
	"github.com/cmrd-a/gophermart/internal/service"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterRoute(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	config.InitConfig()
	repo := repository.Repository{PgxIface: mock}
	svc := service.NewService(context.TODO(), repo)
	router := SetupRouter(svc)

	login := "some@log.in"
	password := "somepassword"
	rows := mock.NewRows([]string{"id"}).AddRow(int64(1))
	mock.ExpectQuery("INSERT INTO users").WithArgs(login, password).WillReturnRows(rows)

	w := httptest.NewRecorder()
	reqBody := fmt.Sprintf(`{"login":"%s", "password":"%s"}`, login, password)
	req, _ := http.NewRequest("POST", "/api/user/register", strings.NewReader(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Authorization"))
}
