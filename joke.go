package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	anyJokeUrl = "https://v2.jokeapi.dev/joke/Any"
)

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

func fetchJoke() (string, error) {
	res, err := http.Get(anyJokeUrl)
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
