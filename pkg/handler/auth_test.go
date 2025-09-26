package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"todo-app"
	"todo-app/pkg/service"
	"todo-app/pkg/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(*mocks.Authorization, todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				UserName: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.Authorization, user todo.User) {
				s.On("CreateUser", user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},

		{
			name:      "Empty Fields",
			inputBody: `{"name":"Test","password":"qwerty"}`,
			mockBehavior: func(s *mocks.Authorization, user todo.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},

		{
			name:      "Service Fail",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				UserName: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.Authorization, user todo.User) {
				s.On("CreateUser", user).Return(0, errors.New("service error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockAuth := mocks.NewAuthorization(t)

			testCase.mockBehavior(mockAuth, testCase.inputUser)

			services := &service.Service{
				Authorization: mockAuth,
			}

			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(mockAuth *mocks.Authorization, userName, password string)

	testTable := []struct {
		name                 string
		requestBody          string
		inputSignIn          signInInput
		mockBehavior         mockBehavior
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name:        "OK",
			requestBody: `{"username":"test", "password":"qwerty"}`,
			inputSignIn: signInInput{
				UserName: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.Authorization, userName, password string) {
				s.On("GenerateToken", userName, password).Return("token", nil)
			},
			expectedResponseBody: `{"token":"token"}`,
			expectedStatusCode:   200,
		},

		{
			name:        "Empty Fields",
			requestBody: `{"username":"test"}`,
			mockBehavior: func(s *mocks.Authorization, userName, password string) {
			},
			expectedResponseBody: `{"message":"invalid request body"}`,
			expectedStatusCode:   400,
		},

		{
			name:        "Service Fail",
			requestBody: `{"username":"test", "password":"qwerty"}`,
			inputSignIn: signInInput{
				UserName: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.Authorization, userName, password string) {
				s.On("GenerateToken", userName, password).Return("", errors.New("service error"))
			},
			expectedResponseBody: `{"message":"service error"}`,
			expectedStatusCode:   500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockAuth := mocks.NewAuthorization(t)

			testCase.mockBehavior(mockAuth, testCase.inputSignIn.UserName, testCase.inputSignIn.Password)

			services := &service.Service{
				Authorization: mockAuth,
			}

			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

		})
	}
}
