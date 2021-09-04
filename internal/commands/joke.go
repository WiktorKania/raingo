package commands

import (
	"encoding/json"
	"errors"
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

type JokeResponse struct {
	Error    bool
	Type     string
	Setup    string
	Delivery string
	Joke     string
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
	var jokeRes JokeResponse
	if err := json.NewDecoder(res.Body).Decode(&jokeRes); err != nil {
		log.Println("Couldn't decode joke: ", err)
		return "", err
	}
	if jokeRes.Type == "single" {
		return jokeRes.Joke, nil
	} else {
		fullJoke := jokeRes.Setup + "\n\n" + jokeRes.Delivery
		return fullJoke, nil
	}
}
