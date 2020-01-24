package gossiper

import (
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/encConversation"
)

func (gsp *Gossiper) CreateConversationState() (cs *encConversation.ConversationState){
	c := &otr3.Conversation{}
	priv := gsp.loadPrivateKey()
	c.SetOurKeys([]otr3.PrivateKey{priv})

	// set the Policies.
	c.Policies.RequireEncryption()
	c.Policies.AllowV3()
	cs =  &encConversation.ConversationState{
		IsFinished:   false,
		Step:         0,
		Conversation: c,
	}


	return cs
}

func (gsp *Gossiper) loadPrivateKey() *otr3.DSAPrivateKey {
	return gsp.privateKey
}
