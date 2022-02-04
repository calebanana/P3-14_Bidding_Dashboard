package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

const classAPI string = "http://10.31.11.11:8041/api/v1/classes/"
const biddingAPI string = "http://10.31.11.11:8221/api/v1/bid/"
const walletAPI string = "http://10.31.11.11:8052/"

var student string
var etiTokens int
var totalStudentBids int

var currentDate = time.Now()
var daysUntilMon = (1 - int(currentDate.Weekday()) + 7) % 7
var semStartDate = currentDate.AddDate(0, 0, daysUntilMon).Format("02-01-2006")
var nextMon = currentDate.AddDate(0, 0, daysUntilMon).Format("02 Jan 2006")

type Info_Class struct {
    ClassCode string
    Schedule string
    Tutor    string
    Capacity int32
    Students []string
}

type Info_Module struct {
    ModuleCode string
	ModuleName string
    ModuleClasses []Info_Class
}

type Info_Semester struct {
    SemesterStartDate string
    SemesterModules []Info_Module
}

type Bid struct {
	StudentID string `bson: "studentID"`
	BidAmt    int  `bson: "bidAmt"`
	BidStatus string `bson: "bidStatus"`
}

type Class struct {
	ClassCode string `bson: "classCode"`
	ClassBids []Bid  `bson: "classBids"`
}

type Module struct {
	ModuleCode    string        `bson: "moduleCode"`
	ModuleName    string        `bson: "moduleName"`
	ModuleClasses []Class       `bson: "moduleClasses"`
}

type Semester struct {
	SemesterStartDate string   
	SemesterModules   []Module
}

type WalletInfo struct {
    WalletID     string `json:"wid"`
    TickerSymbol string `json:"tks"`
    TokenAmount  int    `json:"ta"`
    StudentID    string `json:"stuid"`
}

//////////////////////////////////
//                              //
//          TEMP STUFF          //
//                              //
//////////////////////////////////

func tempLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("web/tempLogin.html"))
		tmpl.Execute(w, nil)
	} else {
		student = r.FormValue("studentid")
		http.Redirect(w, r, "/biddingDashboard/" + student, http.StatusFound)
	}
}

//////////////////////////////////
//                              //
//        TEMP STUFF END        //
//                              //
//////////////////////////////////

