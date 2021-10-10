package discord

import (
	"fmt"
	"log"
	"musec/model"
)

// TextMessage send a message to the current text channel
func TextMessage(sd *model.ServerDetails, format string, args ...interface{}) {
	s, tc := sd.Session, sd.TextChannel

	content := fmt.Sprintf(format, args) + "\n"

	_, err := s.ChannelMessageSend(tc.ID, content)

	if err != nil {
		log.Printf("Error while sending message to [%s] with content [%s]\n", tc.Name, content)
	}
}
