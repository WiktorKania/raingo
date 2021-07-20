package main

import (
	"bytes"
	"encoding/json"
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

var SpartathlonID int = 865167211944345600
var session *discordgo.Session

func tellJoke(session *discordgo.Session, msg *discordgo.MessageCreate) {
	session.ChannelMessageSend(msg.ChannelID, "Przychodzi facet do jasnowidzki.\n- Dzień dobry, Kamilu.\n- Ale ja nie jestem Kamil.\n- Wiem.")
}

func makeItRain(username string, message string, channelID string) {
	channel, err := session.Channel(channelID)
	if err != nil {
		log.Println("Couldn't find channel")
		session.ChannelMessageSend(channelID, "I can't, I don't know where we are.")
		return
	}
	values := map[string]string{"Nickname": username, "Msg": message, "Channel": channel.Name}
	json, err := json.Marshal(values)

	if err != nil {
		log.Println(err)
		session.ChannelMessageSend(channelID, "I can't, I can't write it down.")
	}
	r, err := http.Post("http://localhost:8080/api/raino", "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Println(err)
		session.ChannelMessageSend(channelID, "I tried, I couldn't find them.")
	}
	session.ChannelMessageSend(channelID, "I tossed the message!")
	fmt.Println(r)
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
		case "hermes":
			var message string
			if len(command) > 1 {
				message = strings.Join(command[1:], " ")
			} else {
				session.ChannelMessageSend(msg.ChannelID, "What do you want me to say?")
				return
			}
			makeItRain(msg.Author.Username, message, msg.ChannelID)
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
