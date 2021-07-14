package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func handleMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}
	fmt.Println("Got a message, ", msg.Content)
}

func createHttpServer() {
	router := httprouter.New()
	router.GET("/api/raino", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("endpoint request")
		w.Write([]byte("hello"))
	})
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
