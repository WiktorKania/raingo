package utils

import "github.com/bwmarrin/discordgo"

func tellJoke(session *discordgo.Session, msg *discordgo.MessageCreate, joke string) {
	session.ChannelMessageSend(msg.ChannelID, joke)
}

func replyToChannel(channelID string, msg string) {
	session.ChannelMessageSend(channelID, msg)
}
