package esclient

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
}

type Client struct {
	client *elasticsearch.Client
}

func NewClient(cfg *Config) (*Client, error) {
	config := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetClient() *elasticsearch.Client {
	return c.client
}
