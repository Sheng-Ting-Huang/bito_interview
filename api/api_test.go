package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bito_interview/model"
	"github.com/bito_interview/storage"
)

func TestAddSinglePersonAndMatch(t *testing.T) {
	IdGenerator = storage.FakeIDGenerator{
		FakeID: "1",
	}
	tests := []struct {
		name       string
		person     *storage.Person
		req        *http.Request
		statusCode int
		respBody   *string
	}{
		{
			name:       "no request json body",
			req:        newRequest(http.MethodPost, "/add-and-match", nil),
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid req body",
			req:        newRequest(http.MethodPost, "/add-and-match", bytes.NewBuffer(jsonMarshal(t, &model.Person{}))),
			statusCode: http.StatusBadRequest,
		},
		{
			name: "person created no match",
			req: newRequest(http.MethodPost, "/add-and-match", bytes.NewBuffer(jsonMarshal(t, &model.Person{
				PersonAttributes:    model.PersonAttributes{Name: "abc", Height: 100, Gender: model.GenderMale},
				NumberOfWantedDates: 10,
			}))),
			statusCode: http.StatusOK,
			respBody: func() *string {
				a := `{"self":{"id":"1","name":"abc","height":100,"gender":"male"},"match":null}`
				return &a
			}(),
		},
		{
			name: "person created no match",
			req: newRequest(http.MethodPost, "/add-and-match", bytes.NewBuffer(jsonMarshal(t, &model.Person{
				PersonAttributes:    model.PersonAttributes{Name: "abc", Height: 100, Gender: model.GenderMale},
				NumberOfWantedDates: 10,
			}))),
			statusCode: http.StatusOK,
			respBody: func() *string {
				a := `{"self":{"id":"1","name":"abc","height":100,"gender":"male"},"match":null}`
				return &a
			}(),
		},
		{
			name:   "person created with match",
			person: createPerson("2", model.GenderFemale, 90, 1),
			req: newRequest(http.MethodPost, "/add-and-match", bytes.NewBuffer(jsonMarshal(t, &model.Person{
				PersonAttributes:    model.PersonAttributes{Name: "abc", Height: 100, Gender: model.GenderMale},
				NumberOfWantedDates: 10,
			}))),
			statusCode: http.StatusOK,
			respBody: func() *string {
				a := `{"self":{"id":"1","name":"abc","height":100,"gender":"male"},"match":{"id":"2","name":"","height":90,"gender":"female"}}`
				return &a
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var teardown func(tb testing.TB)
			if test.person != nil {
				teardown = setupTest(t, test.person)
			} else {
				teardown = setupTest(t)
			}
			defer teardown(t)
			rec := executeRequest(t, test.req)
			if got, want := rec.Code, test.statusCode; got != want {
				t.Errorf("%s got status code %v but want %v", t.Name(), got, want)
			}
			if test.respBody != nil {
				if got, want := rec.Body.String(), *test.respBody; got != want {
					t.Errorf("%s got response body %v but want %v", t.Name(), got, want)
				}
			}
		})
	}
}

func TestRemoveSinglePerson(t *testing.T) {
	tests := []struct {
		name       string
		person     *storage.Person
		req        *http.Request
		statusCode int
		respBody   *string
	}{
		{
			name:       "path has no id",
			req:        newRequest(http.MethodDelete, "/person", nil),
			statusCode: http.StatusNotFound,
			respBody: func() *string {
				s := "404 page not found\n"
				return &s
			}(),
		},
		{
			name:       "person not found",
			req:        newRequest(http.MethodDelete, "/person/1", nil),
			statusCode: http.StatusNotFound,
			respBody: func() *string {
				s := "not found\n"
				return &s
			}(),
		},
		{
			name:       "person deleted",
			person:     createPerson("1", model.GenderMale, 1, 1),
			req:        newRequest(http.MethodDelete, "/person/1", nil),
			statusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var teardown func(tb testing.TB)
			if test.person != nil {
				teardown = setupTest(t, test.person)
			} else {
				teardown = setupTest(t)
			}
			defer teardown(t)
			rec := executeRequest(t, test.req)
			if got, want := rec.Code, test.statusCode; got != want {
				t.Errorf("%s got status code %v but want %v", t.Name(), got, want)
			}
			if test.respBody != nil {
				if got, want := rec.Body.String(), *test.respBody; got != want {
					t.Errorf("%s got response body %v but want %v", t.Name(), got, want)
				}
			}
		})
	}
}

