package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	"encoding/json"


	"github.com/gorilla/mux"
)

const classAPI string = "url"

var student string
var etiTokens int

var currentDate = time.Now()
var daysUntilMon = (1 - int(currentDate.Weekday()) + 7) % 7
var nextMon = currentDate.AddDate(0, 0, daysUntilMon).Format("02 Jan 2006")

//////////////////////////////////
//                              //
//          TEMP STUFF          //
//                              //
//////////////////////////////////
type Class struct {
    ClassCode string
    Schedule string
    Tutor    string
    Capacity int32
    Students []string
}

type Module struct {
    ModuleCode string
	ModuleName string
    ModuleClasses []Class
}

type Semester struct {
    SemesterStartDate string
    SemesterModules []Module
}

const class_data_json = `
{"SemesterStartDate":"16-01-2022","SemesterModules":[{"ModuleCode":"ADB","ModuleName":"Advanced Databases","ModuleClasses":[{"ClassCode":"ADB01","Schedule":"17-01-2022 10:00:00","Tutor":"T002","Capacity":5,"Students":["S004","S005","S006"]}]},{"ModuleCode":"DL","ModuleName":"Deep Learning","ModuleClasses":[{"ClassCode":"DL01","Schedule":"18-01-2022 10:00:00","Tutor":"T003","Capacity":11,"Students":["S007","S002","S001"]},{"ClassCode":"DL02","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12,"Students":["S010","S012","S003"]}]},{"ModuleCode":"DSA","ModuleName":"Data Structures and Algorithms","ModuleClasses":[{"ClassCode":"DSA01","Schedule":"18-01-2022 10:00:00","Tutor":"T003","Capacity":11,"Students":["S007","S002","S001"]},{"ClassCode":"DSA02","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12,"Students":["S010","S012","S003"]}, {"ClassCode":"DSA03","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12,"Students":["S010","S012","S003"]}]}]}
`

func tempLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("web/tempLogin.html"))
		tmpl.Execute(w, nil)
	} else {
		student = r.FormValue("studentid")
		if student == "S001" {
			etiTokens = 100
		} else if student == "S002" {
			etiTokens = 150
		}
		http.Redirect(w, r, "/biddingDashboard", http.StatusFound)
	}
}

//////////////////////////////////
//                              //
//        TEMP STUFF END        //
//                              //
//////////////////////////////////

func biddingDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/biddingDashboard.html"))

	var sem Semester
	_ = json.Unmarshal([]byte(class_data_json), &sem)

	for _, i := range sem.SemesterModules {
		fmt.Println(i, "\n")
	}

	data := map[string]interface{}{
		"Student_ID": student,
		"ETI_Tokens": etiTokens,
		"NextMon":    nextMon,
		"SemesterModules": sem.SemesterModules,
	}

	tmpl.Execute(w, data)
}

func editBid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//tmpl := template.Must(template.ParseFiles("web/editBid.html"))

	fmt.Println(params)

	// tmpl.Execute(w, data)
}

func viewBid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//tmpl := template.Must(template.ParseFiles("web/editBid.html"))

	fmt.Println(params)

	// tmpl.Execute(w, data)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", tempLogin)

	router.HandleFunc("/biddingDashboard", biddingDashboard)
	router.HandleFunc("/editBid/{moduleCode}/{classNum}", editBid)
	router.HandleFunc("/viewAll/{moduleCode}/{classNum}", viewBid)

	fmt.Println("Listening on port 8220")
	http.ListenAndServe(":8220", router)
}
