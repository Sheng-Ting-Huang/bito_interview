package storage

import (
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/bito_interview/model"
)

var (
	ErrNotFound = errors.New("not found")
)

var (
	peopleByGender map[model.Gender]People
	All            personById
	rwMutex        *sync.RWMutex
)

type Person struct {
	ID string
	model.Person
}
type People []*Person

var heightCmp = func(p1 *Person, p2 *Person) int {
	if p1.Height == p2.Height {
		return strings.Compare(p1.ID, p2.ID)
	}
	if p1.Height < p2.Height {
		return -1
	}
	return 1
}

func insert(people People, person *Person) People {
	index, _ := slices.BinarySearchFunc(people, person, heightCmp)
	return slices.Insert(people, index, person)
}

func remove(people People, person *Person) People {
	index, found := slices.BinarySearchFunc(people, person, heightCmp)
	if !found {
		return people
	}
	return slices.Delete(people, index, index+1)
}

func queryNMales(person *Person, n int) (People, error) {
	index, _ := slices.BinarySearchFunc(peopleByGender[model.GenderMale], &Person{Person: model.Person{PersonAttributes: model.PersonAttributes{Height: person.Height + 1}}}, heightCmp)
	if index >= len(peopleByGender[model.GenderMale]) {
		return nil, ErrNotFound
	}
	remain := min(len(peopleByGender[model.GenderMale])-index, n)
	return peopleByGender[model.GenderMale][index : index+remain], nil
}

func queryNFemales(person *Person, n int) (People, error) {
	index, _ := slices.BinarySearchFunc(peopleByGender[model.GenderFemale], &Person{Person: model.Person{PersonAttributes: model.PersonAttributes{Height: person.Height}}}, heightCmp)
	return peopleByGender[model.GenderFemale][:min(index, n)], nil
}

func queryN(person *Person, n int) (People, error) {
	if person.Gender == model.GenderFemale {
		return queryNMales(person, n)
	} else if person.Gender == model.GenderMale {
		return queryNFemales(person, n)
	}
	return nil, ErrNotFound
}

type personById map[string]*Person

func (s personById) addPersonWithId(id string, person *model.Person) *Person {
	newPerson := &Person{Person: *person, ID: id}
	s[id] = newPerson
	return newPerson
}

func (s personById) getPerson(id string) (*Person, error) {
	if person, ok := s[id]; ok {
		return person, nil
	}
	return nil, ErrNotFound
}

func (s personById) removePerson(id string) error {
	if _, ok := s[id]; ok {
		delete(s, id)
		return nil
	}
	return ErrNotFound
}

func Add(id string, person *model.Person) *Person {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	newPerson := All.addPersonWithId(id, person)
	peopleByGender[newPerson.Gender] = insert(peopleByGender[newPerson.Gender], newPerson)
	return newPerson
}

func Remove(id string) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	person, err := All.getPerson(id)
	if err == nil {
		All.removePerson(id)
		peopleByGender[person.Gender] = remove(peopleByGender[person.Gender], person)
	}
	return err
}

func Match(id string) (*Person, error) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	person, err := All.getPerson(id)
	if err != nil {
		return nil, err
	}
	possible, err := possibleMatches(id, 1)
	if err != nil {
		return nil, err
	}
	if len(possible) == 0 {
		return nil, ErrNotFound
	}

	match := possible[0]
	person.DecreaseDateCount()
	match.DecreaseDateCount()
	if person.NumberOfWantedDates <= 0 {
		peopleByGender[person.Gender] = remove(peopleByGender[person.Gender], person)
		All.removePerson(person.ID)
	}
	if match.NumberOfWantedDates <= 0 {
		peopleByGender[match.Gender] = remove(peopleByGender[match.Gender], match)
		All.removePerson(match.ID)
	}
	return match, nil
}

func PossibleMatches(id string, maxNum int) (People, error) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return possibleMatches(id, maxNum)
}

func possibleMatches(id string, maxNum int) (People, error) {
	if maxNum <= 0 {
		return nil, ErrNotFound
	}
	person, err := All.getPerson(id)
	if err != nil {
		return nil, err
	}
	if person.NumberOfWantedDates == 0 {
		return nil, ErrNotFound
	}
	matches, err := queryN(person, maxNum)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, ErrNotFound
	}
	return matches, nil
}

// ClearAll reset the storage.
func ClearAll() {
	peopleByGender = map[model.Gender]People{}
	peopleByGender[model.GenderFemale] = People{}
	peopleByGender[model.GenderMale] = People{}
	All = personById{}
}
