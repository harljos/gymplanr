package exerciseapi

import (
	"net/http"
	"time"

	"github.com/harljos/gymplanr/internal/exercisecache"
)

type Client struct {
	httpClient http.Client
	cache      exercisecache.Cache
}

type Exercises []struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
}

type Exercise struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
}

func NewClient(cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		cache: exercisecache.NewCache(cacheInterval),
	}
}
