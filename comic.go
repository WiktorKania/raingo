package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Comic struct {
	Num        int
	Title      string
	Transcript string
	ImageURL   string `json:"img"`
}

func fetchComic(comicURL string) (*Comic, error) {
	res, err := http.Get(comicURL)
	if err != nil {
		log.Println("Couldn't reach xkcd: ", err)
		return nil, err
	}

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Couldn't read xkcd: ", err)
		return nil, err
	}

	var comic Comic
	if err := json.Unmarshal(bodyBytes, &comic); err != nil {
		log.Println("Couldn't unmarshall comic: ", err)
		return nil, err
	}

	return &comic, nil
}

func sendComic(session *discordgo.Session, msg *discordgo.MessageCreate) {
	baseURL := "https://xkcd.com"
	suffixURL := "info.0.json"
	newestURL := fmt.Sprintf("%s/%s", baseURL, suffixURL)
	newestComic, err := fetchComic(newestURL)
	if err != nil {
		log.Println("Couldn't fetch comic: ", err)
	}
	maxNum := newestComic.Num

	randomNum := rand.Intn(maxNum) + 1
	randomURL := fmt.Sprintf("%s/%d/%s", baseURL, randomNum, suffixURL)

	randomComic, err := fetchComic(randomURL)
	if err != nil {
		log.Println("Couldn't fetch comic: ", err)
	}

	imageEmbed := discordgo.MessageEmbedImage{URL: randomComic.ImageURL}
	messageEmbed := discordgo.MessageEmbed{Title: randomComic.Title, Image: &imageEmbed}

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
