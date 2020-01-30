package gossiper

import (
	"github.com/vquelque/SecuriChat/encConversation"
	"github.com/vquelque/SecuriChat/message"
)

func (gsp *Gossiper) sendRumorToUi(text string, rumor *message.RumorMessage) {
	auth := message.NOT_AUTHENTICATED
	if rumor.EncryptedMessage.Step == encConversation.AuthenticationOK {
		auth = message.AUTHENTICATED
	}
	cliMsg := &message.Message{
		Text:          text,
		Origin:        rumor.Origin,
		Room:          rumor.Origin,
		Authenticated: auth,
	}
	select {
	case gsp.UIMessages <- cliMsg:
	default:
	}
}

func (gsp *Gossiper) sendAuthQuestionToUi(msg *message.Message) {
	// a peer wants to authenticate. Send Auth question to UI
	cliMsg := &message.Message{
		Origin:        msg.Origin,
		Room:          msg.Origin,
		Authenticated: message.NOT_AUTHENTICATED,
		AuthQuestion:  msg.AuthQuestion,
	}
	select {
	case gsp.UIMessages <- cliMsg:
	default:
	}
}

func (gsp *Gossiper) sendAuthOK(peerID string) {
	cliMsg := &message.Message{
		Room:          peerID,
		Authenticated: message.AUTHENTICATED,
	}
	select {
	case gsp.UIMessages <- cliMsg:
	default:
	}
}
