package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"todo-app/pkg/service"
	"todo-app/pkg/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(*mocks.Authorization, string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(a *mocks.Authorization, token string) {
				a.On("ParseToken", token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},

		{
			name:                 "No Header",
			mockBehavior:         func(a *mocks.Authorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},

		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			mockBehavior:         func(a *mocks.Authorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},

		{
			name:                 "No token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(a *mocks.Authorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},

		{
			name:        "Service Fail",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(a *mocks.Authorization, token string) {
				a.On("ParseToken", token).Return(0, errors.New("failed to parse token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			mockAuth := mocks.NewAuthorization(t)

			testCase.mockBehavior(mockAuth, testCase.token)

			services := &service.Service{Authorization: mockAuth}

			handler := NewHandler(services)

			r := gin.New()
			r.POST("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)

				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

		})
	}

}
