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

	"github.com/WiktorKania/raingo/internal/commands"
	"github.com/WiktorKania/raingo/internal/raingo"
	"github.com/WiktorKania/raingo/internal/utils"
)

const (
	boomerJoke = "Przychodzi facet do jasnowidzki.\n- Dzień dobry, Kamilu.\n- Ale ja nie jestem Kamil.\n- Wiem."
)

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
			fetchAndRespond := func(jokeType string) {
				joke, err := commands.FetchJoke(jokeType)
				if err != nil {
					utils.ReplyToChannel(msg.ChannelID, "Joke failed")
					log.Println(err)
					return
				}
				utils.ReplyToChannel(msg.ChannelID, joke)
			}
			if len(command) == 2 {
				if command[1] == "boomer" {
					utils.ReplyToChannel(msg.ChannelID, boomerJoke)
				} else {
					fetchAndRespond(command[1])
				}
				break
			}
			fetchAndRespond("Any")
		case "help":
			session.ChannelMessageSend(msg.ChannelID, "I'll look for therapy places for you in my free time")
		case "comic":
			commands.SendComic(session, msg)
		case "meme":
			var subreddit string
			if len(command) > 1 {
				subreddit = command[1]
			}
			commands.SendMeme(subreddit, session, msg)
		}
	}
}

func createHttpServer() {
	router := httprouter.New()
	router.POST("/api/raino", raingo.ListenToRaindrops)
	router.GET("/api/wake", raingo.WakeMeUpInside)
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
	utils.Session = bot
	createHttpServer()
	bot.Close()
}
