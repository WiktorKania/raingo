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

}

func makeEspresso(ch chan<- string) {

}

func makeLatte(ch chan<- string) {

}

func MakeCoffee(session *discordgo.Session, msg *discordgo.MessageCreate) {

	// session.ChannelMessageSend(msg.ChannelID, "Here's your coffee: "+coffeeResult)
}
