package cogs

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func CleanHistory(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandCleanHistory) {
			args := strings.Fields(m.Content)
			if len(args) < 2 {
				d.ChannelMessageSend(m.ChannelID, "Wrong number of arguments use >ban userid or >ban @UserPing")
				return
			}
			if !config.IsAdmin(d, m.Author.ID, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "You are not allowed to use this command.")
				return
			}
			if !storage.IsChannelAllowed(db, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "Messages in this channel are not allowed.")
				return
			}

			targetID, ok := ParseUserId(args[1])
			if !ok {
				d.ChannelMessageSend(m.ChannelID, "Wrong user id.")
				return
			}

			if targetID == d.State.User.ID {
				d.ChannelMessageSend(m.ChannelID, "You can't clean the bot history.")
				return
			}

			if _, err := d.User(targetID); err != nil {
				d.ChannelMessageSend(m.ChannelID, "This user probably doesn't exists lol")
				return
			}

			if !storage.IsUserExist(db, targetID) {
				isAdmin := config.IsAdmin(d, targetID, m.ChannelID)
				storage.AddUser(db, targetID, false, isAdmin)
			}

			err := storage.CleanUserQuestionsHistory(db, targetID)
			if err != nil {
				log.Println(err)
				return
			}

			d.ChannelMessageSend(m.ChannelID, "User questions reputation has been cleaned")
		}
	}

}
