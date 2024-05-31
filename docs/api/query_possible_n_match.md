# Query N Possible Matches

Search at most N possible matching people from the tool based on the given person.

**URL** : `/person/{id}/matches?n={n}`

**Method** : `POST`

**Auth required** : NO

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "matches": [
    {
      "id": "eda2aa1a-a61e-4ccd-a2da-a335bbfa6f51",
      "name": "abc",
      "height": 150,
      "gender": "female"
    },
    {
      "id": "eda2aa1a-a61e-4ccd-a2da-a335bbfa6f52",
      "name": "ccc",
      "height": 160,
      "gender": "female"
    }
  ]
}
```

## Error Response

**Condition** : If person cannot be found by ID or there is no match.

**Code** : `404 NOT FOUND`
