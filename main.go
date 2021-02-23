package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

var ctx = context.Background()

type HealthStatus struct {
	Status string `json:"status"`
}

type HashId struct {
	Id string `json:"id"`
}

func main() {
	e := initEcho()

	e.Logger.Fatal(e.Start(":1323"))
}

func initEcho() *echo.Echo {
	e := echo.New()

	e.GET("/health", healthHandler)

	e.GET("/api/*", getApiHandler)

	e.POST("/api/*", postApiHandler)

	return e
}

func healthHandler(c echo.Context) error {
	healthStatus := HealthStatus{"ok"}
	return c.JSON(http.StatusOK, healthStatus)
}

func getApiHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot read request body")
	}

	requestBodyString := ""
	if len(body) != 0 {
		var requestBodyInterface interface{}
		err = json.Unmarshal([]byte(body), &requestBodyInterface)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot unmarshal request body")
		}

		requestBodyMap := requestBodyInterface.(map[string]interface{})

		r, err := json.Marshal(requestBodyMap)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot marshal request body")
		}

		requestBodyString = string(r)
	}

	uniqueString := concatenateUniqueString(c.Request().URL.String(), string(requestBodyString))

	hashIdString := createHashString(uniqueString)

	response, err := redisClient().Get(ctx, hashIdString).Result()
	if err == redis.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "data needs to be set before GET")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot connect to database")
	}

	encodedResponse := []byte(response)
	return c.JSONBlob(http.StatusOK, encodedResponse)
}

func postApiHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot read request body")
	}

	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "there is no request body")
	}

	var requestBodyInterface interface{}
	err = json.Unmarshal([]byte(body), &requestBodyInterface)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot unmarshal request body")
	}

	requestBodyMap := requestBodyInterface.(map[string]interface{})

	if requestBodyMap["response"] == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing request response")
	}

	requestBodyString := ""
	if requestBodyMap["body"] != nil {
		r, err := json.Marshal(requestBodyMap["body"])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot marshal request body")
		}
		requestBodyString = string(r)
	}

	requestResponseString, err := json.Marshal(requestBodyMap["response"])
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot marshal request response")
	}

	uniqueString := concatenateUniqueString(c.Request().URL.String(), requestBodyString)

	hashIdString := createHashString(uniqueString)

	err = redisClient().Set(ctx, hashIdString, requestResponseString, 7*24*time.Hour).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot connect to database")
	}

	hashId := HashId{hashIdString}

	return c.JSON(http.StatusOK, hashId)
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
		os.Setenv("REDIS_URL", "localhost")
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
