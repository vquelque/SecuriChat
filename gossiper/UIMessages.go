package gossiper

import "github.com/vquelque/SecuriChat/message"

func (gsp *Gossiper) sendRumorToUi(text string, rumor *message.RumorMessage) {
	cliMsg := &message.Message{
		Text:   text,
		Origin: rumor.Origin,
		Room:   rumor.Origin,
	}
	gsp.UIMessages <- cliMsg
}

func (gsp *Gossiper) sendAuthQuestionToUi(msg *message.Message) {
	// a peer wants to authenticate. Send Auth question to UI
	cliMsg := &message.Message{
		Origin:        msg.Origin,
		Room:          msg.Origin,
		Authenticated: message.NOT_AUTHENTICATED,
		AuthQuestion:  msg.AuthQuestion,
	}
	gsp.UIMessages <- cliMsg
}
