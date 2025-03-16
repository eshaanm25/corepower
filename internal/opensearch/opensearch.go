package opensearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	openSearchEndpoint = "https://9a6fd5b868a24c039c29be8d490b3766.ent-search.us-west-1.aws.cloud.es.io"
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

type Client struct {
	SearchClient *http.Client
	Endpoint     *url.URL
}

type SearchRequest struct {
	Query string `json:"query"`
	Page  struct {
		Current int `json:"current"`
		Size    int `json:"size"`
	} `json:"page"`
	Sort struct {
		StartTime string `json:"start_time"`
	} `json:"sort"`
	Filters struct {
		Any []struct {
			StartTimeUtc struct {
				From time.Time `json:"from"`
				To   time.Time `json:"to"`
			} `json:"start_time_utc"`
		} `json:"any"`
		All []struct {
			CenterID []string `json:"center.id"`
		} `json:"all"`
	} `json:"filters"`
}

type SearchResponse struct {
	Meta struct {
		Alerts []interface{} `json:"alerts"`
		Engine struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"engine"`
		Page struct {
			Current      int `json:"current"`
			Size         int `json:"size"`
			TotalPages   int `json:"total_pages"`
			TotalResults int `json:"total_results"`
		} `json:"page"`
		Precision int           `json:"precision"`
		RequestID string        `json:"request_id"`
		Warnings  []interface{} `json:"warnings"`
	} `json:"meta"`
	Results []struct {
		Meta struct {
			Engine string      `json:"engine"`
			ID     string      `json:"id"`
			Score  interface{} `json:"score"`
		} `json:"_meta"`
		AvailableSlots struct {
			Raw float32 `json:"raw"`
		} `json:"available_slots"`
		CanBook struct {
			Raw string `json:"raw"`
		} `json:"can_book"`
		CanBookStatus struct {
			Raw interface{} `json:"raw"`
		} `json:"can_book_status"`
		Capacity struct {
			Raw float32 `json:"raw"`
		} `json:"capacity"`
		CenterID struct {
			Raw string `json:"raw"`
		} `json:"center.id"`
		CenterLocationLatitude struct {
			Raw float64 `json:"raw"`
		} `json:"center.location.latitude"`
		CenterLocationLongitude struct {
			Raw float64 `json:"raw"`
		} `json:"center.location.longitude"`
		CenterName struct {
			Raw string `json:"raw"`
		} `json:"center.name"`
		ClassCategoryID struct {
			Raw float32 `json:"raw"`
		} `json:"class.category.id"`
		ClassCategoryName struct {
			Raw string `json:"raw"`
		} `json:"class.category.name"`
		ClassCategoryParentID struct {
			Raw float32 `json:"raw"`
		} `json:"class.category.parent_id"`
		ClassDescription struct {
			Raw interface{} `json:"raw"`
		} `json:"class.description"`
		ClassID struct {
			Raw float32 `json:"raw"`
		} `json:"class.id"`
		ClassShowInCatalog struct {
			Raw float32 `json:"raw"`
		} `json:"class.show_in_catalog"`
		ClassTags []struct {
			Description struct {
				Raw string `json:"raw"`
			} `json:"description"`
			ID struct {
				Raw string `json:"raw"`
			} `json:"id"`
			Name struct {
				Raw string `json:"raw"`
			} `json:"name"`
		} `json:"class.tags"`
		ClassType struct {
			Raw float32 `json:"raw"`
		} `json:"class.type"`
		Description struct {
			Raw string `json:"raw"`
		} `json:"description"`
		Duration struct {
			Raw float32 `json:"raw"`
		} `json:"duration"`
		EndTime struct {
			Raw time.Time `json:"raw"`
		} `json:"end_time"`
		EndTimeUtc struct {
			Raw time.Time `json:"raw"`
		} `json:"end_time_utc"`
		GuestPassID struct {
			Raw interface{} `json:"raw"`
		} `json:"guest_pass_id"`
		HasDescription struct {
			Raw string `json:"raw"`
		} `json:"hasDescription"`
		ID struct {
			Raw string `json:"raw"`
		} `json:"id"`
		InstructorID struct {
			Raw string `json:"raw"`
		} `json:"instructor_id"`
		Instructors struct {
			Description struct {
				Raw interface{} `json:"raw"`
			} `json:"description"`
			ID struct {
				Raw string `json:"raw"`
			} `json:"id"`
			ImageURL struct {
				Raw string `json:"raw"`
			} `json:"imageUrl"`
			Name struct {
				Raw string `json:"raw"`
			} `json:"name"`
		} `json:"instructors"`
		IsFreeSession struct {
			Raw string `json:"raw"`
		} `json:"is_free_session"`
		IsInstructorSubstituted struct {
			Raw string `json:"raw"`
		} `json:"is_instructor_substituted"`
		IsVirtualClass struct {
			Raw string `json:"raw"`
		} `json:"is_virtual_class"`
		LastUpdateEventDate struct {
			Raw float64 `json:"raw"`
		} `json:"lastUpdateEventDate"`
		Name struct {
			Raw string `json:"raw"`
		} `json:"name"`
		Occupancy struct {
			Raw float32 `json:"raw"`
		} `json:"occupancy"`
		Price struct {
			Raw float32 `json:"raw"`
		} `json:"price"`
		RoomID struct {
			Raw interface{} `json:"raw"`
		} `json:"room_id"`
		SessionGUID struct {
			Raw string `json:"raw"`
		} `json:"session_guid"`
		SessionID struct {
			Raw float32 `json:"raw"`
		} `json:"session_id"`
		StartTime struct {
			Raw time.Time `json:"raw"`
		} `json:"start_time"`
		StartTimeUtc struct {
			Raw time.Time `json:"raw"`
		} `json:"start_time_utc"`
		Status struct {
			Raw float32 `json:"raw"`
		} `json:"status"`
		WaitListedCount struct {
			Raw float32 `json:"raw"`
		} `json:"wait_listed_count"`
	} `json:"results"`
}

func (c *Client) Initialize() {
	c.Endpoint, _ = url.Parse(openSearchEndpoint)
	c.SearchClient = &http.Client{}

	// Set default headers for all requests
	c.SearchClient.Transport = &headerTransport{
		base: http.DefaultTransport,
		headers: map[string]string{
			"Authorization": "Bearer search-w3zgdjpjefcirgda2acp8itx", // Public Token
		},
	}
}

func (c *Client) Search(startTime time.Time, endTime time.Time, centerIds []string) (*SearchResponse, error) {

	payload := SearchRequest{
		Query: "",
		Page: struct {
			Current int `json:"current"`
			Size    int `json:"size"`
		}{
			Current: 1,
			Size:    999,
		},
		Sort: struct {
			StartTime string `json:"start_time"`
		}{
			StartTime: "asc",
		},
		Filters: struct {
			Any []struct {
				StartTimeUtc struct {
					From time.Time `json:"from"`
					To   time.Time `json:"to"`
				} `json:"start_time_utc"`
			} `json:"any"`
			All []struct {
				CenterID []string `json:"center.id"`
			} `json:"all"`
		}{
			Any: []struct {
				StartTimeUtc struct {
					From time.Time `json:"from"`
					To   time.Time `json:"to"`
				} `json:"start_time_utc"`
			}{
				{
					StartTimeUtc: struct {
						From time.Time `json:"from"`
						To   time.Time `json:"to"`
					}{
						From: startTime,
						To:   endTime,
					},
				},
			},
			All: []struct {
				CenterID []string `json:"center.id"`
			}{
				{
					CenterID: centerIds,
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	baseURL := c.Endpoint
	baseURL.Path = "/api/as/v1/engines/schedule-search-prod/search"

	req, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.SearchClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &searchResponse, nil
}
