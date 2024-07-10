package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type JobItem struct {
	ID          string `json:"id"`
	Steps       []Step `json:"steps"`
	CurrentStep int    `json:"currentStep"`
}

type Step struct {
	Type   string `json:"type"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type JobSubmission struct {
	Steps []Step `json:"steps"`
}

type JobSubmissionResponse struct {
	ID string `json:"id"`
}

func main() {
	jobArr := []JobItem{}

	r := mux.NewRouter()
	r.HandleFunc("/job", createJobHandler(&jobArr)).Methods("POST")
	r.HandleFunc("/job/{id}", getJobHandler(&jobArr)).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	srv.ListenAndServe()
}

func createJobHandler(jobArray *[]JobItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job JobSubmission
		err := json.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := uuid.New().String()
		newJob := JobItem{ID: id, Steps: job.Steps, CurrentStep: 0}
		*jobArray = append(*jobArray, newJob)
		fmt.Println("added job", newJob)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JobSubmissionResponse{ID: id})
	}
}

func getJobHandler(jobArray *[]JobItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		for _, job := range *jobArray {
			if job.ID == id {
				json.NewEncoder(w).Encode(job)
				return
			}
		}
		fmt.Printf("jobArr: %v\n", jobArray)
		http.Error(w, "Job not found", http.StatusNotFound)
	}
}
