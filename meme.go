package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Meme struct {
	Title     string
	NSFW      bool
	Author    string
	ImageURLs []string `json:"preview"` // lowest to highest quality
}

func fetchMeme(memeURL string) (*Meme, error) {
	res, err := http.Get(memeURL)
	if err != nil {
		log.Println("Couldn't reach meme-api: ", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		log.Println("Couldn't reach meme-api subreddit: ", memeURL)
		splittedURL := strings.Split(memeURL, "/")
		subreddit := splittedURL[len(splittedURL)-1]
		return nil, fmt.Errorf("there is no subreddit: %s", subreddit)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Couldn't read meme-api: ", err)
		return nil, err
	}

	var meme Meme
	if err := json.Unmarshal(bodyBytes, &meme); err != nil {
		log.Println("Couldn't unmarshall meme: ", err)
		return nil, err
	}

	return &meme, nil
}

func sendMeme(subreddit string, session *discordgo.Session, msg *discordgo.MessageCreate) {
	memeURL := "https://meme-api.herokuapp.com/gimme/" + subreddit
	meme, err := fetchMeme(memeURL)
	if err != nil {
		log.Println("Couldn't fetch meme: ", err)
		replyToChannel(msg.ChannelID, err.Error())
		return
	}

	imageEmbed := discordgo.MessageEmbedImage{URL: meme.ImageURLs[len(meme.ImageURLs)-1]}
	messageEmbed := discordgo.MessageEmbed{Title: meme.Title, Description: "Author: " + meme.Author, Image: &imageEmbed}

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}
