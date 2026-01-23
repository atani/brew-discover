package api

import (
	"fmt"
)

type Cask struct {
	Token       string   `json:"token"`
	FullToken   string   `json:"full_token"`
	Tap         string   `json:"tap"`
	Name        []string `json:"name"`
	Desc        string   `json:"desc"`
	Homepage    string   `json:"homepage"`
	Version     string   `json:"version"`
	Deprecated  bool     `json:"deprecated"`
	Disabled    bool     `json:"disabled"`
}

func (c *Cask) GetName() string {
	if len(c.Name) > 0 {
		return c.Name[0]
	}
	return c.Token
}

func (c *Client) GetCasks() ([]Cask, error) {
	url := fmt.Sprintf("%s/cask.json", BaseURL)
	var casks []Cask
	if err := c.get(url, &casks); err != nil {
		return nil, err
	}
	return casks, nil
}

func (c *Client) GetCasksBytes() ([]byte, error) {
	url := fmt.Sprintf("%s/cask.json", BaseURL)
	return c.getBytes(url)
}

func (c *Client) GetCask(name string) (*Cask, error) {
	url := fmt.Sprintf("%s/cask/%s.json", BaseURL, name)
	var cask Cask
	if err := c.get(url, &cask); err != nil {
		return nil, err
	}
	return &cask, nil
}
