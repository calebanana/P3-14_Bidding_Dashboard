package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
)

//////////////////////////////////
//                              //
//          TEMP STUFF          //
//                              //
//////////////////////////////////



const class_data_json = `
{"SemesterStartDate":"16-01-2022","SemesterModules":[{"ModuleCode":"ADB","ModuleName":"Advanced Databases","ModuleClasses":[{"ClassCode":"ADB01","Schedule":"17-01-2022 10:00:00","Tutor":"T002","Capacity":5}]},{"ModuleCode":"DL","ModuleName":"Deep Learning","ModuleClasses":[{"ClassCode":"DL01","Schedule":"18-01-2022 10:00:00","Tutor":"T003","Capacity":11},{"ClassCode":"DL02","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12}]},{"ModuleCode":"DSA","ModuleName":"Data Structures and Algorithms","ModuleClasses":[{"ClassCode":"DSA01","Schedule":"18-01-2022 10:00:00","Tutor":"T003","Capacity":11},{"ClassCode":"DSA02","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12}, {"ClassCode":"DSA03","Schedule":"18-01-2022 10:00:00","Tutor":"T004","Capacity":12},{"ClassCode":"DSA04","Schedule":"19-01-2022 10:00:00","Tutor":"T007","Capacity":5}]}]}
`

//////////////////////////////////
//                              //
//        TEMP STUFF END        //
//                              //
//////////////////////////////////

func GetQueryParams(r *http.Request) (string, string, string){
	params := mux.Vars(r)
	query := r.URL.Query()

	querySemStart := params["semStart"]
	queryClassCode := query.Get("classCode")
	queryStudentId := query.Get("studentId")
	
	return querySemStart, queryClassCode, queryStudentId
}

func bid(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		querySemStart, queryClassCode, queryStudentId := GetQueryParams(r)
		queryFilter := r.URL.Query().Get("filtered")

		var retrievedSem Semester
		if (queryClassCode == "" && queryStudentId == ""){
			// Get all bids of all classes
			// GET http://localhost:8221/api/v1/bid/:semStartDate

			fmt.Println("GET ALL")
			retrievedSem = GetAllBids(querySemStart, "", "", queryFilter)
			
		} else if (queryClassCode != "" && queryStudentId == "") {
			// Get all bids of class
			// GET http://localhost:8221/api/v1/bid/:semStartDate?classCode=...

			fmt.Println("GET CLASS")
			retrievedSem = GetAllBids(querySemStart, queryClassCode, "", queryFilter)
		} else if (queryClassCode == "" && queryStudentId != ""){
			// Get all bids of student
			// GET http://localhost:8221/api/v1/bid/:semStartDate?studentId=...

			fmt.Println("GET STUDENT")
			retrievedSem = GetAllBids(querySemStart, "", queryStudentId, queryFilter)
		} else if (queryClassCode != "" && queryStudentId != ""){
			// Get bid of student and class
			// GET http://localhost:8221/api/v1/bid/:semStartDate?classCode=...&studentId=...

			fmt.Println("GET CLASS AND STUDENT")
			retrievedSem = GetAllBids(querySemStart, queryClassCode, queryStudentId, queryFilter)
		} else {
			// Get all filtered bids of student
			// GET http://localhost:8221/api/v1/bid/:semStartDate?studentId=...&filtered=true

			fmt.Println("GET STUDENT FILTERED")
			retrievedSem = GetAllBids(querySemStart, "", queryStudentId, queryFilter)
		}

		json.NewEncoder(w).Encode(retrievedSem)
	} else if r.Method == "DELETE" {
		querySemStart, queryClassCode, queryStudentId := GetQueryParams(r)

		// Delete Bid for Class
		// DELETE localhost:8221/api/v1/bid/:semStartDate?classCode=...&studentId=...
		fmt.Println("DELETE")
		DeleteBid(querySemStart, queryClassCode, queryStudentId)
	}

	if r.Header.Get("Content-type") == "application/json" {
        if r.Method == "POST" {
			querySemStart, queryClassCode, queryStudentId := GetQueryParams(r)

			if (queryClassCode == "" && queryStudentId == ""){
				// Add Semester, Modules and Empty Classes
				// POST http://localhost:8221/api/v1/bid/:semStartDate

				// reqBody, err := ioutil.ReadAll(r.Body)
				// var newSemester Semester
				// if err == nil {
				// 	json.Unmarshal(reqBody, &newSemester)
				// 	AddNewSemester(newSemester)
				// }

				jsonFile, _ := os.Open("sampleClasses.json")
				byteValue, _ := ioutil.ReadAll(jsonFile)

				var newSemester Semester
				json.Unmarshal(byteValue, &newSemester)
				AddNewSemester(newSemester)

			} else {
				// Add Bid for Class
				// POST localhost:8221/api/v1/bid/:semStartDate?classCode=...&studentId=...

				reqBody, err := ioutil.ReadAll(r.Body)
				var newBid Bid
				if err == nil {
					json.Unmarshal(reqBody, &newBid)
					AddNewBid(querySemStart, queryClassCode, queryStudentId, newBid.BidAmt)
				}
			}
		} else if r.Method == "PUT" {
			querySemStart, queryClassCode, queryStudentId := GetQueryParams(r)

			// Edit Bid for Class
			// POST localhost:8221/api/v1/bid/:semStartDate?classCode=...&studentId=...

			reqBody, err := ioutil.ReadAll(r.Body)
			var editBid Bid
			if err == nil {
				json.Unmarshal(reqBody, &editBid)
				EditBid(querySemStart, queryClassCode, queryStudentId, editBid.BidAmt)
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/bid/{semStart}", bid).Methods("GET", "PUT", "POST", "DELETE")
	fmt.Println("Listening on port 8221")
	http.ListenAndServe(":8221", router)
}
