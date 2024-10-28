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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

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
