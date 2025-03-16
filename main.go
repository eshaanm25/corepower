package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/eshaanm25/corepower/internal/cognito"
	"github.com/eshaanm25/corepower/internal/corepower"
	"github.com/eshaanm25/corepower/internal/opensearch"
)

func main() {
	// Define command line flags
	username := flag.String("username", "", "CorePower username")
	password := flag.String("password", "", "CorePower password")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Both -username and -password flags are required")
	}

	// Get Central Time location
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}

	// Search For Classes

	// Get start and end times
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(0, 0, 14)

	// Get center IDs
	centerIds := []string{
		"96efaaaa-b040-4d12-8829-51317dd8c1c2", // Cedar Park
		"d810c305-2816-4a92-8c3d-47fe0a69d63a", // Monarch
		"a9d80dd6-0610-4584-b945-f297122d6268", // Mueller
		"5262c2d4-7f8f-4d15-9adf-5ba107409b30", // Triangle
	}

	// Initialize the search client
	searchClient := &opensearch.Client{}
	searchClient.Initialize()

	// Call the search function
	log.Println("Searching for available classes...")
	res, err := searchClient.Search(startTime, endTime, centerIds)
	if err != nil {
		log.Fatalf("Error during search: %v", err)
	}
	log.Printf("Found %d classes...", len(res.Results))

	// Find ideal class based on preferences
	log.Println("Finding ideal class based on preferences...")
	idealClass := corepower.FindIdealClass(res)
	if idealClass == nil {
		log.Println("No ideal class found matching preferences")
		return
	} else {
		// Print ideal class details
		log.Printf("Found ideal class: %s in %s at %s. The available slots are %.0f\n",
			idealClass.ClassCategoryName,
			idealClass.CenterName,
			idealClass.StartTimeUtc.In(location).Format("Mon Jan 2 3:04 PM MST"),
			idealClass.AvailableSlots)

		// Reserve a Class
		log.Println("Attempting to book the class...")

		// Initialize the Cognito token manager with provided credentials
		ctm := cognito.NewCognitoTokenManager(*username, *password)

		// Authenticate to get the token
		err = ctm.Authenticate(context.Background())
		if err != nil {
			log.Fatalf("Error authenticating with CorePower API: %v", err)
		}
		log.Println("Successfully authenticated with CorePower API")

		// Initialize the Reservations Client
		corePowerClient := &corepower.Client{
			Token: ctm.GetIdToken(),
		}
		corePowerClient.Initialize()

		// Reserve the class
		log.Printf("Sending reservation request for %s at %s in %s...",
			idealClass.ClassCategoryName,
			idealClass.StartTimeUtc.In(location).Format("Mon Jan 2 3:04 PM MST"),
			idealClass.CenterName)
		response, err := corePowerClient.Reserve(idealClass.CenterID, idealClass.SessionID)
		if err != nil {
			log.Fatalf("Error reserving: %v\n", err)
		} else {
			log.Printf("Successfully booked: %s at %s in %s", response.ClassName, response.StartTimeUTC, response.Center)
		}

		log.Println("Reservation process completed")
	}
}
