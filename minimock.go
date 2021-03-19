package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	setupPort()

	e := initEcho()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func initEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	e.GET("/", healthHandler)
	e.GET("/api/*", getApiHandler)
	e.POST("/api/*", postApiHandler)
	e.PUT("/api/*", putApiHandler)

	return e
}

func setupPort() {
	_, ok := os.LookupEnv("PORT")
	if !ok {
		os.Setenv("PORT", "1323")
	}
}

func authKey() bool {
	_, ok := os.LookupEnv("AUTH_KEY")
	if !ok {
		return false
	} else {
		return true
	}
}

func unauthorized(authorization string) bool {
	splitToken := strings.Split(authorization, "Bearer")
	if len(splitToken) != 2 {
		return true
	}

	authKey := strings.TrimSpace(splitToken[1])

	if authKey != os.Getenv("AUTH_KEY") {
		return true
	}

	return false
}

func healthHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusOK)
}

func getApiHandler(c echo.Context) error {
	hashIdString := createHashString(c.Request().URL.String())

	response, err := redisClient().Get(context.Background(), hashIdString).Result()
	if err == redis.Nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	encodedResponse := []byte(response)

	return c.JSONBlob(http.StatusOK, encodedResponse)
}

func postApiHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	requestBodyString := ""
	if len(body) != 0 {
		var requestBodyInterface interface{}
		err = json.Unmarshal([]byte(body), &requestBodyInterface)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		requestBodyMap := requestBodyInterface.(map[string]interface{})

		r, err := json.Marshal(requestBodyMap)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		requestBodyString = string(r)
	}

	uniqueString := concatenateUniqueString(c.Request().URL.String(), string(requestBodyString))

	hashIdString := createHashString(uniqueString)

	response, err := redisClient().Get(context.Background(), hashIdString).Result()
	if err == redis.Nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	encodedResponse := []byte(response)

	return c.JSONBlob(http.StatusOK, encodedResponse)
}

func putApiHandler(c echo.Context) error {
	if authKey() {
		authorization := c.Request().Header.Get("Authorization")
		if unauthorized(authorization) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var requestBodyInterface interface{}
	err = json.Unmarshal([]byte(body), &requestBodyInterface)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	requestBodyMap := requestBodyInterface.(map[string]interface{})

	if requestBodyMap["response"] == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	requestBodyString := ""
	if requestBodyMap["body"] != nil {
		r, err := json.Marshal(requestBodyMap["body"])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		requestBodyString = string(r)
	}

	requestResponseString, err := json.Marshal(requestBodyMap["response"])
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	uniqueString := concatenateUniqueString(c.Request().URL.String(), requestBodyString)

	hashIdString := createHashString(uniqueString)

	err = redisClient().Set(context.Background(), hashIdString, requestResponseString, 7*24*time.Hour).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK)
}

func concatenateUniqueString(urlString string, bodyString string) string {
	uniqueString := urlString + bodyString

	return uniqueString
}

func createHashString(uniqueString string) string {
	sha1Hash := sha1.Sum([]byte(uniqueString))
	hashIdString := hex.EncodeToString(sha1Hash[:])

	return hashIdString
}

func redisClient() *redis.Client {
	_, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		os.Setenv("REDIS_URL", "localhost:6379")
	}

	_, ok = os.LookupEnv("REDIS_PASS")
	if !ok {
		os.Setenv("REDIS_PASS", "")
	}
	address := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASS")
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return client
}
