package tasks

import (
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type DetentionNote struct {
	BehaviorNote
	DetenionComplete bool `json:"detention_complete"`
}

func HandleDetentionReportSend(app *pocketbase.PocketBase) error {
	fmt.Println("CRON::DETENTION_REPORT Starting detention report generation")
	
	mailListRecord, err := app.Dao().FindRecordsByFilter("mail_list", "subs~'behavior'", "", 100, 0)
	if err != nil {
		return fmt.Errorf("Could not find mail_list records: %v", err)
	}

	recipients := []mail.Address{}
	for _, record := range mailListRecord {
		recipients = append(recipients, mail.Address{Address: record.GetString("email")})
	}
	
	fmt.Printf("CRON::DETENTION_REPORT Found %d email recipients\n", len(recipients))

	// Get regular detention notes
	notes, err := GetDetentionNotes(app)
	if err != nil {
		return fmt.Errorf("Error fetching detention notes: %v", err)
	}

	// Get double detention alerts
	doubleDetentions, err := CheckDoubleDetentions(app)
	if err != nil {
		return fmt.Errorf("Error checking double detentions: %v", err)
	}

	// Send report if there are any detentions or double detention alerts
	if len(notes) > 0 || len(doubleDetentions) > 0 {
		fmt.Printf("CRON::DETENTION_REPORT Sending report for %d detention notes and %d double detention alerts\n", len(notes), len(doubleDetentions))
		if err := SendEnhancedDetentionReport(app, notes, doubleDetentions); err != nil {
			log.Printf("Error sending detention report: %v", err)
			return err
		}
		fmt.Println("CRON::DETENTION_REPORT Successfully sent detention report")
	} else {
		fmt.Println("CRON::DETENTION_REPORT No pending detentions or double detention alerts found")
	}
	
	return nil
}

func GetDetentionNotes(app *pocketbase.PocketBase) ([]DetentionNote, error) {
	collection, err := app.Dao().FindCollectionByNameOrId("behavior_notes")
	if err != nil {
		return nil, fmt.Errorf("collection not found: %v", err)
	}

	// Calculate the date 7 days ago for rolling window check
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	
	// Query for detention notes within the last 7 days that are not action_complete
	filter := fmt.Sprintf("next_step ~ 'Detention' && action_complete = false && created_at >= '%s'", sevenDaysAgo)
	
	results, err := app.Dao().FindRecordsByFilter(collection.Name, filter, "-created_at", 200, 0)
	if err != nil {
		return nil, fmt.Errorf("error querying detention notes: %v", err)
	}

	fmt.Printf("Found %d detention notes within the last 7 days\n", len(results))

	var detentionNotes []DetentionNote
	for _, record := range results {
		// Convert author_id from string to int
		authorID := 0
		if authorIDStr := record.GetString("author_id"); authorIDStr != "" {
			if id, err := strconv.Atoi(authorIDStr); err == nil {
				authorID = id
			}
		}

		note := DetentionNote{
			BehaviorNote: BehaviorNote{
				ID:                0, // Not using ManageBac ID for grouping, use PocketBase record ID
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
				AuthorID:          authorID,
				ReportedBy:        record.GetString("reported_by"),
				HomeRoomAdvisor:   record.GetString("homeroom_advisor"),
				VisibleToParents:  record.GetString("visible_to_parents") == "true",
				VisibleToStudents: record.GetString("visible_to_students") == "true",
				CreatedAt:         record.GetString("created_at"),
				UpdatedAt:         record.GetString("updated_at"),
				ActionComplete:    record.GetBool("action_complete"),
			},
			DetenionComplete: record.GetBool("detention_complete"),
		}
		detentionNotes = append(detentionNotes, note)
	}

	return detentionNotes, nil
}

// GetAllDetentionNotesInWindow gets ALL detention notes (completed or not) within 7-day window
func GetAllDetentionNotesInWindow(app *pocketbase.PocketBase) ([]DetentionNote, error) {
	collection, err := app.Dao().FindCollectionByNameOrId("behavior_notes")
	if err != nil {
		return nil, fmt.Errorf("collection not found: %v", err)
	}

	// Calculate the date 7 days ago for rolling window check
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	
	// Query for ALL detention notes within the last 7 days (including lunchtime detention, completed or not)
	filter := fmt.Sprintf("(next_step ~ 'Detention' || next_step ~ 'Lunchtime') && created_at >= '%s'", sevenDaysAgo)
	
	results, err := app.Dao().FindRecordsByFilter(collection.Name, filter, "-created_at", 500, 0)
	if err != nil {
		return nil, fmt.Errorf("error querying all detention notes: %v", err)
	}

	var detentionNotes []DetentionNote
	for _, record := range results {
		// Convert author_id from string to int
		authorID := 0
		if authorIDStr := record.GetString("author_id"); authorIDStr != "" {
			if id, err := strconv.Atoi(authorIDStr); err == nil {
				authorID = id
			}
		}

		note := DetentionNote{
			BehaviorNote: BehaviorNote{
				ID:                0, // Not using ManageBac ID for grouping, use PocketBase record ID
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
				AuthorID:          authorID,
				ReportedBy:        record.GetString("reported_by"),
				HomeRoomAdvisor:   record.GetString("homeroom_advisor"),
				VisibleToParents:  record.GetString("visible_to_parents") == "true",
				VisibleToStudents: record.GetString("visible_to_students") == "true",
				CreatedAt:         record.GetString("created_at"),
				UpdatedAt:         record.GetString("updated_at"),
				ActionComplete:    record.GetBool("action_complete"),
			},
			DetenionComplete: record.GetBool("detention_complete"),
		}
		detentionNotes = append(detentionNotes, note)
	}

	return detentionNotes, nil
}

// StudentDetentionSummary represents a student with multiple detentions
type StudentDetentionSummary struct {
	StudentID       string
	FirstName       string
	LastName        string
	Email           string
	Grade           string
	DetentionCount  int
	DetentionNotes  []DetentionNote
}

// CheckDoubleDetentions identifies students with multiple detentions in 7-day window
func CheckDoubleDetentions(app *pocketbase.PocketBase) ([]StudentDetentionSummary, error) {
	allDetentions, err := GetAllDetentionNotesInWindow(app)
	if err != nil {
		return nil, fmt.Errorf("error getting all detention notes: %v", err)
	}

	// Group detentions by student
	studentDetentions := make(map[string][]DetentionNote)
	for _, detention := range allDetentions {
		studentID := detention.StudentID
		studentDetentions[studentID] = append(studentDetentions[studentID], detention)
	}

	var doubleDetentions []StudentDetentionSummary
	for studentID, detentions := range studentDetentions {
		if len(detentions) >= 2 { // Student has 2 or more detentions
			summary := StudentDetentionSummary{
				StudentID:      studentID,
				FirstName:      detentions[0].FirstName,
				LastName:       detentions[0].LastName,
				Email:          detentions[0].Email,
				Grade:          detentions[0].Grade,
				DetentionCount: len(detentions),
				DetentionNotes: detentions,
			}
			doubleDetentions = append(doubleDetentions, summary)
		}
	}

	fmt.Printf("Found %d students with multiple detentions in 7-day window\n", len(doubleDetentions))
	for _, student := range doubleDetentions {
		fmt.Printf("  - %s %s (%s): %d detentions\n", 
			student.FirstName, student.LastName, student.Grade, student.DetentionCount)
	}

	return doubleDetentions, nil
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

func generateEnhancedDetentionReportHTML(notes []DetentionNote, doubleDetentions []StudentDetentionSummary) string {
	var htmlBody string

	htmlBody += fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Detention Report</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 1000px; margin: 0 auto; background-color: #f9f9f9; color: #333; padding: 20px;">
    <div class="header" style="background-color: #232363; color: white; padding: 20px; text-align: center; border-radius: 8px;">
        <h1>Outstanding Detentions Report</h1>
        <p>Generated on: %s</p>
        <a style="color: white; font-weight: bold; font-size: 1em;" href="https://massolit.app/behavior">Massolit Detentions</a>
    </div>
`, time.Now().Format("2006-01-02 15:04:05"))

	// Add double detention alerts section if any exist
	if len(doubleDetentions) > 0 {
		htmlBody += `
    <div class="alert-section" style="background-color: #ffebee; border: 2px solid #f44336; border-radius: 8px; padding: 20px; margin: 20px 0;">
        <h2 style="color: #d32f2f; margin-top: 0;">🚨 DOUBLE DETENTION ALERTS</h2>
        <p style="color: #d32f2f; font-weight: bold;">The following students have received multiple detentions within the last 7 days:</p>
        <table style="width: 100%; border-collapse: collapse; margin-top: 10px;">
            <thead>
                <tr>
                    <th style="border: 1px solid #f44336; padding: 8px; text-align: left; background-color: #f44336; color: white;">Student</th>
                    <th style="border: 1px solid #f44336; padding: 8px; text-align: left; background-color: #f44336; color: white;">Grade</th>
                    <th style="border: 1px solid #f44336; padding: 8px; text-align: left; background-color: #f44336; color: white;">Total Detentions</th>
                    <th style="border: 1px solid #f44336; padding: 8px; text-align: left; background-color: #f44336; color: white;">Detention Types</th>
                </tr>
            </thead>
            <tbody>`

		for _, student := range doubleDetentions {
			detentionTypes := make(map[string]int)
			for _, detention := range student.DetentionNotes {
				detentionTypes[detention.NextStep]++
			}
			
			var typesList string
			for detType, count := range detentionTypes {
				if typesList != "" {
					typesList += ", "
				}
				typesList += fmt.Sprintf("%s (%d)", detType, count)
			}

			htmlBody += fmt.Sprintf(`
                <tr style="background-color: #ffcdd2;">
                    <td style="border: 1px solid #f44336; padding: 8px;">%s %s</td>
                    <td style="border: 1px solid #f44336; padding: 8px;">%s</td>
                    <td style="border: 1px solid #f44336; padding: 8px; font-weight: bold;">%d</td>
                    <td style="border: 1px solid #f44336; padding: 8px;">%s</td>
                </tr>`,
				student.FirstName, student.LastName, student.Grade, student.DetentionCount, typesList)
		}

		htmlBody += `
            </tbody>
        </table>
    </div>`
	}

	// Add regular detention table if any exist
	if len(notes) > 0 {
		htmlBody += `
    <h2 style="margin-top: 30px;">Outstanding Detentions</h2>
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
        <tbody>`

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

		htmlBody += `
        </tbody>
    </table>`
	}

	// Summary section
	htmlBody += fmt.Sprintf(`
    <div class="summary" style="margin-top: 30px; padding: 15px; background-color: #e3f2fd; border-radius: 8px;">
        <h3 style="margin-top: 0;">Summary</h3>
        <p><strong>Total Pending Detentions:</strong> %d</p>
        <p><strong>Students with Multiple Detentions:</strong> %d</p>
    </div>
</body>
</html>
`, len(notes), len(doubleDetentions))

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

func SendEnhancedDetentionReport(app *pocketbase.PocketBase, notes []DetentionNote, doubleDetentions []StudentDetentionSummary) error {
	mailListRecord, err := app.Dao().FindRecordsByFilter("mail_list", "subs~'behavior'", "", 100, 0)

	if err != nil {
		return fmt.Errorf("Could not find mail_list records")
	}

	recipients := []mail.Address{}

	for _, record := range mailListRecord {
		recipients = append(recipients, mail.Address{Address: record.GetString("email")})
	}

	htmlBody := generateEnhancedDetentionReportHTML(notes, doubleDetentions)

	// Create subject line that indicates if there are double detention alerts
	subject := fmt.Sprintf("Detention Report - %s", time.Now().Format("2006-01-02"))
	if len(doubleDetentions) > 0 {
		subject = fmt.Sprintf("🚨 DETENTIONS: Detention Report with %d Double Detention Alerts - %s", len(doubleDetentions), time.Now().Format("2006-01-02"))
	}

	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      recipients,
		Subject: subject,
		HTML:    htmlBody,
	}

	return app.NewMailClient().Send(message)
}

// Keep the original function for backward compatibility
func generateDetentionReportHTML(notes []DetentionNote) string {
	// Call the enhanced version with empty double detentions
	return generateEnhancedDetentionReportHTML(notes, []StudentDetentionSummary{})
}