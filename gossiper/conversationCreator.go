package encConversation

import (
	"crypto/rand"
	"github.com/coyim/otr3"
)

func CreateConversation() (c *otr3.Conversation){
	var priv *otr3.DSAPrivateKey
	priv.Generate(rand.Reader)
	c.SetOurKeys([]otr3.PrivateKey{priv})

	// set the Policies.
	c.Policies.RequireEncryption()
	c.Policies.AllowV2()
	c.Policies.AllowV3()
	c.Policies.SendWhitespaceTag()
	c.Policies.WhitespaceStartAKE()
	return c
}