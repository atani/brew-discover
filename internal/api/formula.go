package api

import (
	"fmt"
)

type Formula struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Tap         string   `json:"tap"`
	Desc        string   `json:"desc"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Versions    Versions `json:"versions"`
	Deprecated  bool     `json:"deprecated"`
	Disabled    bool     `json:"disabled"`
	Dependencies []string `json:"dependencies"`
}

type Versions struct {
	Stable string `json:"stable"`
	Head   string `json:"head"`
}

func (c *Client) GetFormulae() ([]Formula, error) {
	url := fmt.Sprintf("%s/formula.json", BaseURL)
	var formulae []Formula
	if err := c.get(url, &formulae); err != nil {
		return nil, err
	}
	return formulae, nil
}

func (c *Client) GetFormulaeBytes() ([]byte, error) {
	url := fmt.Sprintf("%s/formula.json", BaseURL)
	return c.getBytes(url)
}

func (c *Client) GetFormula(name string) (*Formula, error) {
	url := fmt.Sprintf("%s/formula/%s.json", BaseURL, name)
	var formula Formula
	if err := c.get(url, &formula); err != nil {
		return nil, err
	}
	return &formula, nil
}
