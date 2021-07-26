package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

var (
	SpartathlonID int = 865167211944345600
	session       *discordgo.Session
)

const (
	boomerJoke = "Przychodzi facet do jasnowidzki.\n- Dzień dobry, Kamilu.\n- Ale ja nie jestem Kamil.\n- Wiem."
)

func tellJoke(session *discordgo.Session, msg *discordgo.MessageCreate, joke string) {
	session.ChannelMessageSend(msg.ChannelID, joke)
}

func replyToChannel(channelID string, msg string) {
	session.ChannelMessageSend(channelID, msg)
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
			if len(command) == 2 && command[1] == "boomer" {
				tellJoke(session, msg, boomerJoke)
				break
			}
			joke, err := fetchJoke()
			if err != nil {
				replyToChannel(msg.ChannelID, "Joke failed")
				fmt.Println(err)
				break
			}
			tellJoke(session, msg, joke)
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
	router := httprouter.New()
	router.POST("/api/raino", listenToRaindrops)
	router.GET("/api/raino", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("get request")
		rw.Write([]byte("Hello"))
	})
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	godotenv.Load(".env")
	botToken, present := os.LookupEnv("BOT_TOKEN")
	if !present {
		panic("No bot token found!")
	}
	bot, err := discordgo.New("Bot " + botToken)
	fmt.Println("API version:", discordgo.APIVersion)
	if err != nil {
		fmt.Println("Error creating bot session!")
		panic(err)
	}
	bot.AddHandler(handleMessage)
	bot.Open()
	session = bot
	createHttpServer()
	bot.Close()
}
