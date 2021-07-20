package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

type UserMsg struct {
	Nickname string
	Msg      string
	Channel  string
	ImageURL string `json:",omitempty"`
}

func listenToRaindrops(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var userMsg UserMsg
	err := decoder.Decode(&userMsg)
	if err != nil {
		log.Fatalln(err)
	}
	fields := []*discordgo.MessageEmbedField{
		{Name: "Author", Value: userMsg.Nickname, Inline: true},
		{Name: "Message", Value: userMsg.Msg, Inline: true},
		{Name: "Channel", Value: userMsg.Channel, Inline: true},
	}
	var imageEmbed discordgo.MessageEmbedImage
	if userMsg.ImageURL != "" {
		imageEmbed = discordgo.MessageEmbedImage{URL: userMsg.ImageURL}
	}
	messageEmbed := discordgo.MessageEmbed{Fields: fields, Image: &imageEmbed}
	session.ChannelMessageSend(strconv.Itoa(SpartathlonID), "Someone from **Raino** is calling:")
	session.ChannelMessageSendEmbed(strconv.Itoa(SpartathlonID), &messageEmbed)
}
