package main

import (
	"context"
	"fmt"
	"log"

	// "math"
	"net/http"
	"personal-project/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/hi", helloworld).Methods("GET")

	route.HandleFunc("/home", home).Methods("GET")

	route.HandleFunc("/project", project).Methods("GET")

	route.HandleFunc("/submit", submit).Methods("POST")

	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/delete/{index}", delete).Methods("GET")

	route.HandleFunc("/edit/{index}", edit).Methods("GET")

	route.HandleFunc("/editButton/{id}", editButton).Methods("POST")

	route.HandleFunc("/detail/{id}", detail).Methods("GET")

	fmt.Println("server is Running")
	http.ListenAndServe("localhost:8000", route)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/home.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// response := map[string]interface{}{
	// 	"Project": projectData,
	// }

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, description, duration FROM db_project")
	//fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(&each.Id, &each.Title, &each.Description, &each.Duration)
		if err != nil {
			fmt.Println(err.Error)
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Project": result,
	}

	// fmt.Println(result)

	tmplt.Execute(w, resData)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/project.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, "")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/contact.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, "")
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/detail.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	var DetailProject = Project{}

	// for i, project := range projectData {
	// 	if i == index {
	// 		DetailProject = Project{
	// 			Title:       project.Title,
	// 			Description: project.Description,
	// 			StartDate:   project.StartDate,
	// 			EndDate:     project.EndDate,
	// 			Duration:    project.Duration,
	// 			Golang:      project.Golang,
	// 			JavaScript:  project.JavaScript,
	// 			React:       project.React,
	// 			Node:        project.Node,
	// 		}
	// 	}
	// }

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, duration FROM db_project WHERE id=$1", id).Scan(
		&DetailProject.Id, &DetailProject.Title, &DetailProject.StartDate, &DetailProject.EndDate, &DetailProject.Description, &DetailProject.Duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	DetailProject.Formatstartdate = DetailProject.StartDate.Format("2 january 2006")
	DetailProject.Formatenddate = DetailProject.EndDate.Format("2 january 2006")

	// DetailProject.Duration = Duration

	data := map[string]interface{}{
		"DetailProject": DetailProject,
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, data)
}

type Project struct {
	Title           string
	StartDate       time.Time
	EndDate         time.Time
	Description     string
	Duration        string
	Golang          string
	JavaScript      string
	React           string
	Node            string
	Id              int
	Formatstartdate string
	Formatenddate   string
}

var projectData = []Project{}

func submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("addTitle")
	startDate := r.PostForm.Get("addStartDate")
	endDate := r.PostForm.Get("addEndDate")
	description := r.PostForm.Get("addDescription")

	layout := "2006-01-02"
	parsingstartdate, _ := time.Parse(layout, startDate)
	parsingenddate, _ := time.Parse(layout, endDate)

	hours := parsingenddate.Sub(parsingstartdate).Hours()
	days := hours / 24
	// weeks := math.Round(days / 7)
	// month := math.Round(days / 30)
	// year := math.Round(days / 365)

	var duration string
	// if year > 0 {
	// 	duration = strconv.FormatFloat(year, 'f', 0, 64) + " years"
	// } else if month > 0 {
	// 	duration = strconv.FormatFloat(month, 'f', 0, 64) + " month"
	// } else if weeks > 0 {
	// 	duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " weeks"
	// } else
	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO db_project(title, start_date, end_date, description, duration) VALUES($1, $2, $3, $4, $5)", title, parsingstartdate, parsingenddate, description, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

	// golang := r.PostForm.Get("addGolang")
	// javaScript := r.PostForm.Get("addJavaScript")
	// react := r.PostForm.Get("addReact")
	// nodejs := r.PostForm.Get("addNode")

	// newProject := Project{
	// 	Title:       title,
	// 	StartDate:   startDate,
	// 	EndDate:     endDate,
	// 	Duration:    duration,
	// 	Description: description,
	// 	Golang:      golang,
	// 	JavaScript:  javaScript,
	// 	React:       react,
	// 	Node:        nodejs,
	// 	Id:          len(projectData),
	// }

	//projectData = append(projectData, newProject)

	// fmt.Println(projectData)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["index"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM db_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/project-edit.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	var DetailProject = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, FROM db_project WHERE id=$1", id).Scan(
		&DetailProject.Id, &DetailProject.Title, &DetailProject.StartDate, &DetailProject.EndDate, &DetailProject.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// for i, project := range projectData {
	// 	if i == index {
	// 		DetailProject = Project{
	// 			Title:       project.Title,
	// 			Description: project.Description,
	// 			StartDate:   project.StartDate,
	// 			EndDate:     project.EndDate,
	// 			Duration:    project.Duration,
	// 			Golang:      project.Golang,
	// 			JavaScript:  project.JavaScript,
	// 			React:       project.React,
	// 			Node:        project.Node,
	// 		}
	// 	}
	// }

	data := map[string]interface{}{
		"ProjectEdit": DetailProject,
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, data)
}

func editButton(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

	title := r.PostForm.Get("addTitle")
	startDate := r.PostForm.Get("addStartDate")
	endDate := r.PostForm.Get("addEndDate")
	description := r.PostForm.Get("addDescription")

	golang := r.PostForm.Get("addGolang")
	javaScript := r.PostForm.Get("addJavaScript")
	react := r.PostForm.Get("addReact")
	nodejs := r.PostForm.Get("addNode")

	layout := "2006-01-02"
	parsingstartdate, _ := time.Parse(layout, startDate)
	parsingenddate, _ := time.Parse(layout, endDate)

	hours := parsingenddate.Sub(parsingstartdate).Hours()
	days := hours / 24
	// weeks := math.Round(days / 7)
	// month := math.Round(days / 30)
	// year := math.Round(days / 365)

	var duration string
	// if year > 0 {
	// 	duration = strconv.FormatFloat(year, 'f', 0, 64) + " years"
	// } else if month > 0 {
	// 	duration = strconv.FormatFloat(month, 'f', 0, 64) + " month"
	// } else if weeks > 0 {
	// 	duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " weeks"
	// } else
	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}

	newProject := Project{
		Title:       title,
		Duration:    duration,
		Description: description,
		Golang:      golang,
		JavaScript:  javaScript,
		React:       react,
		Node:        nodejs,
		Id:          id,
	}
	projectData[id] = newProject

	// fmt.Println(projectData)
}
