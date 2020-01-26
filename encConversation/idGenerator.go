package encConversation

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"

	"github.com/vquelque/SecuriChat/constant"
)

const base62Char = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GetRandomName() string {
	var id strings.Builder
	for c := 0; c < constant.NameLength; c++ {
		gen, err := rand.Int(rand.Reader, big.NewInt(62))
		if err != nil {
			log.Fatal("Error generating unique name")
		}
		id.WriteByte(base62Char[gen.Int64()])
	}
	return id.String()
}
