package model

type Gender string

const (
	GenderFemale Gender = "female"
	GenderMale   Gender = "male"
)

type PersonAttributes struct {
	Name   string `json:"name" validate:"required"`
	Height int    `json:"height" validate:"gt=0,lte=250"`
	Gender Gender `json:"gender" validate:"oneof=female male"`
}

type Person struct {
	PersonAttributes
	NumberOfWantedDates int `json:"number_of_wanted_dates" validate:"gt=0"`
}

func (p *Person) DecreaseDateCount() {
	p.NumberOfWantedDates--
}
