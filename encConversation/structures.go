package encConversation

import (
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/utils"
	"sync"
)

const (
	QueryMsg         = iota
	DHCommit         = iota
	DHKey            = iota
	RevealSig        = iota
	Sig              = iota
	AKE_Finished     = iota
	QueryTextMessage = "?OTRv3?"
)

type EncryptedMessage struct {
	Message otr3.ValidMessage
	Step    int
	Dest    string
}

func (enc *EncryptedMessage) Encode() []byte {
	b := utils.EncodeUint64(uint64(enc.Step))
	b = append(b, []byte(enc.Dest)...)
	b = append(b, enc.Message...)
	return b
}

type ConversationState struct {
	Step         int // step of the Auth. Key. Exchange (AKE)
	Conversation *otr3.Conversation
	Buffer       chan string
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
