package encConversation

import (
	"github.com/coyim/otr3"
	"sync"
)

const(
	QueryMsg = iota
	DHCommit = iota
	DHKey = iota
	RevealSig = iota
	Sig = iota
	AKE_Finished = iota
)


type EncryptedMessage struct {
	Message otr3.ValidMessage
	Step int
	Dest string
}


type ConversationState struct {
	IsFinished bool
	Step int // step of the Auth. Key. Exchange (AKE)
	Conversation *otr3.Conversation
}

type ConvStateMap struct {
	Map  map[string]*ConversationState // key: name
	Lock sync.RWMutex
}

func InitConvStateMap() (convMap *ConvStateMap) {
	convMap = &ConvStateMap{}
	convMap.Map = make(map[string]*ConversationState)
	convMap.Lock = sync.RWMutex{}
	return convMap
}

func (convMap *ConvStateMap) Update(k string, v *ConversationState) {
	convMap.Lock.Lock()
	convMap.Map[k] = v
	convMap.Lock.Unlock()
}

func (convMap *ConvStateMap) Load(k string) (v *ConversationState, ok bool) {
	convMap.Lock.RLock()
	v, ok = convMap.Map[k]
	ok = ok && v != nil
	convMap.Lock.RUnlock()
	return v, ok
}
