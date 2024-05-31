package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bito_interview/model"
	storage "github.com/bito_interview/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var validate = validator.New()

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/add-and-match", AddSinglePersonAndMatch).Methods(http.MethodPost)
	router.HandleFunc("/person/{id}", RemoveSinglePerson).Methods(http.MethodDelete)
	router.HandleFunc("/person/{id}/matches", QuerySinglePeople).Methods(http.MethodGet)
	router.HandleFunc("/swagger/{any}", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/v1/swagger/doc.json"))).Methods(http.MethodGet)
	return router
}

func AddSinglePersonAndMatch(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "request json body missing", http.StatusBadRequest)
		return
	}
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPerson := &model.Person{}
	if err := json.Unmarshal(jsonBytes, &newPerson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(newPerson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storagePerson := storage.Add(IdGenerator.GenerateKey(), newPerson)

	resp := AddAndMatchResponse{Self: &PersonResponse{ID: storagePerson.ID, PersonAttributes: storagePerson.PersonAttributes}}
	matchPerson, err := storage.Match(storagePerson.ID)
	if err == nil {
		resp.Match = &PersonResponse{ID: matchPerson.ID, PersonAttributes: matchPerson.PersonAttributes}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResp))
}

func RemoveSinglePerson(w http.ResponseWriter, r *http.Request) {
	var id string
	var ok bool
	if id, ok = mux.Vars(r)["id"]; !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	err := storage.Remove(id)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprint(w, "removed")
}

func QuerySinglePeople(w http.ResponseWriter, r *http.Request) {
	var id string
	var ok bool
	if id, ok = mux.Vars(r)["id"]; !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if !r.URL.Query().Has("n") {
		http.Error(w, "query parameter `n` is required", http.StatusBadRequest)
		return
	}
	maxNum, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		http.Error(w, "query parameter `n` number is expected", http.StatusBadRequest)
		return
	}
	if maxNum <= 0 {
		http.Error(w, "query parameter `n` positive number is expected", http.StatusBadRequest)
		return
	}

	matches, err := storage.PossibleMatches(id, maxNum)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	resp := &PossibleMatches{}
	for _, match := range matches {
		resp.Matches = append(resp.Matches, PersonResponse{ID: match.ID, PersonAttributes: match.PersonAttributes})
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResp))
}
