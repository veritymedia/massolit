package tasks

import (
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// DetentionNote extends BehaviorNote with additional tracking
type DetentionNote struct {
	BehaviorNote
	DetenionComplete bool `json:"detention_complete"`
}

// getDetentionNotes retrieves incomplete detention notes
func GetDetentionNotes(app *pocketbase.PocketBase) ([]DetentionNote, error) {
	// Ensure the collection exists
	collection, err := app.Dao().FindCollectionByNameOrId("behavior_notes")
	if err != nil {
		return nil, fmt.Errorf("collection not found: %v", err)
	}

	// Prepare query to find detention notes not marked as complete
	results, err := app.Dao().FindRecordsByFilter(collection.Name,
		"behavior_type = 'Detention' AND detention_complete = false", "", 200, 0)
	if err != nil {
		return nil, fmt.Errorf("error querying detention notes: %v", err)
	}

	var detentionNotes []DetentionNote
	for _, record := range results {
		note := DetentionNote{
			BehaviorNote: BehaviorNote{
				ID:                record.GetInt("id"),
				StudentID:         record.GetString("student_id"),
				FirstName:         record.GetString("first_name"),
				LastName:          record.GetString("last_name"),
				Email:             record.GetString("email"),
				Grade:             record.GetString("grade"),
				IncidentTime:      record.GetString("incident_time"),
				BehaviorType:      record.GetString("behavior_type"),
				Notes:             record.GetString("notes"),
				NextStep:          record.GetString("next_step"),
				NextStepDate:      record.GetString("next_step_date"),
				AuthorID:          record.GetInt("author_id"),
				ReportedBy:        record.GetString("reported_by"),
				HomeRoomAdvisor:   record.GetString("homeroom_advisor"),
				VisibleToParents:  record.GetBool("visible_to_parents"),
				VisibleToStudents: record.GetBool("visible_to_students"),
				CreatedAt:         record.GetString("created_at"),
				UpdatedAt:         record.GetString("updated_at"),
			},
			DetenionComplete: record.GetBool("detention_complete"),
		}
		detentionNotes = append(detentionNotes, note)
	}

	return detentionNotes, nil
}

// generateDetentionReportHTML creates an HTML report for detention notes
func generateDetentionReportHTML(notes []DetentionNote) string {
	var htmlBody string

	htmlBody += fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Detention Notes Report</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .header { background-color: #4CAF50; color: white; padding: 10px; text-align: center; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Outstanding Detention Notes</h1>
        <p>Generated on: %s</p>
    </div>

    <table>
        <thead>
            <tr>
                <th>Student</th>
                <th>Grade</th>
                <th>Incident Time</th>
                <th>Reported By</th>
                <th>Notes</th>
                <th>Next Step</th>
            </tr>
        </thead>
        <tbody>
`, time.Now().Format("2006-01-02 15:04:05"))

	for _, note := range notes {
		htmlBody += fmt.Sprintf(`
            <tr>
                <td>%s %s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
            </tr>
`, note.FirstName, note.LastName, note.Grade, note.IncidentTime,
			note.ReportedBy, note.Notes, note.NextStep)
	}

	htmlBody += fmt.Sprintf(`
        </tbody>
    </table>

    <p>Total Outstanding Detentions: %d</p>
</body>
</html>
`, len(notes))

	return htmlBody
}

// sendDetentionReport sends a detention report using PocketBase mailer
func SendDetentionReport(app *pocketbase.PocketBase, notes []DetentionNote) error {
	// Get the report recipient from environment or app settings
	recipient := os.Getenv("DETENTION_REPORT_RECIPIENT")
	if recipient == "" {
		return fmt.Errorf("no recipient email set")
	}

	// Generate HTML report
	htmlBody := generateDetentionReportHTML(notes)

	// Prepare the message
	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To: []mail.Address{
			{Address: recipient},
		},
		Subject: fmt.Sprintf("Detention Notes Report - %s", time.Now().Format("2006-01-02")),
		HTML:    htmlBody,
	}

	// Send the email
	return app.NewMailClient().Send(message)
}
