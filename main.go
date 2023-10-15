package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var statushistory map[string]bool

func init() {
	statushistory = make(map[string]bool)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating discord session:", err)
		return
	}

	dg.LogLevel = discordgo.LogDebug

	// _, err = dg.ChannelMessageSend(os.Getenv("CHANNEL_ID"), "Bot is online!")
	// if err != nil {
	// 	log.Println("Error sending startup message:", err)
	// }

	dg.AddHandler(onlinenotification)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildPresences
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onlinenotification(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	fmt.Println("status-update received")
	if m.Presence.Status == discordgo.StatusOnline {

		if statushistory[m.User.ID] == true {
			return
		}

		user, _ := s.User(m.User.ID)
		if m.Presence.ClientStatus.Desktop == discordgo.StatusOnline {
			statushistory[m.User.ID] = true
			_, err := s.ChannelMessageSend(os.Getenv("CHANNEL_ID"), fmt.Sprintf("%sがオンラインです", user.Username))
			if err != nil {
				log.Println("Error sending message:", err)
			}
		} else if m.Presence.ClientStatus.Mobile == discordgo.StatusOnline {
			_, err := s.ChannelMessageSend(os.Getenv("CHANNEL_ID"), fmt.Sprintf("%sが偵察しています", user.Username))
			if err != nil {
				log.Println("Error sending message:", err)
			}
		}
	} else if m.Presence.Status == discordgo.StatusOffline {
		statushistory[m.User.ID] = false
	}
}
