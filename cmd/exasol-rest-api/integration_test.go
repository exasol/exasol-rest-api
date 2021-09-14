package exasol_rest_api_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	exasol_rest_api "main/cmd/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/query/:query", exasol_rest_api.Query)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT 1 FROM DUAL", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	fmt.Println(responseRecorder.Body)

	if responseRecorder.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, responseRecorder.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, responseRecorder.Code)
	}
}
