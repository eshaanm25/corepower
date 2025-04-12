package corepower

import (
	"strings"
	"time"

	"github.com/eshaanm25/corepower/internal/opensearch"
)

type Result struct {
	AvailableSlots    float32   `json:"available_slots"`
	CanBook           string    `json:"can_book"`
	CenterName        string    `json:"center_name"`
	ClassCategoryName string    `json:"class_category_name"`
	StartTimeUtc      time.Time `json:"start_time_utc"`
	CenterID          string    `json:"center_id"`
	SessionID         float32   `json:"session_id"`
}

type ClassPreference struct {
	CenterName string
	StartTime  time.Time
	EndTime    time.Time
	Preference int // Lower number means higher preference
}

// FindIdealClass finds the best available class based on user preferences
func FindIdealClass(searchResponse *opensearch.SearchResponse) *Result {
	// Get Central Time location
	centralTime, _ := time.LoadLocation("America/Chicago")

	// Filter for Yoga Sculpt classes 14 days from now that are available and bookable
	var validClasses []Result
	for _, class := range searchResponse.Results {
		if class.ClassCategoryName.Raw == "Yoga Sculpt" &&
			class.CanBook.Raw == "true" &&
			class.Status.Raw == 2 &&
			class.AvailableSlots.Raw > 0 &&
			class.StartTimeUtc.Raw.After(time.Now().AddDate(0, 0, 13)) {
			result := Result{
				AvailableSlots:    class.AvailableSlots.Raw,
				CanBook:           class.CanBook.Raw,
				CenterName:        strings.Split(class.CenterName.Raw, ` - `)[0],
				ClassCategoryName: class.ClassCategoryName.Raw,
				StartTimeUtc:      class.StartTimeUtc.Raw,
				CenterID:          class.CenterID.Raw,
				SessionID:         class.SessionID.Raw,
			}
			validClasses = append(validClasses, result)
		}
	}

	var preferences []ClassPreference
	var bestClass *Result
	bestPreference := 999

	for _, class := range validClasses {
		// Convert UTC class time to Central Time for comparison
		classTimeCT := class.StartTimeUtc.In(centralTime)
		weekday := classTimeCT.Weekday()
		isWeekend := weekday == time.Saturday || weekday == time.Sunday || weekday == time.Monday || weekday == time.Friday

		if isWeekend {
			// Weekend preferences
			preferences = []ClassPreference{
				{
					CenterName: "Monarch",
					StartTime:  time.Date(0, 1, 1, 13, 0, 0, 0, centralTime), // 1:00 PM CT
					EndTime:    time.Date(0, 1, 1, 15, 0, 0, 0, centralTime), // 3:00 PM CT
					Preference: 1,
				},
				{
					CenterName: "Mueller",
					StartTime:  time.Date(0, 1, 1, 13, 0, 0, 0, centralTime), // 1:00 PM CT
					EndTime:    time.Date(0, 1, 1, 15, 0, 0, 0, centralTime), // 3:00 PM CT
					Preference: 2,
				},
				{
					CenterName: "Triangle",
					StartTime:  time.Date(0, 1, 1, 11, 0, 0, 0, centralTime), // 11:00 AM CT
					EndTime:    time.Date(0, 1, 1, 15, 0, 0, 0, centralTime), // 3:00 PM CT
					Preference: 2,
				},
				{
					CenterName: "Cedar Park",
					StartTime:  time.Date(0, 1, 1, 12, 0, 0, 0, centralTime), // 12:00 PM CT
					EndTime:    time.Date(0, 1, 1, 17, 0, 0, 0, centralTime), // 5:00 PM CT
					Preference: 3,
				},
			}
		} else {
			// Weekday preferences
			preferences = []ClassPreference{
				{
					CenterName: "Triangle",
					StartTime:  time.Date(0, 1, 1, 17, 30, 0, 0, centralTime), // 5:30 PM CT
					EndTime:    time.Date(0, 1, 1, 18, 30, 0, 0, centralTime), // 6:30 PM CT
					Preference: 1,
				},
				{
					CenterName: "Mueller",
					StartTime:  time.Date(0, 1, 1, 18, 0, 0, 0, centralTime),  // 6:00 PM CT
					EndTime:    time.Date(0, 1, 1, 18, 45, 0, 0, centralTime), // 6:45 PM CT
					Preference: 2,
				},
				{
					CenterName: "Monarch",
					StartTime:  time.Date(0, 1, 1, 18, 0, 0, 0, centralTime),  // 6:00 PM CT
					EndTime:    time.Date(0, 1, 1, 18, 45, 0, 0, centralTime), // 6:45 PM CT
					Preference: 2,
				},
				{
					CenterName: "Cedar Park",
					StartTime:  time.Date(0, 1, 1, 18, 30, 0, 0, centralTime), // 6:30 PM CT
					EndTime:    time.Date(0, 1, 1, 19, 30, 0, 0, centralTime), // 7:30 PM CT
					Preference: 3,
				},
			}
		}

		// Get class time of day for comparison (in Central Time)
		classTimeOfDay := time.Date(0, 1, 1,
			classTimeCT.Hour(),
			classTimeCT.Minute(),
			0, 0, centralTime)

		// Check against preferences
		for _, pref := range preferences {
			if strings.Contains(strings.ToLower(class.CenterName), strings.ToLower(pref.CenterName)) &&
				!classTimeOfDay.Before(pref.StartTime) &&
				!classTimeOfDay.After(pref.EndTime) {
				if pref.Preference < bestPreference {
					bestPreference = pref.Preference
					classCopy := class
					bestClass = &classCopy
				}
				break
			}
		}
	}

	return bestClass
}