func biddingDashboard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	filtered := query.Get("filtered")
	student = params["studentId"] 

	// temp --------------------------------------------------------------------------------------
	if student == "S001" {
		etiTokens = 100
	} else if student == "S002" {
		etiTokens = 150
	} else {
		etiTokens = 50
	}

	jsonFile, _ := os.Open("../BiddingAPI/sampleClasses.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var infoSem Info_Semester
	json.Unmarshal(byteValue, &infoSem)
	// temp end -----------------------------------------------------------------------------------

	// get number of eti tokens
	fmt.Println("GET FROM WALLET ", walletAPI + student + "/Wallet/ETI")
	// walletAPIresponse, _ := http.Get(walletAPI + student + "/Wallet/ETI")
	// walletData, _ := ioutil.ReadAll(walletAPIresponse.Body)

	// var walletInfo WalletInfo
	// _ = json.Unmarshal([]byte(walletData), &walletInfo)

	// walletAPIresponse.Body.Close()

	// etiTokens = walletInfo.TokenAmount

	// get all classes from class api
	fmt.Println("GET FROM CLASS ", classAPI + semStartDate)
	// classAPIresponse, _ := http.Get(classAPI + semStartDate)
	// infoSemData, _ := ioutil.ReadAll(classAPIresponse.Body)

	// // var infoSem Info_Semester
	// _ = json.Unmarshal([]byte(infoSemData), &infoSem)

	// classAPIresponse.Body.Close()	

	// get all bids for student
	biddingAPIresponse, _ := http.Get(biddingAPI + semStartDate + "?studentId=" + student)
	semData, _ := ioutil.ReadAll(biddingAPIresponse.Body)
	var sem Semester
	_ = json.Unmarshal([]byte(semData), &sem)
	biddingAPIresponse.Body.Close()

	var displayInfoSem Info_Semester
	var displaySem Semester
	if filtered == "true"{
		var displayInfoMods []Info_Module
		for modIndex, mod := range sem.SemesterModules{
			var infoModClasses []Info_Class
			for clsIndex, cls := range mod.ModuleClasses{
				if len(cls.ClassBids) != 0{
					var displayInfoCls = infoSem.SemesterModules[modIndex].ModuleClasses[clsIndex]
					fmt.Println(displayInfoCls)
					infoModClasses = append(infoModClasses, displayInfoCls)
				}
			}
			var displayInfoMod = Info_Module{
				ModuleCode: mod.ModuleCode,
				ModuleName: mod.ModuleName,
				ModuleClasses: infoModClasses,
			}

			if len(displayInfoMod.ModuleClasses) > 0{
				displayInfoMods = append(displayInfoMods, displayInfoMod)
			}
		}
		displayInfoSem = Info_Semester{
			SemesterStartDate: semStartDate,
			SemesterModules: displayInfoMods,
		}

		biddingAPIresponse, _ := http.Get(biddingAPI + semStartDate + "?studentId=" + student + "&filtered=true")
		displaySemData, _ := ioutil.ReadAll(biddingAPIresponse.Body)
		_ = json.Unmarshal([]byte(displaySemData), &displaySem)
		biddingAPIresponse.Body.Close()
		fmt.Println(displayInfoSem, "\n\n")
		fmt.Println(displaySem)
	} else {
		displayInfoSem = infoSem
		displaySem = sem
	}

	// calculate total student bids
	totalStudentBids = 0
	for _, mod := range sem.SemesterModules{
		for _, cls := range mod.ModuleClasses{
			if len(cls.ClassBids) > 0{
				totalStudentBids += cls.ClassBids[0].BidAmt
			}
		}
	}
	etiTokens -= totalStudentBids

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("web/biddingDashboard.html"))	

		data := map[string]interface{}{
			"StudentID": student,
			"ETITokens": etiTokens,
			"NextMon": nextMon,
			"SemInfo": displayInfoSem.SemesterModules,
			"SemBids": displaySem.SemesterModules,
			"Filtered": filtered,
		}

		tmpl.Execute(w, data)
	} else if r.Method == "POST"{
		moduleSearch := r.FormValue("moduleSearch")

		if moduleSearch != "all"{
			var displayInfoMods []Info_Module
			var displayMods []Module
			for modIndex, mod := range sem.SemesterModules{
				if mod.ModuleCode == moduleSearch{
					infoMod := infoSem.SemesterModules[modIndex]
					displayInfoMods = append(displayInfoMods, infoMod)

					mod := sem.SemesterModules[modIndex]
					displayMods = append(displayMods, mod)
				}
			}
			displayInfoSem = Info_Semester{
				SemesterStartDate: semStartDate,
				SemesterModules: displayInfoMods,
			}

			displaySem = Semester{
				SemesterStartDate: semStartDate,
				SemesterModules: displayMods,
			}
		} else {
			displayInfoSem = infoSem
			displaySem = sem
		}

		tmpl := template.Must(template.ParseFiles("web/biddingDashboard.html"))
		data := map[string]interface{}{
			"StudentID": student,
			"ETITokens": etiTokens,
			"NextMon": nextMon,
			"SemInfo": displayInfoSem.SemesterModules,
			"SemBids": displaySem.SemesterModules,
			"Filtered": filtered,
		}

		tmpl.Execute(w, data)
	}
}

