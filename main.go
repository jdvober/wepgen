package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jdvober/gauth"
	"github.com/jdvober/gsheets"
	"github.com/joho/godotenv"
)

// Student contains all imformation about a student relating to WEP Status
type Student struct {
	firstName string
	lastName  string
	name      string
	AP        string
	honors    string
	grade     string
	target    string
	start     string
	end       string
	review    string
	superCog  string
	creative  string
	science   string
	academic  string
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("config.env"); err != nil {
		_, err := os.Create("./config.env")
		if err != nil {
			fmt.Println(err)
		}
		log.Println("Error loading .env")
	}
}
func main() {
	// Load data from env file
	giftedSSID := os.Getenv("GIFTED_SSID")
	wepTemplateSSID := os.Getenv("WEP_TEMPLATE_SSID")
	gr := os.Getenv("GRADE")
	grad := os.Getenv("GRADUATION_YEAR")
	st := os.Getenv("START")
	e := os.Getenv("END")
	rev := os.Getenv("REVIEW")

	client := gauth.Authorize()

	data := gsheets.GetValues(client, giftedSSID, "Sheet1!A2:O")

	// Iterate through data and format for paste in WEP Template Google Sheet
	values := make([][]interface{}, len(data))
	for i := range data {
		// Create student profile
		profile := Student{
			firstName: data[i][1].(string),
			lastName:  data[i][0].(string),
			superCog:  data[i][2].(string),
			science:   data[i][12].(string),
			name:      strings.Join([]string{data[i][1].(string), data[i][0].(string)}, " "),
			grade:     gr,
			target:    grad,
			start:     st,
			end:       e,
			review:    rev,
		}
		if gr == "12" {
			profile.AP = "true"
			profile.honors = "false"
		} else {
			profile.AP = "false"
			profile.honors = "true"
		}

		// Calculate Academic
		switch {
		case len(profile.superCog) > 0 && len(profile.science) > 0:
			profile.academic = "Superior Cognitive Ability, Superior Academic Ability"
		case len(profile.superCog) > 0:
			profile.academic = "Superior Cognitive Ability"
		case len(profile.science) > 0:
			profile.academic = "Superior Academic Ability"
		}
		if len(profile.creative) > 0 {
			profile.academic = strings.Join([]string{profile.academic, "Creative"}, ", ")
		}
		// Add to array to post to sheet
		values[i] = []interface{}{profile.AP, profile.name, profile.grade, profile.target, profile.start, profile.end, profile.review, profile.academic}
	}
	gsheets.BatchUpdateValues(client, wepTemplateSSID, "Sheet1!A2:H", "ROWS", values)
}
