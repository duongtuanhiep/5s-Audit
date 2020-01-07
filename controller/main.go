package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/duongtuanhiep/fliqaudit/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Global varriable for error and database connection
var db *sql.DB
var err error

func main() {
	//Oppening database connection
	// db, err = sql.Open("mysql", "fliqtesting:Deoconao1@tcp(localhost:3306)/fliqaudit?parseTime=true&group_concat_max_len=8192")
	db, err = sql.Open("mysql", "root:ditmetuoiloz@tcp(database-2.cluster-cyu24tc3dqek.eu-north-1.rds.amazonaws.com:3306)/fliqaudit?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	// Closing database connection
	defer db.Close()

	//Creating router
	router := mux.NewRouter()

	//Endpoint mappping
	ver := "/fliq/v3"
	router.HandleFunc(ver+"/stat/{id}", getStat).Methods("GET")
	router.HandleFunc(ver+"/stat", getStats).Methods("GET")
	router.HandleFunc(ver+"/audits", getAudits).Methods("GET")
	router.HandleFunc(ver+"/shortaudits", getShortAudits).Methods("GET")
	router.HandleFunc(ver+"/audits/{id}", getAudit).Methods("GET")
	router.HandleFunc(ver+"/activequestions", getActiveQuestion).Methods("GET")
	router.HandleFunc(ver+"/audits", createAudit).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))

	// http.ListenAndServe(":8080", router)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var auditStat model.AuditsStat

	auditStatResult, err := db.Query("select auditphase, round(avg(answer),2) as averageScore from checkitemfull group by auditphase")
	if err != nil {
		panic(err.Error())
	}
	defer auditStatResult.Close()
	for auditStatResult.Next() {
		var phaseAverage model.PhaseAverage
		err := auditStatResult.Scan(&phaseAverage.AuditPhase, &phaseAverage.Score)
		if err != nil {
			panic(err.Error())
		}
		auditStat.PhaseScores = append(auditStat.PhaseScores, phaseAverage)
	}
	json.NewEncoder(w).Encode(auditStat)
}

func getStat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var auditStat model.AuditStat

	params := mux.Vars(r)
	auditStat.Key, err = strconv.Atoi(params["id"])
	if err != nil {
		panic(err.Error())
	}
	auditStatResult, err := db.Query("select auditphase, round(avg(answer),2) as averageScore from checkitemfull where auditKey = ? group by auditphase", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer auditStatResult.Close()
	for auditStatResult.Next() {
		var phaseAverage model.PhaseAverage
		err := auditStatResult.Scan(&phaseAverage.AuditPhase, &phaseAverage.Score)
		if err != nil {
			panic(err.Error())
		}
		auditStat.PhaseScores = append(auditStat.PhaseScores, phaseAverage)
	}
	json.NewEncoder(w).Encode(auditStat)
}

func getAudits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Initialize a slice of audits
	var audits []model.Audit

	//Select all audit from audit table
	auditResult, err := db.Query("SELECT * from audit ")
	if err != nil {
		panic(err.Error())
	}
	defer auditResult.Close()

	// Cycling through the array to assemble the checklist array from 2 table checkitem and question
	//Look into database diagram for further understanding this process
	for auditResult.Next() {
		var audit model.Audit
		var checkitems []model.CheckItem

		//fetching audit data onto our current Audit object
		err := auditResult.Scan(&audit.Key, &audit.DepartmentName, &audit.CompanyName, &audit.AuditorName, &audit.AssessedDate)
		if err != nil {
			panic(err.Error())
		}

		//Querying checkitem
		checkItemResult, err := db.Query("select checkItemID,questionID,checkItem,checkItemDescription,answer,auditPhase,status from checkitemfull where auditKey = ?", audit.Key)
		if err != nil {
			panic(err.Error())
		}
		defer checkItemResult.Close()
		for checkItemResult.Next() {
			var checkitem model.CheckItem
			err := checkItemResult.Scan(&checkitem.CheckItemID, &checkitem.QuestionID, &checkitem.CheckItem, &checkitem.CheckItemDescription, &checkitem.CheckItemAnswer, &checkitem.AuditPhase, &checkitem.Status)
			if err != nil {
				panic(err.Error())
			}
			checkitems = append(checkitems, checkitem)
		}
		//Jsoning
		// json.Unmarshal([]byte(holder), &checkitems)
		audit.CheckItems = checkitems
		audits = append(audits, audit)
	}
	json.NewEncoder(w).Encode(audits)
}

