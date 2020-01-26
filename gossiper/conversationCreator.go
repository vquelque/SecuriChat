package gossiper

import (
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/encConversation"
)

const maxBufferSize  = 100

func (gsp *Gossiper) createConversationState() (cs *encConversation.ConversationState){
	c := &otr3.Conversation{}
	priv := gsp.loadPrivateKey()
	c.SetOurKeys([]otr3.PrivateKey{priv})

	// set the Policies.
	c.Policies.RequireEncryption()
	c.Policies.AllowV3()
	cs =  &encConversation.ConversationState{
		Step:         0,
		Conversation: c,
		Buffer: make(chan string,maxBufferSize),
	}


	return cs
}

func (gsp *Gossiper) createOrLoadConversationState(dest string) (*encConversation.ConversationState,bool) {
	cs, ok := gsp.convStateMap.Load(dest)
	if !ok {
		cs = gsp.createConversationState()
		gsp.convStateMap.Update(dest, cs)
	}
	return cs,ok
}

func (gsp *Gossiper) loadPrivateKey() *otr3.DSAPrivateKey {
	return gsp.privateKey
}
