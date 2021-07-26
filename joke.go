package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	jokeApiUrl = "https://v2.jokeapi.dev/joke/"
)

var (
	allowedJokeCategories = [...]string{
		"Any", "Programming", "Miscellaneous", "Dark", "Pun", "Spooky", "Christmas",
	}
	jokeCategoriesMap = map[string]string{
		"code":   "Programming",
		"any":    "Any",
		"misc":   "Miscellaneous",
		"dark":   "Dark",
		"pun":    "Pun",
		"spooky": "Spooky",
		"xmas":   "Christmas",
	}
)

func validJokeCategory(jokeType string) bool {
	for _, _type := range allowedJokeCategories {
		if jokeType == _type {
			return true
		}
	}
	_, exists := jokeCategoriesMap[jokeType]
	return exists
}

type JokeFlags struct {
	Nsfw      bool
	Religious bool
	Political bool
	Racist    bool
	Sexist    bool
	Explicit  bool
}

type JokeResponse struct {
	Error    bool
	Category string
	Type     string
	Setup    string
	Delivery string
	Joke     string
	Flags    JokeFlags
	Id       int
	Safe     bool
	Lang     string
}

func fetchJoke(jokeType string) (string, error) {
	if !validJokeCategory(jokeType) {
		return "", errors.New("incorrect joke category")
	}
	category, ok := jokeCategoriesMap[jokeType]
	if ok {
		jokeType = category
	}
	res, err := http.Get(jokeApiUrl + jokeType)
	if err != nil {
		log.Println("Couldn't reach joke: ", err)
		return "", err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Couldn't read joke: ", err)
		return "", err
	}
	var jokeRes JokeResponse
	if err := json.Unmarshal(bodyBytes, &jokeRes); err != nil {
		log.Println("Couldn't unmarshall joke: ", err)
		return "", err
	}
	if jokeRes.Type == "single" {
		return jokeRes.Joke, nil
	} else {
		fullJoke := jokeRes.Setup + "\n\n" + jokeRes.Delivery
		return fullJoke, nil
	}
}
