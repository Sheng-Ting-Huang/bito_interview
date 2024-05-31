package storage

import (
	"sync"

	"github.com/bito_interview/model"
)

func init() {
	peopleByGender = map[model.Gender]People{}
	peopleByGender[model.GenderFemale] = People{}
	peopleByGender[model.GenderMale] = People{}
	All = personById{}
	rwMutex = &sync.RWMutex{}
}
