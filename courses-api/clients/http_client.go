package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient struct {
	inscriptionsAPIURL string
	client             *http.Client
}

func NewHTTPClient(inscriptionsAPIURL string) *HTTPClient {
	return &HTTPClient{
		inscriptionsAPIURL: inscriptionsAPIURL,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type Inscription struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}

func (c *HTTPClient) GetInscriptionsByCourse(courseID uint) ([]Inscription, error) {
	url := fmt.Sprintf("%s/courses/%d/inscriptions", c.inscriptionsAPIURL, courseID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to inscriptions API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get inscriptions for course %d: status code %d", courseID, resp.StatusCode)
	}

	var inscriptions []Inscription
	if err := json.NewDecoder(resp.Body).Decode(&inscriptions); err != nil {
		return nil, fmt.Errorf("error decoding inscriptions: %v", err)
	}

	return inscriptions, nil
}
