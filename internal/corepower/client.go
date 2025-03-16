package corepower

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	reservationEndpoint = "https://api2.corepoweryoga.com"
)

type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.base.RoundTrip(req)
}

type Reservation struct {
	CenterId  string  `json:"centerId"`
	SessionId float32 `json:"sessionId"`
}

type Client struct {
	Endpoint          *url.URL
	ReservationClient *http.Client
	Token             string
}

type ReservationResponse struct {
	ID                 int    `json:"id"`
	StartTime          string `json:"startTime"`
	EndTime            string `json:"endTime"`
	StartTimeUTC       string `json:"startTimeUTC"`
	EndTimeUTC         string `json:"endTimeUTC"`
	SessionName        string `json:"sessionName"`
	RegistrationStatus int    `json:"registrationStatus"`
	InstructorID       string `json:"instructorId"`
	Instructor         string `json:"instructor"`
	ImagePaths         struct {
		Px64  string `json:"px64"`
		Px100 string `json:"px100"`
		Px200 string `json:"px200"`
		Px400 string `json:"px400"`
		Px800 string `json:"px800"`
	} `json:"imagePaths"`
	Center                   string `json:"center"`
	ClassName                string `json:"className"`
	CanCancel                bool   `json:"canCancel"`
	CurrentWaitlistPosition  any    `json:"currentWaitlistPosition"`
	ClassType                int    `json:"classType"`
	ClassID                  int    `json:"classId"`
	SessionID                int    `json:"sessionId"`
	CenterID                 string `json:"centerId"`
	StudentVirtualLink       string `json:"studentVirtualLink"`
	IsVirtualClass           bool   `json:"isVirtualClass"`
	VirtualType              int    `json:"virtualType"`
	VirtualGuestReservations any    `json:"virtualGuestReservations"`
	InvoiceID                string `json:"invoiceId"`
	Price                    int    `json:"price"`
	IsInstructorSubstituted  bool   `json:"is_instructor_substituted"`
	Status                   string `json:"status"`
}

func (c *Client) Initialize() {
	c.Endpoint, _ = url.Parse(reservationEndpoint)
	c.ReservationClient = &http.Client{}

	c.ReservationClient.Transport = &headerTransport{
		base: http.DefaultTransport,
		headers: map[string]string{
			"authorization": c.Token,
		},
	}
}

func (c *Client) Reserve(centerId string, sessionId float32) (ReservationResponse, error) {

	reservation := Reservation{
		CenterId:  centerId,
		SessionId: sessionId,
	}
	payload, err := json.Marshal(reservation)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("error marshaling reservation: %v", err)
	}

	baseURL := c.Endpoint
	baseURL.Path = "/reservation"

	req, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-api-version", "2.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.ReservationClient.Do(req)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ReservationResponse{}, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var reservationResponse ReservationResponse
	err = json.Unmarshal(body, &reservationResponse)
	if err != nil {
		return ReservationResponse{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return reservationResponse, nil
}