func editBid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	classCode := params["classCode"]

	biddingAPIresponse, _ := http.Get(biddingAPI + semStartDate + "?classCode=" + classCode + "&studentId=" + student)
	semData, _ := ioutil.ReadAll(biddingAPIresponse.Body)

	var sem Semester
	_ = json.Unmarshal([]byte(semData), &sem)

	biddingAPIresponse.Body.Close()

	var studentBid Bid 
	if len(sem.SemesterModules[0].ModuleClasses[0].ClassBids) > 0{
		studentBid = sem.SemesterModules[0].ModuleClasses[0].ClassBids[0]
	} else {
		studentBid = Bid{
			StudentID: student,
			BidAmt: 0,
			BidStatus: "Pending",
		}
	}

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("web/editBid.html"))

		data := map[string]interface{}{
			"StudentID": student,
			"ETITokens": etiTokens,
			"ClassCode": classCode,
			"StudentBid": studentBid,
		}

		tmpl.Execute(w, data)
	} else if r.Method == "POST" {

		bidAmt, _ := strconv.Atoi(r.FormValue("bidAmt"))

		var editBid = Bid{
			StudentID: r.FormValue("studentId"),
			BidAmt: bidAmt,
			BidStatus: "Pending",
		}

		if studentBid.BidAmt == 0 && editBid.BidAmt > 0 {
			// Add bid
			editBid_json, _ := json.Marshal(editBid)

			response, _ := http.Post(biddingAPI + semStartDate + "?classCode=" + classCode + "&studentId=" + editBid.StudentID, "application/json", bytes.NewBuffer(editBid_json))
			response.Body.Close()

			http.Redirect(w, r, "/biddingDashboard/" + student, http.StatusFound)

		} else if studentBid.BidAmt > 0 && editBid.BidAmt == 0 {
			http.Redirect(w, r, "/deleteBid/" + classCode + "/" + studentBid.StudentID, http.StatusFound)
		} else {
			editBid_json, _ := json.Marshal(editBid)

			request, _ := http.NewRequest(http.MethodPut,
				biddingAPI + semStartDate + "?classCode=" + classCode + "&studentId=" + editBid.StudentID,
				bytes.NewBuffer(editBid_json))
	
			request.Header.Set("Content-Type", "application/json")
	
			client := &http.Client{}
			response, _ := client.Do(request)

			response.Body.Close()

			http.Redirect(w, r, "/biddingDashboard/" + student, http.StatusFound)
		}
	}
}

func viewAll(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tmpl := template.Must(template.ParseFiles("web/viewAll.html"))

	classCode := params["classCode"]

	biddingAPIresponse, _ := http.Get(biddingAPI + semStartDate + "?classCode=" + classCode)
	semData, _ := ioutil.ReadAll(biddingAPIresponse.Body)

	var sem Semester
	_ = json.Unmarshal([]byte(semData), &sem)

	biddingAPIresponse.Body.Close()

	var retrievedClassCode = sem.SemesterModules[0].ModuleClasses[0].ClassCode
	var retrievedClassBids = sem.SemesterModules[0].ModuleClasses[0].ClassBids

	sort.Slice(retrievedClassBids, func(i, j int) bool {
		return retrievedClassBids[i].BidAmt > retrievedClassBids[j].BidAmt
	})

	fmt.Println(retrievedClassBids)

	data := map[string]interface{}{
		"StudentID": student,
		"ClassCode": retrievedClassCode,
		"ClassBids": retrievedClassBids,
	}

	tmpl.Execute(w, data)
}

func deleteBid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	classCode := params["classCode"]
	studentId := params["studentId"]

	//fmt.Println(biddingAPI + semStartDate + "?classCode=" + classCode + "&studentId=" + studentId)
	request, _ := http.NewRequest(http.MethodDelete,
		biddingAPI + semStartDate + "?classCode=" + classCode + "&studentId=" + studentId,
		nil)

	client := &http.Client{}
	response, _ := client.Do(request)

	response.Body.Close()
	
	http.Redirect(w, r, "/biddingDashboard/" + student, http.StatusFound)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", tempLogin)

	router.HandleFunc("/biddingDashboard/{studentId}", biddingDashboard)
	router.HandleFunc("/editBid/{classCode}", editBid)
	router.HandleFunc("/viewAll/{classCode}", viewAll)
	router.HandleFunc("/deleteBid/{classCode}/{studentId}", deleteBid)

	fmt.Println("Listening on port 8220")
	http.ListenAndServe(":8220", router)
}
