package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

// BehaviorNote represents the structure of behavior notes from ManageBac
type BehaviorNote struct {
	ID                int    `json:"id"`
	StudentID         string `json:"student_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Grade             string `json:"grade"`
	IncidentTime      string `json:"incident_time"`
	BehaviorType      string `json:"behavior_type"`
	Notes             string `json:"notes"`
	NextStep          string `json:"next_step"`
	NextStepDate      string `json:"next_step_date"`
	AuthorID          int    `json:"author_id"`
	ReportedBy        string `json:"reported_by"`
	HomeRoomAdvisor   string `json:"homeroom_advisor"`
	VisibleToParents  bool   `json:"visible_to_parents"`
	VisibleToStudents bool   `json:"visible_to_students"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ActionComplete    bool   `json:"action_complete"`
}

// ManageBacResponse represents the full API response
type ManageBacResponse struct {
	BehaviorNotes []BehaviorNote `json:"behavior_notes"`
	Meta          struct {
		CurrentPage int `json:"current_page"`
		TotalPages  int `json:"total_pages"`
		TotalCount  int `json:"total_count"`
		PerPage     int `json:"per_page"`
	} `json:"meta"`
}

// fetchBehaviorNotes retrieves behavior notes from ManageBac API
func FetchBehaviorNotes(authToken string, modifiedSince string, managebacUrl string) (*ManageBacResponse, error) {
	fmt.Println("CRON::BEHAVIOUR_NOTES::FETCH_BEHAVIOUR_NOTES")

	resource := "/v2/behavior/notes"

	queryParams := []string{}
	var rawQuery string

	if modifiedSince != "" {
		queryParams = append(queryParams, fmt.Sprintf("modified_since=%s", url.PathEscape(modifiedSince)))
	}

	queryParams = append(queryParams, "per_page=150")

	rawQuery = strings.Join(queryParams, "&")

	u, _ := url.ParseRequestURI(managebacUrl)
	u.Path = resource
	u.RawQuery = rawQuery
	qualifiedUrl := fmt.Sprintf("%v", u)

	req, err := http.NewRequest(http.MethodGet, qualifiedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("auth-token", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var manageBacResp ManageBacResponse
	err = json.NewDecoder(resp.Body).Decode(&manageBacResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &manageBacResp, nil
}

func SaveBehaviorNotes(app *pocketbase.PocketBase, notes []BehaviorNote) error {
	fmt.Println("CRON::BEHAVIOUR_NOTES::SAVE_BEHAVIOUR_NOTES")
	collection, err := app.Dao().FindCollectionByNameOrId("behavior_notes")
	if err != nil {
		return fmt.Errorf("collection not found: %v", err)
	}

	for _, note := range notes {
		// Check if note already exists
		existingNote, err := app.Dao().FindRecordById(
			collection.Name,
			fmt.Sprintf("id = '%d'", note.ID),
		)

		if err != nil {
			// Create new record if not exists
			record := models.NewRecord(collection)
			record.Set("id", fmt.Sprintf("%d", note.ID))
			record.Set("student_id", note.StudentID)
			record.Set("first_name", note.FirstName)
			record.Set("last_name", note.LastName)
			record.Set("email", note.Email)
			record.Set("grade", note.Grade)
			record.Set("incident_time", note.IncidentTime)
			record.Set("behavior_type", note.BehaviorType)
			record.Set("notes", note.Notes)
			record.Set("next_step", note.NextStep)
			record.Set("next_step_date", note.NextStepDate)
			record.Set("author_id", note.AuthorID)
			record.Set("reported_by", note.ReportedBy)
			record.Set("homeroom_advisor", note.HomeRoomAdvisor)
			record.Set("visible_to_parents", note.VisibleToParents)
			record.Set("visible_to_students", note.VisibleToStudents)
			record.Set("created_at", note.CreatedAt)
			record.Set("updated_at", note.UpdatedAt)

			if err := app.Dao().Save(record); err != nil {
				log.Printf("Error saving record: %v", err)
			}
		} else {
			// Update existing record
			existingNote.Set("updated_at", note.UpdatedAt)

			if err := app.Dao().Save(existingNote); err != nil {
				log.Printf("Error updating record: %v", err)
			}
		}
	}

	return nil
}
