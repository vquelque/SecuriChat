package test

import (
	"github.com/vquelque/SecuriChat/encConversation"
	"github.com/vquelque/SecuriChat/pow"
	"testing"
)

func TestPow(t *testing.T) {
	h := []byte("Testytgvbhycfg vhjkuytvcgh bjkuvgh bjhgvh bnhgb ")
	newPoW := pow.NewProofOfWork(h)
	if newPoW.Validator(h) != true{
		t.Errorf("Incorrect Hash")

	}
}

func TestPow2(t *testing.T) {
	h := encConversation.EncryptedMessage{
		Message: []byte("hellohjvvgezqegvdsvedkhgzvhefqesfcvghbjigjcfhbjkgf"),
		Step:    2,
		Dest:    "ee",
	}

	newPoW := pow.NewProofOfWork(h.Encode())
	if newPoW.Validator(h.Encode()) != true{
		t.Errorf("Incorrect Hash")

	}
}
