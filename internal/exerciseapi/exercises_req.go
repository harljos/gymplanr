package exerciseapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (c *Client) GetExercise(muscle, difficulty, exerciseType string) (Exercise, error) {
	url := fmt.Sprintf("https://api.api-ninjas.com/v1/exercises?muscle=%s&difficulty=%s&type=%s", muscle, difficulty, exerciseType)

	godotenv.Load(".env")

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return Exercise{}, errors.New("no api key")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Exercise{}, err
	}
	req.Header.Add("X-Api-Key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Exercise{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return Exercise{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Exercise{}, err
	}

	exercise := Exercises{}
	err = json.Unmarshal(data, &exercise)
	if err != nil {
		return Exercise{}, err
	}

	exerciseChosen := rand.Intn(9 - 0) + 0

	return exercise[exerciseChosen], nil
}
