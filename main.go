package main

import (
	"fmt"
	"io/ioutil"
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
	// wepTemplateSSID := os.Getenv("WEP_TEMPLATE_SSID")
	gr := os.Getenv("GRADE")
	grad := os.Getenv("GRADUATION_YEAR")
	st := os.Getenv("START")
	e := os.Getenv("END")
	rev := os.Getenv("REVIEW")

	client := gauth.Authorize()

	data := gsheets.GetValues(client, giftedSSID, "Sheet1!A2:O")

	// Iterate through data and format for paste in WEP Template Google Sheet
	values := make([][]interface{}, len(data))
	profiles := make([]Student, len(data))
	for i := range data {
		// Create student profile
		profile := Student{
			firstName: data[i][1].(string),
			lastName:  data[i][0].(string),
			superCog:  data[i][2].(string),
			creative:  data[i][3].(string),
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

		profiles = append(profiles, profile)
		// Make local WEP XML file
	}
	makeFiles(profiles)
	// Post to sheet
	// gsheets.BatchUpdateValues(client, wepTemplateSSID, "Sheet1!A2:H", "ROWS", values)
}

func makeFiles(profiles []Student) {

	for _, profile := range profiles {
		var filename string = profile.firstName + "_" + profile.lastName + "_" + "WEP.xml"

		if _, err := os.Stat("./weps"); os.IsNotExist(err) {
			fmt.Println("./weps does not exist. Mkdir will create ./weps")
			os.Mkdir("./weps/", 0777)
		}
		// Open the proper template as a byte array
		if profile.AP == "true" {
			src, err := ioutil.ReadFile("./templates/AP.xml")
			if err != nil {
				log.Println("Error reading file ./templates/AP.XML")
				return
			}
			// Find and replace
			newFile := strings.ReplaceAll(string(src), "$NAME$", profile.name)
			newFile = strings.ReplaceAll(newFile, "$GRADE$", profile.grade)
			newFile = strings.ReplaceAll(newFile, "$TARGET$", profile.target)
			newFile = strings.ReplaceAll(newFile, "$START$", profile.start)
			newFile = strings.ReplaceAll(newFile, "$END$", profile.end)
			newFile = strings.ReplaceAll(newFile, "$REVIEW$", profile.review)
			newFile = strings.ReplaceAll(newFile, "$ACADEMIC$", profile.academic)

			f, err := os.Create("./weps/" + filename)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("\nCreated file ./weps/%s\n", filename)
			l, err := f.WriteString(newFile)
			if err != nil {
				fmt.Printf("Problem writing %v bytes", l)
				log.Println(err)
				f.Close()
				return
			}
			err = f.Close()
			if err != nil {
				fmt.Println(err)
			}
		} else if profile.honors == "true" {
			src, err := ioutil.ReadFile("./templates/Honors.xml")
			if err != nil {
				log.Println("Error reading file ./templates/Honors.XML")
				return
			}
			// Find and replace
			newFile := strings.ReplaceAll(string(src), "$NAME$", profile.name)
			newFile = strings.ReplaceAll(newFile, "$GRADE$", profile.grade)
			newFile = strings.ReplaceAll(newFile, "$TARGET$", profile.target)
			newFile = strings.ReplaceAll(newFile, "$START$", profile.start)
			newFile = strings.ReplaceAll(newFile, "$END$", profile.end)
			newFile = strings.ReplaceAll(newFile, "$REVIEW$", profile.review)
			newFile = strings.ReplaceAll(newFile, "$ACADEMIC$", profile.academic)

			f, err := os.Create("./weps/" + filename)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("\nCreated file ./weps/%s\n", filename)
			l, err := f.WriteString(newFile)
			if err != nil {
				fmt.Printf("Problem writing %v bytes", l)
				log.Println(err)
				f.Close()
				return
			}
			err = f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