func TestAddSinglePeople(t *testing.T) {
	IdGenerator = storage.FakeIDGenerator{
		FakeID: "1",
	}
	tests := []struct {
		name       string
		people     storage.People
		req        *http.Request
		statusCode int
		respBody   *string
	}{
		{
			name:       "no query parameter",
			req:        newRequest(http.MethodGet, "/person/1/matches", nil),
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid query parameter value",
			req:        newRequest(http.MethodGet, "/person/1/matches?n=-1", nil),
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "person not found",
			req:        newRequest(http.MethodGet, "/person/1/matches?n=1", nil),
			statusCode: http.StatusNotFound,
		},
		{
			name: "no match",
			people: storage.People{
				createPerson("1", model.GenderFemale, 90, 1),
				createPerson("2", model.GenderFemale, 90, 1),
			},
			req:        newRequest(http.MethodGet, "/person/1/matches?n=1", nil),
			statusCode: http.StatusNotFound,
		},
		{
			name: "match",
			people: storage.People{
				createPerson("1", model.GenderFemale, 90, 1),
				createPerson("2", model.GenderMale, 100, 1),
			},
			req:        newRequest(http.MethodGet, "/person/1/matches?n=1", nil),
			statusCode: http.StatusOK,
			respBody: func() *string {
				s := `{"matches":[{"id":"2","name":"","height":100,"gender":"male"}]}`
				return &s
			}(),
		},
		{
			name: "match many",
			people: storage.People{
				createPerson("1", model.GenderFemale, 90, 1),
				createPerson("2", model.GenderMale, 100, 1),
				createPerson("3", model.GenderMale, 100, 1),
				createPerson("4", model.GenderMale, 100, 1),
				createPerson("5", model.GenderMale, 100, 1),
			},
			req:        newRequest(http.MethodGet, "/person/1/matches?n=3", nil),
			statusCode: http.StatusOK,
			respBody: func() *string {
				s := `{"matches":[{"id":"2","name":"","height":100,"gender":"male"},{"id":"3","name":"","height":100,"gender":"male"},{"id":"4","name":"","height":100,"gender":"male"}]}`
				return &s
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var teardown func(tb testing.TB)
			if len(test.people) > 0 {
				teardown = setupTest(t, test.people...)
			} else {
				teardown = setupTest(t)
			}
			defer teardown(t)
			rec := executeRequest(t, test.req)
			if got, want := rec.Code, test.statusCode; got != want {
				t.Errorf("%s got status code %v but want %v", t.Name(), got, want)
			}
			if test.respBody != nil {
				if got, want := rec.Body.String(), *test.respBody; got != want {
					t.Errorf("%s got response body %v but want %v", t.Name(), got, want)
				}
			}
		})
	}
}

func executeRequest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	router := NewRouter()
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

func newRequest(method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	return req
}

func jsonMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func setupTest(tb testing.TB, people ...*storage.Person) func(tb testing.TB) {
	setupPeople(tb, people...)
	return func(tb testing.TB) {
		storage.ClearAll()
	}
}
func setupPeople(tb testing.TB, people ...*storage.Person) {
	tb.Helper()
	for _, person := range people {
		storage.Add(person.ID, &person.Person)
	}
}

func createPerson(id string, gender model.Gender, height int, numDates int) *storage.Person {
	return &storage.Person{ID: id, Person: model.Person{PersonAttributes: model.PersonAttributes{
		Height: height,
		Gender: gender},
		NumberOfWantedDates: numDates,
	}}
}
