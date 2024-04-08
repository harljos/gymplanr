package exerciseapi

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

type Exercise []struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
	}
}