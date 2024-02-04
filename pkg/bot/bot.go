package bot

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// Utility function to create an embed in response to an interaction
func create_embed(name string, session *discordgo.Session, interaction *discordgo.InteractionCreate, description string, Fields []*discordgo.MessageEmbedField) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xFFE41E,
		Description: description,

		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     name,
		Fields:    Fields,
	}

	// Send the embed as a response to the provided interaction
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}})
}

// Utility function to create an embed in response to an interaction
func send_embed(name string, session *discordgo.Session, user string, description string, Fields []*discordgo.MessageEmbedField) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xFFE41E,
		Description: description,

		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     name,
		Fields:    Fields,
	}

	// Send the embed as a response to the provided interaction
	channel, err := session.UserChannelCreate(user)

	if err != nil {
		fmt.Println(err)
	} else {
		session.ChannelMessageSendEmbed(channel.ID, embed)
	}
}

const cmdUsage = "USAGE: /wiki [function/callback]"

type Results struct {
	Status struct {
		Total      int `json:"total"`
		Failed     int `json:"failed"`
		Successful int `json:"successful"`
	} `json:"status"`
	Request struct {
		Query struct {
			Query string `json:"query"`
		} `json:"query"`
		Size      int `json:"size"`
		From      int `json:"from"`
		Highlight struct {
			Style  interface{} `json:"style"`
			Fields interface{} `json:"fields"`
		} `json:"highlight"`
		Fields           interface{} `json:"fields"`
		Facets           interface{} `json:"facets"`
		Explain          bool        `json:"explain"`
		Sort             []string    `json:"sort"`
		IncludeLocations bool        `json:"includeLocations"`
		SearchAfter      interface{} `json:"search_after"`
		SearchBefore     interface{} `json:"search_before"`
	} `json:"request"`
	Hits      []Hit `json:"hits"`
	TotalHits int   `json:"total"`
	Took      int64 `json:"took"`
}

type Hit struct {
	Url                  string  `json:"url"`
	Title                string  `json:"title"`
	Description          string  `json:"desc"`
	TitleFragments       string  `json:"title_fragment"`
	DescriptionFragments string  `json:"desc_fragment"`
	Score                float64 `json:"score"`
}

func Run(BotToken string) {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Bot broken")
		os.Exit(1)
	}

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		// add more cases if required
	}

}