func getShortAudits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// HANDLING PAPGINATION
	queryMap := r.URL.Query()

	//Default value for page number and page size
	pageNumber := 0
	pageSize := 0
	//TESST
	pageNumber, err := strconv.Atoi(queryMap.Get("pageNumber"))
	if err != nil {
		pageNumber = 0
	}
	pageSize, err = strconv.Atoi(queryMap.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	//Initalize an array of short audits
	var shortAudits []model.ShortAudit

	//Select all the properties needed for short version of audit from audit table : key, auditor name and assessed date
	shortAuditResult, err := db.Query("select audit.key, audit.auditorName, audit.assessedDate from audit order by assessedDate DESC LIMIT ?,?", (pageNumber * pageSize), pageSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err.Error())
	}
	defer shortAuditResult.Close()
	// Looping through the result set to assemble the data to the slice
	for shortAuditResult.Next() {
		var shortAudit model.ShortAudit
		err := shortAuditResult.Scan(&shortAudit.Key, &shortAudit.AuditorName, &shortAudit.AssessedDate)
		if err != nil {
			panic(err.Error())
		}
		averageScoreResult, err := db.Query("select round(avg(checkitem.answer),2) as averageScore from checkitem where checkitem.auditkey = ? group by checkitem.auditKey", shortAudit.Key)
		if err != nil {
			panic(err.Error())
		}
		defer averageScoreResult.Close()
		for averageScoreResult.Next() {
			err := averageScoreResult.Scan(&shortAudit.AverageScore)
			if err != nil {
				panic(err.Error())
			}
		}
		// select round(avg(checkitem.answer),2) as averageScore from checkitem where checkitem.auditkey = ? group by checkitem.auditKey
		shortAudits = append(shortAudits, shortAudit)
	}
	json.NewEncoder(w).Encode(shortAudits)
}

func getAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	auditResult, err := db.Query("SELECT * FROM audit WHERE audit.key = ?", params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(string(err.Error()))
		panic(err.Error())
	}
	defer auditResult.Close()
	// Create dummy outside of the loop since this will only loops one - best practice
	var audit model.Audit
	var checkitems []model.CheckItem
	for auditResult.Next() {
		//fetching audit data onto our current Audit object
		err := auditResult.Scan(&audit.Key, &audit.DepartmentName, &audit.CompanyName, &audit.AuditorName, &audit.AssessedDate)
		if err != nil {
			panic(err.Error())
		}

		//Querying checkitem
		checkItemResult, err := db.Query("select checkItemID,questionID,checkItem,checkItemDescription,answer,auditPhase,status from checkitemfull where auditKey = ?", audit.Key)
		if err != nil {
			panic(err.Error())
		}
		defer checkItemResult.Close()
		for checkItemResult.Next() {
			var checkitem model.CheckItem
			err := checkItemResult.Scan(&checkitem.CheckItemID, &checkitem.QuestionID, &checkitem.CheckItem, &checkitem.CheckItemDescription, &checkitem.CheckItemAnswer, &checkitem.AuditPhase, &checkitem.Status)
			if err != nil {
				panic(err.Error())
			}
			checkitems = append(checkitems, checkitem)
		}
		//Jsoning
		// json.Unmarshal([]byte(holder), &checkitems)
		audit.CheckItems = checkitems
	}
	json.NewEncoder(w).Encode(audit)
}

func getActiveQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var questions []model.Question
	questionResult, err := db.Query("SELECT * FROM question WHERE status =\"ACTIVE\" ")
	if err != nil {
		panic(err.Error())
	}
	defer questionResult.Close()
	for questionResult.Next() {
		var question model.Question
		err := questionResult.Scan(&question.QuestionID, &question.CheckItem, &question.CheckItemDescription, &question.AuditPhase, &question.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err.Error())
		}
		questions = append(questions, question)
	}
	json.NewEncoder(w).Encode(questions)
}

func createAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var auditsend model.AuditSend

	//Read request Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	// Un marshall into object suitable for sending back to database
	err = json.Unmarshal(body, &auditsend)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	//Create time instance
	timeNow := time.Now().Format("2006-01-02 15:04:05")

	res, err := db.Exec("INSERT INTO audit (departmentName,companyName,auditorName,assessedDate) VALUES(?,?,?,?)", auditsend.DepartmentName, auditsend.CompanyName, auditsend.AuditorName, timeNow)
	if err != nil {
		panic(err.Error())
	}

	//new way
	stmt, err := db.Prepare("INSERT INTO checkitem (questionID,auditKey,answer) VALUES(?,?,?)")
	for _, checkitemsend := range auditsend.CheckItemsSend {
		defer stmt.Close()
		holder, err := res.LastInsertId()
		_, err = stmt.Exec(checkitemsend.QuestionID, holder, checkitemsend.CheckItemAnswer)
		if err != nil {
			panic(err.Error())
		}
	}
	stmt.Close()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(auditsend)
	// fmt.Fprintf(w, "Audit Created")
}
