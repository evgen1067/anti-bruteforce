package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/evgen1067/anti-bruteforce/internal/repository/psql"
	"github.com/evgen1067/anti-bruteforce/internal/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	blacklistURL = "/list/blacklist"
	whitelistURL = "/list/whitelist"
)

func TestCustomNotFoundHandler(t *testing.T) {
	t.Run("custom not found handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		CustomNotFoundHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		var ex common.APIException
		err = ex.UnmarshalJSON(data)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, ex.Code)
		require.Equal(t, "The page you are looking for does not exist.", ex.Message)
	})
}

func TestAdd(t *testing.T) {
	cfg, err := config.Parse("../../configs/local.json")
	require.NoError(t, err)
	ctx := context.Background()
	repo := psql.NewRepo(cfg)
	err = repo.Connect(ctx)
	require.NoError(t, err)
	defer repo.Close()
	s := service.NewServices(ctx, repo)
	NewServer(s, cfg)

	request := common.APIListRequest{Address: "127.0.12.1/25"}

	t.Run("Successful addition to the blacklist", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Add).Methods(http.MethodPost)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, blacklistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("Unsuccessful addition to the blacklist (the address already exists in the list)", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Add).Methods(http.MethodPost)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, blacklistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusConflict, recorder.Code)
	})

	t.Run("Successful addition to the whitelist", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Add).Methods(http.MethodPost)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, whitelistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("Unsuccessful addition to the blacklist (the address already exists in the list)", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Add).Methods(http.MethodPost)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, whitelistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusConflict, recorder.Code)
	})

	t.Run("Successful removal from the blacklist", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Delete).Methods(http.MethodDelete)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, blacklistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusAccepted, recorder.Code)
	})

	t.Run("Successful removal from the whitelist", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Delete).Methods(http.MethodDelete)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, whitelistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusAccepted, recorder.Code)
	})

	t.Run("Unsuccessful removal from the blacklist (address does not exist)", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Delete).Methods(http.MethodDelete)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, blacklistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("Unsuccessful removal from the whitelist (address does not exist)", func(t *testing.T) {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(&request)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/list/{value}", Delete).Methods(http.MethodDelete)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, whitelistURL, b)
		require.NoError(t, err)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
