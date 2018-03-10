package main

import (
	"fmt"
	"os" // For JPEG decoding
	// For PNG decoding
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token         = "TOKEN"
	Messages      = []Person{}
	LengthOfArray int
	MsgChannel    = []Channel{}
)

//Person is used for know the name and the NumberOfMessages of a person in the server
type Person struct {
	name        string
	NumberOfMsg int
	isChecked   bool
}

//Channel represents the name and the number of Messages of that Channel
type Channel struct {
	name        string
	NumberOfMsg int
	isChecked   bool
	ID          string
}

func main() {

	// Create a new Discord session using the provided bot token.

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.

	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.

	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	if m.Author.ID == s.State.User.ID {
		return
	}

	//Count how many times a person type a message in the server :)
	if m.Content != "" {
		//var NumberOfMessages int
		//NumberOfMessages++

		Aux4, _ := s.Channel(m.ChannelID)

		if len(Messages) == 0 {
			Aux := Person{name: m.Author.Username, NumberOfMsg: 1, isChecked: true}
			Messages = append(Messages, Aux)
		} else {
			//To create or add messages at the list
			for index := range Messages {
				//To determinate if is the same person to add a msg point to it
				if Messages[index].name == m.Author.Username && Messages[index].isChecked == true {
					Messages[index].NumberOfMsg++
					break
					//To determinate if there aren't that person checked and create the instance
				} else if index+1 == len(Messages) {
					Aux2 := Person{name: m.Author.Username, NumberOfMsg: 1, isChecked: true}
					Messages = append(Messages, Aux2)
					break
				}

			}
		}

		if len(MsgChannel) == 0 {
			Aux := Channel{name: Aux4.Name, NumberOfMsg: 1, isChecked: true, ID: m.ChannelID}
			MsgChannel = append(MsgChannel, Aux)
		} else {
			//To create or add messages at the list
			for index := range MsgChannel {
				//To determinate if is the same Channel
				if MsgChannel[index].ID == m.ChannelID && MsgChannel[index].isChecked == true {
					MsgChannel[index].NumberOfMsg++
					break
					//To determinate if there aren't that person checked and create the instance
				} else if index+1 == len(MsgChannel) {
					Aux2 := Channel{name: Aux4.Name, NumberOfMsg: 1, isChecked: true, ID: m.ChannelID}
					MsgChannel = append(MsgChannel, Aux2)
					break
				}

			}
		}

	}

	if m.Content == "Mensajes" {
		for index := range Messages {
			Aux3 := strconv.Itoa(Messages[index].NumberOfMsg)
			s.ChannelMessageSend(m.ChannelID, Messages[index].name+" "+Aux3)
		}
		for index := range MsgChannel {
			Aux3 := strconv.Itoa(MsgChannel[index].NumberOfMsg)
			s.ChannelMessageSend(m.ChannelID, MsgChannel[index].name+Aux3)
		}
	}
	if m.Content == "Personas" {
		Aux, _ := s.Channel(m.ChannelID)

		Aux2, _ := s.GuildMembers(Aux.GuildID, "", 100)

		for index := range Aux2 {
			Aux3 := Aux2[index].JoinedAt

			s.ChannelMessageSend(m.ChannelID, Aux2[index].User.Username+"\n"+Aux2[index].User.Email+"\n"+Aux3)
		}
	}
}
