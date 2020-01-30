package message

//Message corresponds to the message send from the UI client to gossiper

const AUTHENTICATED = "AUTHENTICATION_OK"
const NOT_AUTHENTICATED = "AUTHENTICATION_NOK"

type Message struct {
	Text          string
	Origin        string //for UI
	Room          string //for UI
	Authenticated string //for UI
	Encrypted     bool
	Destination   string
	File          string
	Request       []byte
	AuthQuestion  string
	AuthAnswer    string
}

//Print client message
func (msg *Message) String() string {
	str := "CLIENT MESSAGE "
	str += msg.Text
	return str
}
