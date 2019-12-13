package observer

import (
	"fmt"
	"sync"

	"github.com/vquelque/Peerster/vector"
)

// Observer structure used for callback to routine
type Observer struct {
	waitingForAck map[string]chan bool
	lock          sync.RWMutex
}

func Init() *Observer {
	obs := &Observer{waitingForAck: make(map[string]chan bool)}
	return obs
}

func (obs *Observer) Register(sender string) chan bool {
	obs.lock.Lock()
	defer obs.lock.Unlock()
	ackChan := make(chan bool)
	obs.waitingForAck[sender] = ackChan
	return ackChan
}

func (obs *Observer) Unregister(sender string) {
	obs.lock.Lock()
	defer obs.lock.Unlock()
	delete(obs.waitingForAck, sender)
}

func (obs *Observer) GetObserver(sp *vector.StatusPacket, peer string) chan bool {
	obs.lock.RLock()
	defer obs.lock.RUnlock()
	for _, ps := range sp.Want {
		id := peer + fmt.Sprintf("%s : %d", ps.Identifier, ps.NextID-1)
		ackChan, found := obs.waitingForAck[id]
		if found {
			return ackChan
		}
	}
	return nil
}
