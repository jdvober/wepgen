import (
	"github.com/nguyenthenguyen/docx"
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
		profiles = append(profiles, profile)
	}
	makeFiles(profiles)
}

func makeFiles(profiles []Student) {

	for _, profile := range profiles {
		var filename string = profile.firstName + "_" + profile.lastName + "_" + "WEP.docx"

		// Make sure ./weps folder exists
		if _, err := os.Stat("./weps"); os.IsNotExist(err) {
			fmt.Println("./weps does not exist. Mkdir will create ./weps")
			os.Mkdir("./weps/", 0777)
		}
		// Specify the correct template
		if profile.AP == "true" {
			templateFile := "./templates/AP.docx"
		} else {
			templateFile := "./templates/Honors.docx"
		}
		replaceAndCreate(templateFile, filename, profile)
	}
}

func replaceAndCreate(templateFile string, filename string, profile Student) {
	// Read from docx file
	fmt.Printf("Using %s to create WEP for %s...", templateFile, profile.name)
	t, err := docx.ReadDocxFile(templateFile)
	if err != nil {
		panic(err)
	}
	newDocx := t.Editable()
	// Replace parts of template
	newDocx.Replace("$NAME$", profile.name, -1)
	newDocx.Replace("$GRADE$", profile.grade, -1)
	newDocx.Replace("$TARGET$", profile.target, -1)
	newDocx.Replace("$START$", profile.start, -1)
	newDocx.Replace("$END$", profile.end, -1)
	newDocx.Replace("$REVIEW$", profile.review, -1)
	newDocx.Replace("$ACADEMIC$", profile.academic, -1)
	newDocx.WriteToFile("./weps/" + filename)

	r.Close()
	fmt.Printf(" Done.\n")
}
