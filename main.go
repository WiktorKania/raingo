package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

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
		command := message[3:]
		fmt.Println("Got command, ", command)
		switch command {
		case "joke":
			tellJoke(session, msg)
		case "help":
			session.ChannelMessageSend(msg.ChannelID, "I'll look for therapy places for you in my free time")
		}
	}
}

func createHttpServer() {
	fmt.Println("hey", httprouter.New())
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
