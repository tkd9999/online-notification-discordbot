package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var statushistory map[string]bool
var mobilehistory map[string]time.Time

func init() {
	statushistory = make(map[string]bool)
	mobilehistory = make(map[string]time.Time)
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

	// members, err := dg.GuildMembers(os.Getenv("GUILD_ID"), "", 1000)
	// if err != nil {
	// 	log.Println("Error getting members:", err)
	// }

	// guild, err := dg.Guild(os.Getenv("GUILD_ID"))
	// if err != nil {
	// 	log.Println("Error getting guild:", err)
	// }

	// onlinemembers := make([]string, 0)
	// for _, member := range members {
	// 	if status, _ := guild.GetPresence(member.User.ID).ClientStatus.Desktop; status == discordgo.StatusOnline {
	// 		onlinemembers = append(onlinemembers, member.User.Username)
	// 	}
	// }

	// if len(onlinemembers) > 0 {
	// 	onlinemembersstring := strings.Join(onlinemembers, " ")
	// 	_, err = dg.ChannelMessageSend(os.Getenv("CHANNEL_ID"), fmt.Sprintf("がオンラインです: %s", onlinemembersstring))
	// 	if err != nil {
	// 		log.Println("Error sending startup message:", err)
	// 	}
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
			if time.Now().Sub(mobilehistory[m.User.ID]).Minutes() < 5 {
				return
			}
			mobilehistory[m.User.ID] = time.Now()
			_, err := s.ChannelMessageSend(os.Getenv("CHANNEL_ID"), fmt.Sprintf("%sがオンラインです(モバイル)", user.Username))
			if err != nil {
				log.Println("Error sending message:", err)
			}
		}
	} else if m.Presence.Status == discordgo.StatusOffline {
		statushistory[m.User.ID] = false
	}
}
