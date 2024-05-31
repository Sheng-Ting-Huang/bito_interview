package storage

import (
	"testing"

	"github.com/bito_interview/model"
	"github.com/google/go-cmp/cmp"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		ID          string
		person      *model.Person
		addedPerson *Person
	}{
		{
			name: "success",
			ID:   "id-1",
			person: &model.Person{
				PersonAttributes: model.PersonAttributes{
					Height: 1,
					Gender: model.GenderFemale,
					Name:   "abc",
				},
			},
			addedPerson: &Person{
				ID: "id-1",
				Person: model.Person{
					PersonAttributes: model.PersonAttributes{
						Height: 1,
						Gender: model.GenderFemale,
						Name:   "abc",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			teardown := setupTest(t)
			defer teardown(t)
			got := Add(test.ID, test.person)
			if diff := cmp.Diff(got, test.addedPerson); diff != "" {
				t.Errorf("%s got want:\n%s", t.Name(), diff)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		expectedErr error
		numPeople   int
		numMales    int
		numFemales  int
	}{
		{
			name:        "not_found",
			id:          "not_found",
			expectedErr: ErrNotFound,
			numPeople:   2,
			numMales:    1,
			numFemales:  1,
		},
		{
			name:       "success",
			id:         "id-1",
			numPeople:  1,
			numMales:   0,
			numFemales: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			teardown := setupTest(
				t,
				createPerson("id-1", model.GenderMale, 10, 1),
				createPerson("id-2", model.GenderFemale, 10, 1),
			)
			defer teardown(t)
			if got, want := Remove(test.id), test.expectedErr; got != want {
				t.Errorf("%s got %v but want: %v", t.Name(), got, want)
			}
			if got, want := len(All), test.numPeople; got != want {
				t.Errorf("%s got %v people but want: %v", t.Name(), got, want)
			}
			if got, want := len(peopleByGender[model.GenderMale]), test.numMales; got != want {
				t.Errorf("%s got %v males but want: %v", t.Name(), got, want)
			}
			if got, want := len(peopleByGender[model.GenderFemale]), test.numFemales; got != want {
				t.Errorf("%s got %v females but want: %v", t.Name(), got, want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name   string
		people People
		id     string
		match  *Person
		err    error
	}{
		{
			name: "person_not_found",
			id:   "not_found",
			err:  ErrNotFound,
		},
		{
			name: "success",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 1, 1),
			},
			match: createPerson("2", model.GenderFemale, 1, 0),
		},
		{
			name: "male_no_match",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "girl_no_match",
			id:   "2",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "no_more_dates",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 0),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "match_2",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 9, 1),
				createPerson("3", model.GenderFemale, 9, 0),
			},
			match: createPerson("2", model.GenderFemale, 9, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			teardown := setupTest(t, test.people...)
			defer teardown(t)
			got, gotErr := Match(test.id)
			if diff := cmp.Diff(got, test.match); diff != "" {
				t.Errorf("%s got want:\n%s", t.Name(), diff)
			}
			if got, want := gotErr, test.err; got != want {
				t.Errorf("%s got %v but want: %v", t.Name(), got, want)
			}
		})
	}
}

func TestPossibleMatches(t *testing.T) {
	tests := []struct {
		name    string
		people  People
		id      string
		n       int
		matches People
		err     error
	}{
		{
			name: "person_not_found",
			id:   "not_found",
			err:  ErrNotFound,
		},
		{
			name: "success",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 1, 1),
			},
			n: 1,
			matches: People{
				createPerson("2", model.GenderFemale, 1, 1),
			},
		},
		{
			name: "male_no_match",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "girl_no_match",
			id:   "2",
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "no_more_dates",
			id:   "1",
			people: People{
				createPerson("1", model.GenderMale, 10, 0),
				createPerson("2", model.GenderFemale, 11, 1),
			},
			err: ErrNotFound,
		},
		{
			name: "matches_3females",
			id:   "1",
			n:    3,
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 9, 1),
				createPerson("3", model.GenderFemale, 1, 1),
				createPerson("4", model.GenderFemale, 2, 1),
				createPerson("5", model.GenderFemale, 3, 1),
				createPerson("6", model.GenderFemale, 4, 1),
				createPerson("7", model.GenderFemale, 5, 1),
			},
			matches: People{
				createPerson("3", model.GenderFemale, 1, 1),
				createPerson("4", model.GenderFemale, 2, 1),
				createPerson("5", model.GenderFemale, 3, 1),
			},
		},
		{
			name: "matches_3males",
			id:   "1",
			n:    3,
			people: People{
				createPerson("1", model.GenderFemale, 1, 1),
				createPerson("2", model.GenderMale, 9, 1),
				createPerson("3", model.GenderMale, 1, 1),
				createPerson("4", model.GenderMale, 2, 1),
				createPerson("5", model.GenderMale, 3, 1),
				createPerson("6", model.GenderMale, 4, 1),
				createPerson("7", model.GenderMale, 5, 1),
			},
			matches: People{
				createPerson("4", model.GenderMale, 2, 1),
				createPerson("5", model.GenderMale, 3, 1),
				createPerson("6", model.GenderMale, 4, 1),
			},
		},
		{
			name: "ask_more_than_available_males",
			id:   "1",
			n:    3,
			people: People{
				createPerson("1", model.GenderFemale, 1, 1),
				createPerson("2", model.GenderMale, 9, 1),
				createPerson("3", model.GenderMale, 2, 1),
			},
			matches: People{
				createPerson("3", model.GenderMale, 2, 1),
				createPerson("2", model.GenderMale, 9, 1),
			},
		},
		{
			name: "ask_more_than_available_females",
			id:   "1",
			n:    3,
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
				createPerson("2", model.GenderFemale, 9, 1),
				createPerson("3", model.GenderFemale, 2, 1),
			},
			matches: People{
				createPerson("3", model.GenderFemale, 2, 1),
				createPerson("2", model.GenderFemale, 9, 1),
			},
		},
		{
			name: "negativeN",
			id:   "1",
			n:    -1,
			people: People{
				createPerson("1", model.GenderMale, 10, 1),
			},
			err: ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			teardown := setupTest(t, test.people...)
			defer teardown(t)
			got, gotErr := PossibleMatches(test.id, test.n)
			if diff := cmp.Diff(got, test.matches); diff != "" {
				t.Errorf("%s got want:\n%s", t.Name(), diff)
			}
			if got, want := gotErr, test.err; got != want {
				t.Errorf("%s got %v but want: %v", t.Name(), got, want)
			}
		})
	}
}

func setupTest(tb testing.TB, people ...*Person) func(tb testing.TB) {
	setupPeople(tb, people...)
	return func(tb testing.TB) {
		peopleByGender = map[model.Gender]People{}
		peopleByGender[model.GenderFemale] = People{}
		peopleByGender[model.GenderMale] = People{}
		All = personById{}
	}
}
func setupPeople(tb testing.TB, people ...*Person) {
	tb.Helper()
	for _, person := range people {
		Add(person.ID, &person.Person)
	}
}

func createPerson(id string, gender model.Gender, height int, numDates int) *Person {
	return &Person{ID: id, Person: model.Person{PersonAttributes: model.PersonAttributes{
		Height: height,
		Gender: gender},
		NumberOfWantedDates: numDates,
	}}
}
