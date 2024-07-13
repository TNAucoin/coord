package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	Port          string
	mux           *mux.Router
	JobSubmission map[string]JobItem
}

type JobItem struct {
	ID          string `json:"id"`
	Steps       []Step `json:"steps"`
	CurrentStep int    `json:"currentStep"`
}

type Step struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

type JobSubmission struct {
	Steps []Step `json:"steps"`
}

type JobSubmissionResponse struct {
	ID string `json:"id"`
}

func main() {
	JobSubmission := make(map[string]JobItem)
	r := mux.NewRouter()

	App := App{Port: ":8080", mux: r, JobSubmission: JobSubmission}
	App.RegisterRoutes()

	srv := &http.Server{
		Addr:    App.Port,
		Handler: App.mux,
	}

	srv.ListenAndServe()
}

func (a *App) RegisterRoutes() {
	a.mux.HandleFunc("/job", a.createJobHandler()).Methods("POST")
	a.mux.HandleFunc("/job/{id}", a.getJobHandler()).Methods("GET")
}

func (a *App) createJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job JobSubmission
		err := json.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j := a.createJobSubmissionItem(job.Steps)
		a.JobSubmission[j.ID] = *j
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JobSubmissionResponse{ID: j.ID})
	}
}

func (a *App) getJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if job, ok := a.JobSubmission[id]; ok {
			json.NewEncoder(w).Encode(job)
		} else {
			http.Error(w, "Job not found", http.StatusNotFound)
		}
	}
}

func (a *App) createJobSubmissionItem(steps []Step) *JobItem {
	id := uuid.New().String()
	return &JobItem{ID: id, Steps: steps, CurrentStep: 0}
}
