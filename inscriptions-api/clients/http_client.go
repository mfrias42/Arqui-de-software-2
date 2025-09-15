package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPClient struct {
	usersAPIURL   string
	coursesAPIURL string
	client        *http.Client
}

func NewHTTPClient(usersAPIURL, coursesAPIURL string) *HTTPClient {
	return &HTTPClient{
		usersAPIURL:   usersAPIURL,
		coursesAPIURL: coursesAPIURL,
		client:        &http.Client{},
	}
}

func (c *HTTPClient) CheckUserExists(userID uint) error {
	resp, err := c.client.Get(fmt.Sprintf("%s/users/%d", c.usersAPIURL, userID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("user does not exist")
	}

	return nil
}

func (c *HTTPClient) CheckCourseExists(courseID uint) error {
	resp, err := c.client.Get(fmt.Sprintf("%s/courses/%d", c.coursesAPIURL, courseID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("course with ID %d not found", courseID)
	}

	return nil
}

type CourseDetails struct {
	ID        uint `json:"id"`
	Capacity  int  `json:"capacity"`
	Available bool `json:"available"`
	// Add other fields if needed
}

func (c *HTTPClient) GetCourseDetails(courseID uint) (*CourseDetails, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/courses/%d", c.coursesAPIURL, courseID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("course with ID %d not found", courseID)
	}

	var course CourseDetails
	if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
		return nil, fmt.Errorf("error decoding course details: %v", err)
	}

	return &course, nil
}

func (c *HTTPClient) UpdateCourseAvailability(courseID int64) error {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/courses/%d/availability", c.coursesAPIURL, courseID), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update course availability for course ID %d", courseID)
	}

	return nil
}
