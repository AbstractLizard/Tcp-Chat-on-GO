package AxmorChat

type Message struct {
	Client *Client
	Msg    string
}

func ReturnMessage(msg *Message) string{
	if msg == nil {
		return "сообщение не найдено"
	}
	return msg.Client.Name + " -> " + msg.Msg
}
