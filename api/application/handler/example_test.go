package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/maronfranc/poc-golang-ddd/application"
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
)

const endpoint = "/examples"

var testId string

func TestExampleCRUD(t *testing.T) {
	envfile, err := infrastructure.EnvGetFileName()
	if err != nil {
		panic(err)
	}
	err = infrastructure.EnvLoad(envfile)
	if err != nil {
		panic(err)
	}

	rootPathFile := fmt.Sprintf("../../%s", envfile)
	err = database.Start(rootPathFile)
	if err != nil {
		log.Fatalf("Test setup error: %s", err)
	}
	rt := newTestServer()

	t.Run("POST '/': should create a new example", func(t *testing.T) {
		d := dto.CreateExampleDto{
			Title:       "Test title: create",
			Description: "Test description: create"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(d)
		_, res, _ := request(rt, http.MethodPost, endpoint, b)
		var body dto.Response[dto.CreateExampleResponseDto]
		json.Unmarshal(res.Body.Bytes(), &body)
		assertEqual(t, body.Data.Description, d.Description)

		assertEqual(t, res.Code, http.StatusOK)
		testId = body.Data.Id
	})

	t.Run("PATCH '/': should update example", func(t *testing.T) {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(dto.CreateExampleDto{
			Title:       "Update test title",
			Description: "Update test description"})
		endpoint := fmt.Sprintf("%s/%s", endpoint, testId)
		_, res, _ := request(rt, http.MethodPatch, endpoint, b)
		assertEqual(t, res.Code, http.StatusOK)

		var expected dto.ResponsePaginated[dto.CreateExampleDto]
		json.Unmarshal(res.Body.Bytes(), &expected)

		for _, d := range *expected.Data {
			assertStringNotEmpty(t, "description", d.Description)
		}
	})

	t.Run("GET '/': should return paginated data", func(t *testing.T) {
		_, res, _ := request(rt, http.MethodGet, endpoint, nil)
		assertEqual(t, res.Code, http.StatusOK)

		var body dto.ResponsePaginated[dto.ManyExampleResponseDto]
		json.Unmarshal(res.Body.Bytes(), &body)

		for _, d := range *body.Data {
			assertStringNotEmpty(t, "id", d.Id)
			assertStringNotEmpty(t, "title", d.Title)
		}
	})

	t.Run("GET '/{id}': should return expected data by id", func(t *testing.T) {
		endpoint := fmt.Sprintf("%s/%s", endpoint, testId)
		_, res, _ := request(rt, http.MethodGet, endpoint, nil)
		assertEqual(t, res.Code, http.StatusOK)

		var body dto.Response[dto.CreateExampleResponseDto]
		json.Unmarshal(res.Body.Bytes(), &body)

		assertEqual(t, body.Data.Id, testId)
	})

	// t.Run("DELETE '/{id}': should delete correct example by id", func(t *testing.T) {
	// 	endpoint := fmt.Sprintf("%s/%s", ENDPOINT, testId)
	// 	_, res, _ := request(rt, http.MethodDelete, endpoint, nil)
	// 	assertEqual(t, res.Code, http.StatusOK)
	//
	// 	var body dto.Response[bool]
	// 	json.Unmarshal(res.Body.Bytes(), &body)
	//
	// 	assertEqual(t, body.Data, true)
	// })
}

func assertEqual(t testing.TB, expected, got any) {
	t.Helper()
	if reflect.DeepEqual(got, expected) {
		return
	}
	t.Errorf("expected: %s, got: %s", got, expected)
}

func assertStringNotEmpty(t testing.TB, label, v string) {
	t.Helper()
	if v != "" {
		return
	}
	t.Errorf("value '%s' is empty", label)
}

func request(
	rt *chi.Mux, method, url string, body io.Reader,
) (*http.Request, *httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, nil, err
	}
	res := httptest.NewRecorder()
	rt.ServeHTTP(res, req)
	return req, res, nil
}

func newTestServer() *chi.Mux {
	app := &application.Application{}
	rt := app.Setup()
	return rt
}
