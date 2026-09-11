package exasol_rest_api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWriteEndpointsReturnAPIErrorResponseForMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	application := Application{}
	cases := []struct {
		name    string
		method  string
		path    string
		handler func(*gin.Context)
	}{
		{name: "execute statement", method: http.MethodPost, path: "/api/v1/statement", handler: application.ExecuteStatement},
		{name: "insert row", method: http.MethodPost, path: "/api/v1/row", handler: application.InsertRow},
		{name: "delete rows", method: http.MethodDelete, path: "/api/v1/rows", handler: application.DeleteRows},
		{name: "update rows", method: http.MethodPut, path: "/api/v1/rows", handler: application.UpdateRows},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, bytes.NewBufferString("{"))
			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.Request = request

			testCase.handler(context)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "{\"status\":\"error\",\"exception\":\"unexpected EOF\"}", response.Body.String())
			assert.False(t, context.IsAborted())
		})
	}
}
