package discord

import (
	"bangarang/musec/model"
	"fmt"
)

// TextMessage send a message to the current text channel
func TextMessage(sd *model.ServerDetails, content string) {
	s, tc := sd.Session, sd.TextChannel

	_, err := s.ChannelMessageSend(tc.ID, content)

	if err != nil {
		println(fmt.Sprintf("Error while sending message to %s; %s\n", tc.Name, content))
	}
}