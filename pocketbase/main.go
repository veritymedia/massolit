package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/veritymedia/massolit/pocketbase/tasks"
)

type ConfigRecord struct {
	id    string
	name  string
	value string
}

//go:embed all:.output/public
var public embed.FS

func main() {
	app := pocketbase.New()

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	managebacUrl := "https://api.managebac.com/v2"
	managebacApiKey := os.Getenv("MANAGEBAC_API")

	if len(managebacApiKey) == 0 {
		log.Panic("No Managebac Key has been found. Exiting.")
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		// Add a job to send detention reports daily at 8 AM
		scheduler.MustAdd("sendReport", "0 8 * * *", func() {
			// Fetch outstanding detention notes
			notes, err := tasks.GetDetentionNotes(app)
			if err != nil {
				log.Printf("Error fetching detention notes: %v", err)
				return
			}

			// If there are outstanding detentions, send report
			if len(notes) > 0 {
				if err := tasks.SendDetentionReport(app, notes); err != nil {
					log.Printf("Error sending detention report: %v", err)
				}
			}
		})

		if err != nil {
			return fmt.Errorf("failed to add detention report cron job: %v", err)
		}

		// Start the scheduler
		scheduler.Start()

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		go func() {
			for {

				modifiedSinceRecord, err := e.App.Dao().FindRecordsByFilter("config", "name='last_behavior_sync_datetime'", "", 1, 0)

				if err != nil {
					fmt.Println("CRON::BEHAVIOUR_NOTES Error: Could not fetch latest edit time")
				}

				if len(modifiedSinceRecord) == 0 {
					fmt.Println("CRON::BEHAVIOUR_NOTES Error: Could not fetch, length 0")

				}

				var modifiedSinceValue string = modifiedSinceRecord[0].GetString("value")

				fmt.Printf("CRON::BEHAVIOUR_NOTES Last modified %s\n", modifiedSinceValue)

				resp, err := tasks.FetchBehaviorNotes(managebacApiKey, modifiedSinceValue, managebacUrl)

				if err != nil {
					log.Printf("Error fetching behavior notes: %v", err)
				} else {

					if err := tasks.SaveBehaviorNotes(app, resp.BehaviorNotes); err != nil {
						log.Printf("Error saving behavior notes: %v", err)
					}

					updatedSinceRecord := modifiedSinceRecord[0]
					updatedSinceRecord.Set("value", time.Now().Format(time.RFC3339))

					if err := app.Dao().SaveRecord(updatedSinceRecord); err != nil {
						fmt.Printf("Could not save new modified_since param: %s", err)
					}
				}

				time.Sleep(5 * time.Minute)
			}
		}()
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.GET("/managebac/students", func(c echo.Context) error {

			searchParam := c.QueryParam("q")
			resource := "/v2/students"
			params := url.Values{}
			params.Add("per_page", "400")
			params.Add("q", searchParam)

			u, _ := url.ParseRequestURI(managebacUrl)
			u.Path = resource
			u.RawQuery = params.Encode()
			qualifiedUrl := fmt.Sprintf("%v", u)
			req, err := http.NewRequest(http.MethodGet, qualifiedUrl, nil)

			fmt.Println("URL: ", qualifiedUrl)

			if err != nil {
				fmt.Println(err)
				c.JSON(500, map[string]string{
					"message": "Could not create request",
				})
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("auth-token", managebacApiKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				// Handle the error if the request fails
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reach external server"})
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// Handle error reading the body
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
			}

			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				// Handle the case where the response is not valid JSON
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid JSON response from server"})
			}

			// Log the external response for debugging (optional)
			log.Printf("Response from external API: %v", jsonResponse)

			// Return the decoded JSON as a response
			return c.JSON(resp.StatusCode, jsonResponse)
		})

		e.Router.GET("/managebac/students/:studentId", func(c echo.Context) error {

			studentId := c.PathParam("studentId")

			req, err := http.NewRequest(http.MethodGet, fmt.Sprint(managebacUrl, "/students/", studentId), nil)

			if err != nil {
				fmt.Println(err)
				c.JSON(500, map[string]string{
					"message": "Could not create request",
				})
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("auth-token", managebacApiKey)

			fmt.Println("API KEY: ", managebacApiKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				// Handle the error if the request fails
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reach external server"})
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// Handle error reading the body
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
			}

			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				// Handle the case where the response is not valid JSON
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid JSON response from server"})
			}

			// Log the external response for debugging (optional)
			log.Printf("Response from external API: %s", body)

			// Return the response from the external server
			// Return the decoded JSON as a response
			return c.JSON(resp.StatusCode, jsonResponse)
		})

		e.Router.GET("/*", apis.StaticDirectoryHandler(echo.MustSubFS(public, ".output/public"), true))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
