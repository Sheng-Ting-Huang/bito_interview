package api

import "github.com/bito_interview/model"

type AddAndMatchResponse struct {
	Self  *PersonResponse `json:"self"`
	Match *PersonResponse `json:"match"`
}

type PersonResponse struct {
	ID string `json:"id"`
	model.PersonAttributes
}

type PossibleMatches struct {
	Matches []PersonResponse `json:"matches"`
}
