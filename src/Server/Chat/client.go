package AxmorChat

type Client struct {
	Name string
	Connected bool
}

//Конструктор Client
func NewClient (name string, connected bool) *Client{
	return &Client {name, connected, }
}

//Функция для возврата состояния состояния
func StatusClient (client *Client) string{
	if client.Connected == true{
		return client.Name + " онлайн \n"
	}
	return client.Name + " оффлайн \n"
}
