package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupPort(t *testing.T) {
	t.Run("UnsetEnv", func(t *testing.T) {
		setupPort()

		defer os.Unsetenv("PORT")

		got := os.Getenv("PORT")
		want := "1323"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("SetEnv", func(t *testing.T) {
		os.Setenv("PORT", "5000")

		setupPort()

		defer os.Unsetenv("PORT")

		got := os.Getenv("PORT")
		want := "5000"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestHealthHandler(t *testing.T) {
	e := initEcho()

	req := httptest.NewRequest(echo.GET, "/health", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetHandler(t *testing.T) {
	var tests = []struct {
		name        string
		route       string
		body        string
		message     string
		redisKey    string
		redisValue  string
		dbConnected bool
		statusCode  int
	}{
		{
			name:        "EmptyRequestBody",
			route:       "/api/test-route",
			body:        ``,
			message:     `{"test": 20}`,
			redisKey:    "9074f62003bcbed6e87000ad55c501754308685b",
			redisValue:  `{"test": 20}`,
			dbConnected: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "WithRequestBody",
			route:       "/api/test-route",
			body:        `{"test": 10}`,
			message:     `{"test": 20}`,
			redisKey:    "1f1eee663738854c4e53bf7be7902de982f22255",
			redisValue:  `{"test": 20}`,
			dbConnected: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "DataUnsetInRedis",
			route:       "/api/test-route",
			body:        `{"body": {"test": 10}}`,
			message:     `{"message":"data needs to be set before GET"}`,
			redisKey:    "b99c071333d4dbca0d9298e5c8d7480f176cafdc",
			redisValue:  `{"test": 20}`,
			dbConnected: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "BadRoute",
			route:       "/bad-api/test-route",
			body:        ``,
			message:     `{"message":"Not Found"}`,
			dbConnected: true,
			statusCode:  http.StatusNotFound,
		},
		{
			name:        "MissingConnection",
			route:       "/api/test-route",
			body:        ``,
			message:     `{"message":"cannot connect to database"}`,
			redisKey:    "9074f62003bcbed6e87000ad55c501754308685b",
			redisValue:  `{"test": 20}`,
			dbConnected: false,
			statusCode:  http.StatusInternalServerError,
		},
	}

	e := initEcho()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.dbConnected {
				miniredis, err := miniredis.Run()
				if err != nil {
					panic(err)
				}
				defer miniredis.Close()

				os.Setenv("REDIS_URL", miniredis.Addr())
				defer os.Unsetenv("REDIS_URL")
			}

			redisClient().Set(ctx, test.redisKey, test.redisValue, 7*24*time.Hour)
			req := httptest.NewRequest(echo.GET, test.route, strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			if test.statusCode != http.StatusOK {
				assert.Equal(t, test.statusCode, rec.Code)
				assert.Equal(t, test.message+"\n", rec.Body.String())
			} else {
				assert.Equal(t, test.statusCode, rec.Code)
				assert.Equal(t, test.message, rec.Body.String())
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	var tests = []struct {
		name        string
		route       string
		body        string
		message     string
		dbConnected bool
		statusCode  int
	}{
		{
			name:        "EmptyRequestBody",
			route:       "/api/test-route",
			body:        ``,
			message:     `{"message":"there is no request body"}`,
			dbConnected: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "NoBodyKey",
			route:       "/api/test-route",
			body:        `{"response": {"test": 20}}`,
			message:     `{"id":"9074f62003bcbed6e87000ad55c501754308685b"}`,
			dbConnected: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "MissingBodyKey",
			route:       "/api/test-route",
			body:        `{"bod": {"test": 10}, "response": {"test": 20}}`,
			message:     `{"id":"9074f62003bcbed6e87000ad55c501754308685b"}`,
			dbConnected: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "BodyAndResponseKeys",
			route:       "/api/test-route",
			body:        `{"body": {"test": 10}, "response": {"test": 20}}`,
			message:     `{"id":"1f1eee663738854c4e53bf7be7902de982f22255"}`,
			dbConnected: true,
			statusCode:  http.StatusOK,
		},
		{
			name:        "MissingResponseKey",
			route:       "/api/test-route",
			body:        `{"body": {"test": 10}}`,
			message:     `{"message":"missing request response"}`,
			dbConnected: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "BadRoute",
			route:       "/bad-api/test-route",
			body:        ``,
			message:     `{"message":"Not Found"}`,
			dbConnected: true,
			statusCode:  http.StatusNotFound,
		},
		{
			name:        "MissingConnection",
			route:       "/api/test-route",
			body:        `{"body": {"test": 10}, "response": {"test": 20}}`,
			message:     `{"message":"cannot connect to database"}`,
			dbConnected: false,
			statusCode:  http.StatusInternalServerError,
		},
	}

	e := initEcho()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.dbConnected {
				miniredis, err := miniredis.Run()
				if err != nil {
					panic(err)
				}
				defer miniredis.Close()

				os.Setenv("REDIS_URL", miniredis.Addr())
				defer os.Unsetenv("REDIS_URL")
			}

			req := httptest.NewRequest(echo.POST, test.route, strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, test.statusCode, rec.Code)
			assert.Equal(t, test.message+"\n", rec.Body.String())
		})
	}
}

func TestRedisClient(t *testing.T) {
	t.Run("UnsetEnv", func(t *testing.T) {
		got := redisClient().Options().Addr
		want := "localhost"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("SetEnv", func(t *testing.T) {
		os.Setenv("REDIS_URL", "192.168.0.1")

		defer os.Unsetenv("REDIS_URL")

		got := redisClient().Options().Addr
		want := "192.168.0.1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestConcatenateUniqueString(t *testing.T) {
	got := concatenateUniqueString("/api/get-test-api", "0")
	want := "/api/get-test-api0"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateHashString(t *testing.T) {
	got := createHashString("/api/get-test-api0")
	want := "5fec1e0adec6e3e056ede3b57b8f95b16d9d72c3"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
