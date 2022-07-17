package client

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c Client) Connect() error {
	return nil
}

func (c Client) Subscribe(channel chan string) {
}
