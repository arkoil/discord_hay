package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var Store = make(map[string]bool)

var (
	Token    = flag.String("t", "", "Bot token")
	GuildID  = flag.String("gi", "", "Guild id")
	ChanelId = flag.String("ci", "", "Chanel id, for message")
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + *Token)
	checkErr(err)
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = discord.Open()
	checkErr(err)
	discord.AddHandler(CheckState)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

//CheckState Handler for checking voice state
func CheckState(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.GuildID != *GuildID {
		return
	}
	g, _ := s.State.Guild(event.GuildID)
	newStore := make(map[string]bool)
	for _, vs := range g.VoiceStates {
		if !Store[vs.UserID] {
			name := ""
			for _, member := range g.Members {
				if member.User.ID == vs.UserID {
					name = member.User.Username
					if member.Nick != "" {
						name = member.Nick
					}
				}

			}

			_, err := s.ChannelMessageSend(*ChanelId, "Привет "+name)
			checkErr(err)
		}
		newStore[vs.UserID] = true
	}
	Store = newStore

}

// Checking error
func checkErr(err error) {
	if err != nil {
		ErrorLog.Println("ERROR", err)
	}
}
