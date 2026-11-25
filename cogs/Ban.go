package cogs

import (
	"TheNeonCenturyBot2/Filters"
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func BanUser(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandBanUser) {
			args := strings.Fields(m.Content)
			if len(args) < 2 {
				d.ChannelMessageSend(m.ChannelID, "Wrong number of arguments use >ban userid or >ban @UserPing")
				return
			}
			if !Filters.IsAdminFilter(m.ChannelID, d, m.Author.ID) {
				return
			}

			if !Filters.IsChannelAllowedFilter(db, m.ChannelID, d) {
				return
			}

			targetID, ok := ParseUserId(args[1])
			if !ok {
				d.ChannelMessageSend(m.ChannelID, "Wrong user id.")
				return
			}
			if targetID == m.Author.ID {
				d.ChannelMessageSend(m.ChannelID, "You can't ban yourself [baka](https://klipy.com/gifs/nanami-aoyama-nanami)")
				return
			}
			if targetID == d.State.User.ID {
				d.ChannelMessageSend(m.ChannelID, "You can't ban the bot.")
				return
			}
			if config.IsBanned(db, targetID) {
				d.ChannelMessageSend(m.ChannelID, "User already banned")
				return
			}

			if _, err := d.User(targetID); err != nil {
				d.ChannelMessageSend(m.ChannelID, "This user probably doesn't exists lol")
				return
			}
			if !storage.IsUserExist(db, targetID) {
				isAdmin := config.IsAdmin(d, targetID, m.ChannelID)
				storage.AddUser(db, targetID, false, isAdmin)
				if isAdmin {
					d.ChannelMessageSend(m.ChannelID, "You are not allowed to ban admins.")
					return
				}
			}

			err := storage.BanUser(db, targetID)
			if err != nil {
				log.Println(err)
				return
			}

			d.ChannelMessageSend(m.ChannelID, "User has been banned.")
		}
	}

}
