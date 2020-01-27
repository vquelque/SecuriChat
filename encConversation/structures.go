package encConversation

import (
	"sync"

	"github.com/coyim/otr3"
)

const (
	QueryMsg         = iota
	DHCommit         = iota
	DHKey            = iota
	RevealSig        = iota
	Sig              = iota
	AkeFinished      = iota
	SMP1             = iota
	SMP2             = iota
	SMP3             = iota
	SMP4             = iota
	SMP5             = iota
	AuthenticationOK = iota
	QueryTextMessage = "?OTRv3?"
)

type ConversationState struct {
	Step         int // step of the Auth. Key. Exchange (AKE)
	Conversation *otr3.Conversation
	Buffer       chan string
	AnswerChan   chan string
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

func (convMap *ConvStateMap) DestroyConversation(dest string) {
	convMap.Update(dest,nil)
}


