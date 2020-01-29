package gossiper

import "github.com/vquelque/SecuriChat/message"

func (gsp *Gossiper) sendRumorToUi(rumor *message.RumorMessage) {
	cliMsg := &message.Message{
		Text:   rumor.Text,
		Origin: rumor.Origin,
		Room:   rumor.Origin,
	}
	gsp.UIMessages <- cliMsg
}
