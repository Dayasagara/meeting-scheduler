package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Dayasagara/meeting-scheduler/model"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var dateFormat = regexp.MustCompile(`([2]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))`)
var timeFormat = regexp.MustCompile(`(?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d)`)

func ValidateDate(date string) bool {
	return dateFormat.MatchString(date)
}

func ValidateTime(time string) bool {
	return timeFormat.MatchString(time)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func SyncWithGCalendar(events []model.ScheduleEvent) error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return err
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		return err
	}

	var calEvent *calendar.Event
	calendarId := "primary"

	//map the database events to calendar event and sync with google calendars
	for _, event := range events {
		startTime := &calendar.EventDateTime{
			DateTime: event.Date[0:10] + "T" + event.StartingFrom[11:19] + "+05:30",
		}
		endTime := &calendar.EventDateTime{
			DateTime: event.Date[0:10] + "T" + event.EndingTill[11:19] + "+05:30",
		}
		calEvent = &calendar.Event{
			Summary:     event.EventName,
			Description: event.Description,
			Start:       startTime,
			End:         endTime,
			Location:    event.Location,
		}

		//insert events to google calendar
		calEvent, err = srv.Events.Insert(calendarId, calEvent).Do()
		if err != nil {
			return err
		}
	}
	return nil
}
