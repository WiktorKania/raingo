package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
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
		return nil, fmt.Errorf("There is no subreddit: %s", subreddit)
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
		session.ChannelMessageSend(msg.ChannelID, err.Error())
		return
	}

	imageEmbed := discordgo.MessageEmbedImage{URL: meme.ImageURLs[len(meme.ImageURLs)-1]}
	messageEmbed := discordgo.MessageEmbed{Title: meme.Title, Description: "Author: " + meme.Author, Image: &imageEmbed}

	session.ChannelMessageSendEmbed(msg.ChannelID, &messageEmbed)
}

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

func tellJoke(session *discordgo.Session, msg *discordgo.MessageCreate) {
	session.ChannelMessageSend(msg.ChannelID, "Przychodzi facet do jasnowidzki.\n- Dzień dobry, Kamilu.\n- Ale ja nie jestem Kamil.\n- Wiem.")
}

func handleMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}
	fmt.Println("Got a message, ", msg.Content)
	message := strings.ToLower(msg.Content)

	if strings.HasPrefix(message, "go ") {
		command := strings.Split(message[3:], " ")
		fmt.Println("Got command, ", command)
		switch command[0] {
		case "joke":
			tellJoke(session, msg)
		case "help":
			session.ChannelMessageSend(msg.ChannelID, "I'll look for therapy places for you in my free time")
		case "comic":
			sendComic(session, msg)
		case "meme":
			var subreddit string
			if len(command) > 1 {
				subreddit = command[1]
			}
			sendMeme(subreddit, session, msg)
		}
	}
}

func createHttpServer() {
	fmt.Println("hey", httprouter.New())
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file!")
	}
	bot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	fmt.Println("API version:", discordgo.APIVersion)
	if err != nil {
		fmt.Println("Error creating bot session!")
		panic(err)
	}
	bot.AddHandler(handleMessage)
	bot.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Close()
}
