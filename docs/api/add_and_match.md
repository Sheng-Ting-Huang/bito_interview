# Add and Match

Add the given person and find a matching person.

**URL** : `/add-and-match/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
    "name": "person name",
    "height": 100, //between 1 to 250
    "gender": "male", // male or female only
    "number_of_wanted_dates": 1 // any integer number greater than 0
}
```

**Data example**

```json
{
    "name": "Jason",
    "height": 180,
    "gender": "male",
    "number_of_wanted_dates": 10
}
```

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "self": {
    "id": "ec6cf230-a113-4102-b3e2-b335391a8304",
    "name":"abc",
    "height":100,
    "gender":"male"
  },
  "match": {
    "id":"eda2aa1a-a61e-4ccd-a2da-a335bbfa6f51",
    "name":"abc",
    "height":9,
    "gender":"female"
  }
}
```

## Error Response

**Condition** : If person information is invalid.

**Code** : `400 BAD REQUEST`
