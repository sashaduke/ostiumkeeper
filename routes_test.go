package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandleData(t *testing.T) {
	// Setup Redis client with test data
	testData := Data{Timestamp: time.Now().UTC(), Value: "0.12618"}
	storeDataRedis(testData)

	req, _ := http.NewRequest("GET", "/data", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleData)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandleContracts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/contracts", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleContracts)

	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
