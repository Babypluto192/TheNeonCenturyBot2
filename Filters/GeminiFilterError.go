package Filters

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func GeminiFilterError(channelId string, d *discordgo.Session, filteredQuestion string) bool {
	if filteredQuestion == "" {
		log.Println("Something went wrong with filter")
		d.ChannelMessageSend(channelId, "Something went wrong with filter, try again later")
		return false
	}
	return true
}
