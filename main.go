package main

import (
	"TheNeonCenturyBot2/cogs"
	"TheNeonCenturyBot2/storage"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func GeminiClient(geminiAPIKey string) (*genai.Client, error) {
	return genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  geminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
}
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db := storage.DbStart()

	token := os.Getenv("BOT_TOKEN")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	ai, err := GeminiClient(geminiAPIKey)
	if err != nil {
		log.Fatal("Gemini error: ", err)
	}

	if token == "" {
		log.Fatal("BOT_TOKEN environment variable not set")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	discord.AddHandler(allowedGuilds)
	discord.AddHandler(cogs.Help)
	discord.AddHandler(cogs.AllowChannel(db))
	discord.AddHandler(cogs.DisAllowChannel(db))
	discord.AddHandler(cogs.BanUser(db))
	discord.AddHandler(cogs.UnBanUser(db))
	discord.AddHandler(cogs.Ask(db, ai))
	discord.AddHandler(cogs.CleanHistory(db))
	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			log.Fatal("Error closing Discord session: ", err)
		}
	}(discord)

	err = discord.UpdateCustomStatus(">help")

	fmt.Println("The Neon Century bot is running. Ready to summarize lore")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func allowedGuilds(d *discordgo.Session, g *discordgo.GuildCreate) {
	publicServer := os.Getenv("PUBLIC_SERVER")
	devServer := os.Getenv("DEV_SERVER")
	if g.ID != publicServer && g.ID != devServer {
		log.Println("Unauthorized discord server found, Leaving")

		err := d.GuildLeave(g.ID)
		if err != nil {
			log.Fatal("Error leaving discord server: ", err)
		}
	}

	if g.ID == devServer {
		log.Println("Working on dev server")

	}
}
