package commands

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func wait() {
	time.Sleep(time.Second * time.Duration(rand.Intn(3)))
}

func makeAmericano(ch chan<- string) {
	wait()
	ch <- "Americano"
}

func makeEspresso(ch chan<- string) {
	wait()
	ch <- "Espresso"
}

func makeLatte(ch chan<- string) {
	wait()
	ch <- "Latte"
}

func MakeCoffee(session *discordgo.Session, msg *discordgo.MessageCreate) {
	coffeeChan := make(chan string, 3)
	go makeAmericano(coffeeChan)
	go makeEspresso(coffeeChan)
	go makeLatte(coffeeChan)
	coffeeResult := <-coffeeChan
	session.ChannelMessageSend(msg.ChannelID, "Here's your coffee: "+coffeeResult)
}
