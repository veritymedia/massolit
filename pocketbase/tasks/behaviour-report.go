package tasks

import (
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type DetentionNote struct {
	BehaviorNote
	DetenionComplete bool `json:"detention_complete"`
}

func HandleDetentionReportSend(app *pocketbase.PocketBase) error {
	mailListRecord, err := app.Dao().FindRecordsByFilter("mail_list", "subs~'behavior'", "", 100, 0)

	if err != nil {
		return fmt.Errorf("Could not find mail_list records")
	}

	recipients := []mail.Address{}

	for _, record := range mailListRecord {
		recipients = append(recipients, mail.Address{Address: record.GetString("email")})
	}

	notes, err := GetDetentionNotes(app)
	if err != nil {
		return fmt.Errorf("Error fetching detention notes: %v", err)
	}

	if len(notes) > 0 {
		if err := SendDetentionReport(app, notes); err != nil {
			log.Printf("Error sending detention report: %v", err)
		}
	}
	return err
}

func GetDetentionNotes(app *pocketbase.PocketBase) ([]DetentionNote, error) {
	collection, err := app.Dao().FindCollectionByNameOrId("behavior_notes")
	if err != nil {
		return nil, fmt.Errorf("collection not found: %v", err)
	}

	results, err := app.Dao().FindRecordsByFilter(collection.Name,
		"next_step ~ 'Detention' && action_complete = false", "-created_at", 200, 0)
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
func PrettyFormatDate(dateStr string) string {
	parsedTime, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		return "Invalid date"
	}

	now := time.Now()

	daysDiff := int(now.Sub(parsedTime).Hours() / 24)

	month := parsedTime.Format("January")
	day := parsedTime.Day()

	ordinalSuffix := func(day int) string {
		switch {
		case day%10 == 1 && day%100 != 11:
			return "st"
		case day%10 == 2 && day%100 != 12:
			return "nd"
		case day%10 == 3 && day%100 != 13:
			return "rd"
		default:
			return "th"
		}
	}

	dateFormat := fmt.Sprintf("%s %d%s", month, day, ordinalSuffix(day))

	if daysDiff > 0 {
		dateFormat += fmt.Sprintf(" (%d day%s ago)", daysDiff, func() string {
			if daysDiff != 1 {
				return "s"
			}
			return ""
		}())
	}

	return dateFormat
}

func generateDetentionReportHTML(notes []DetentionNote) string {
	var htmlBody string

	htmlBody += fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Detention Report</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; background-color: #f9f9f9; color: #333; padding: 20px;">
    <div class="header" style="background-color: #232363; color: white; padding: 20px; text-align: center; border-radius: 8px;">
        <h1>Outstanding Detentions</h1>
        <p>Generated on: %s</p>
        <a style="color: white; font-weight: bold; font-size: 1em;" href="https://massolit.app/behavior">Massolit Detentions</a>
    </div>

    <table style="width: 100%; border-collapse: collapse; margin-top: 20px; border-radius: 8px; overflow: hidden;">
        <thead>
            <tr>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Student</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Grade</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Incident Date</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Reported By</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Notes</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left; background-color: #232363; color: white; text-transform: uppercase; font-weight: bold;">Next Step</th>
            </tr>
        </thead>
        <tbody>
`, time.Now().Format("2006-01-02 15:04:05"))

	for i, note := range notes {
		var bgColor string
		if i%2 == 0 {
			bgColor = "#f2f8fc"
		} else {
			bgColor = "#ffffff"
		}
		htmlBody += fmt.Sprintf(`
            <tr style="background-color: %s;">
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s %s</td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s</td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s</td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s</td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s</td>
                <td style="border: 1px solid #ddd; padding: 12px; text-align: left;">%s</td>
            </tr>
`, bgColor, note.FirstName, note.LastName, note.Grade, PrettyFormatDate(note.IncidentTime), note.ReportedBy, note.Notes, note.NextStep)
	}

	htmlBody += fmt.Sprintf(`
        </tbody>
    </table>

    <p>Total Pending Detentions: %d</p>
</body>
</html>
`, len(notes))

	return htmlBody
}

func SendDetentionReport(app *pocketbase.PocketBase, notes []DetentionNote) error {
	mailListRecord, err := app.Dao().FindRecordsByFilter("mail_list", "subs~'behavior'", "", 100, 0)

	if err != nil {
		return fmt.Errorf("Could not find mail_list records")
	}

	recipients := []mail.Address{}

	for _, record := range mailListRecord {
		recipients = append(recipients, mail.Address{Address: record.GetString("email")})
	}

	htmlBody := generateDetentionReportHTML(notes)

	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      recipients,
		Subject: fmt.Sprintf("Detention Report - %s", time.Now().Format("2006-01-02")),
		HTML:    htmlBody,
	}

	return app.NewMailClient().Send(message)
}
