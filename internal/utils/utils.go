package utils

import (
	"github.com/bwmarrin/discordgo"
)

var (
	SpartathlonID int = 865167211944345600
	Session       *discordgo.Session
)

func ReplyToChannel(channelID string, msg string) {
	Session.ChannelMessageSend(channelID, msg)
}
