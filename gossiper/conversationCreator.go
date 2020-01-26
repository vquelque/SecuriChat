package gossiper

import (
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/encConversation"
	"log"
)

const maxBufferSize  = 100



type DebugSMPEventHandler struct{
	dest string
	convStateMap *encConversation.ConvStateMap
}

// HandleSMPEvent dumps all SMP events
func (d DebugSMPEventHandler) HandleSMPEvent(event otr3.SMPEvent, progressPercent int, question string) {
	if event.String() == "SMPEventAbort"{
		log.Println("Warning, error with SMP Protocol, conversation will be destroyed")
		d.convStateMap.DestroyConversation(d.dest)
	}
}



func (gsp *Gossiper) createConversationState(dest string) (cs *encConversation.ConversationState){
	c := &otr3.Conversation{}
	priv := gsp.loadPrivateKey()
	c.SetOurKeys([]otr3.PrivateKey{priv})

	// set the Policies.
	c.Policies.RequireEncryption()
	c.Policies.AllowV3()
	handler := DebugSMPEventHandler{dest:dest,convStateMap:gsp.convStateMap}
	c.SetSMPEventHandler(handler)
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
		cs = gsp.createConversationState(dest)
		gsp.convStateMap.Update(dest, cs)
	}
	return cs,ok
}

func (gsp *Gossiper) loadPrivateKey() *otr3.DSAPrivateKey {
	return gsp.privateKey
}
