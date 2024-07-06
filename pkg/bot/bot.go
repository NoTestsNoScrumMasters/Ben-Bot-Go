package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

func main() {
	BotToken := os.Getenv("TOKEN")
	GuildID := os.Getenv("GUILD_ID")

	Run(BotToken, GuildID)
}

func Run(BotToken string, GuildID string) {
	var (
		botCommands = []*discordgo.ApplicationCommand{
			{
				Name:        "help",
				Description: "haha, you're lost.",
			},
			{
				Name:        "channels",
				Description: "List text channels and their threads",
			},
			{
				Name:        "forums",
				Description: "Channels, but with attention spans and interests.",
			},
			{
				Name:        "voice",
				Description: "Voice channels for screeching",
			},
		}

		commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				embed := &discordgo.MessageEmbed{
					Title:       "Help",
					Description: "You probably need it tbh:",
					Color:       0xFFD700, // Gold
					Fields: []*discordgo.MessageEmbedField{
						{Name: "/channels", Value: "List current server text channels and their threads.", Inline: false},
						{Name: "/forums", Value: "List forums and their posts.", Inline: false},
						{Name: "/voice", Value: "List current server voice channels.", Inline: false},
					},
				}
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					log.Printf("Error responding to help command: %v", err)
				}
			},
			"channels": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				embed := &discordgo.MessageEmbed{
					Title:       "Channels",
					Description: "Listing text channels and their threads:",
					Color:       0x00FF00, // Green
				}
				guild, err := s.State.Guild(i.GuildID)
				if err != nil {
					log.Printf("Error fetching guild: %v", err)
					return
				}
				for _, channel := range guild.Channels {
					if channel.Type == discordgo.ChannelTypeGuildText {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   channel.Name,
							Value:  channel.Mention(),
							Inline: false,
						})
						activeThreads, err := s.ChannelThreadsActive(channel.ID)
						if err == nil {
							for _, thread := range activeThreads.Threads {
								embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
									Name:   "Thread: " + thread.Name,
									Value:  thread.Mention(),
									Inline: true,
								})
							}
						} else {
							log.Printf("Error fetching active threads for channel %v: %v", channel.Name, err)
						}
					}
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					log.Printf("Error responding to channels command: %v", err)
				}
			},
			"forums": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				embed := &discordgo.MessageEmbed{
					Title:       "Forums",
					Description: "Message board style channels:",
					Color:       0x00FF00, // Green
				}
				guild, err := s.State.Guild(i.GuildID)
				if err != nil {
					log.Printf("Error fetching guild: %v", err)
					return
				}
				for _, channel := range guild.Channels {
					if channel.Type == discordgo.ChannelTypeGuildForum {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   channel.Name,
							Value:  channel.Mention(),
							Inline: false,
						})
					}
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					log.Printf("Error responding to forums command: %v", err)
				}
			},
			"voice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				embed := &discordgo.MessageEmbed{
					Title:       "Voice Channels",
					Description: "Voice channels where collectivist chants are sung in perfect harmony:",
					Color:       0x00FF00, // Green
				}
				guild, err := s.State.Guild(i.GuildID)
				if err != nil {
					log.Printf("Error fetching guild: %v", err)
					return
				}
				for _, channel := range guild.Channels {
					if channel.Type == discordgo.ChannelTypeGuildVoice {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   channel.Name,
							Value:  channel.Mention(),
							Inline: false,
						})
					}
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					log.Printf("Error responding to voice command: %v", err)
				}
			},
		}
	)

	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	log.Println("Adding commands...")
	log.Printf("User %s", s.State.User.Username)
	for _, v := range botCommands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
