# api-imbgo
API that allows consumers to access the movies data. Most of the queries will be
against local database but you can use imdb npm package
(​ https://godoc.org/github.com/eefret/gomdb​ ) to query for movies when data isn’t available in the
local database and store it for future reference. Movie object in database will have following
properties as well as movie object in response of the apis:
- title
- released year
- rating
- id
- genres (array of strings)

## Language

I've used Golang because I want learn it.

## Installation

I assume that you have Golang installed in your host.

* Clone repository
```
git clone  https://github.com/lexisvar/api-imbgo.git
```
* Get all required packages
```
go get
```
* Run go server (please go to project root)

```
go run .
```
## Data Storage

Data is stored in memory, data is lost during server restarts.

## API Definition

* Returns a list of movies
```
http://{{server}}/movies/

{
  "imdbid": "tt0452694",
  "title": "The Time Traveler's Wife",
  "released": "14 Aug 2009",
  "rating": "7.1",
  "genres": [
      "Comedy",
      "Drama",
      "Fantasy"
  ]
},
```
* Returns a list of movies by Id
```
http://{{server}}/movies/tt0328107

{
  "imdbid": "tt0328107",
  "title": "Man on Fire",
  "released": "23 Apr 2004",
  "rating": "7.7",
  "genres": [
      "Action",
      " Crime",
      " Drama"
  ]
}
```

* Find movie in local storage if is not found then go to de imb API and save
```
http://{{server}}/movies/title/joker

{
  "imdbid": "tt7286456",
  "title": "Joker",
  "released": "04 Oct 2019",
  "rating": "8.4",
  "genres": [
      "Crime",
      "Drama",
      "Thriller"
  ]
}
```

* Filter by many fields
```
http://{{server}}/movies/filter?start_release=1995&end_release=2020&genres=Horror&rating=6

[
  {
      "imdbid": "tt0452694",
      "title": "The Time Traveler's Wife",
      "released": "14 Aug 2009",
      "rating": "7.1",
      "genres": [
          "Comedy",
          "Drama",
          "Fantasy"
      ]
  }
]
```

* Delete a movie by id from local storage
```
http://{{server}}/movies/tt2884018
```
