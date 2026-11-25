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
	"github.com/revrost/go-openrouter"
	"google.golang.org/genai"
	"gorm.io/gorm"
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
	clients := MainAI()

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}
	LoadCogs(discord, db, ai, clients)
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

func LoadCogs(discord *discordgo.Session, db *gorm.DB, ai *genai.Client, clients []*openrouter.Client) {
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	discord.AddHandler(allowedGuilds)
	discord.AddHandler(cogs.Help)
	discord.AddHandler(cogs.AllowChannel(db))
	discord.AddHandler(cogs.DisAllowChannel(db))
	discord.AddHandler(cogs.BanUser(db))
	discord.AddHandler(cogs.UnBanUser(db))
	discord.AddHandler(cogs.Ask(db, ai, clients))
	discord.AddHandler(cogs.CleanHistory(db))

	log.Println("Loaded Cogs")
}

func MainAI() []*openrouter.Client {
	apiKeys := []string{
		"OPEN_ROUTER_API_KEY_1",
		"OPEN_ROUTER_API_KEY_2",
		"OPEN_ROUTER_API_KEY_3",
		"OPEN_ROUTER_API_KEY_4",
	}
	var clients []*openrouter.Client

	for _, apiKey := range apiKeys {
		if key := os.Getenv(apiKey); key != "" {
			c := openrouter.NewClient(key, openrouter.WithXTitle("TheNeonCenturyBot"))
			clients = append(clients, c)
		}
	}

	if len(clients) == 0 {
		log.Println("No Open Router API keys found")
		return nil
	}

	return clients
}
