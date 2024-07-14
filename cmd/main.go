package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/tnaucoin/coord/config"
	"github.com/tnaucoin/coord/job"
)

type App struct {
	mux           *mux.Router
	dbPool        *pgxpool.Pool
	riverClient   *river.Client[pgx.Tx]
	JobSubmission map[string]JobItem
	Port          string
}

type JobItem struct {
	ID          string `json:"id"`
	Steps       []Step `json:"steps"`
	CurrentStep int    `json:"currentStep"`
}

type Step struct {
	Args map[string]interface{} `json:"args"`
	Type string                 `json:"type"`
}

type JobSubmission struct {
	Steps []Step `json:"steps"`
}

type JobSubmissionResponse struct {
	ID string `json:"id"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables %v", err)
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, config.GetDatabaseConnectionURL())
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}
	defer dbPool.Close()

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{})
	if err != nil {
		log.Fatalf("Failed to create river client %v", err)
	}

	JobSubmission := make(map[string]JobItem)
	r := mux.NewRouter()
	applicationPort := fmt.Sprintf(":%d", config.GetApplicationPort())
	App := App{Port: applicationPort, mux: r, JobSubmission: JobSubmission, riverClient: riverClient, dbPool: dbPool}
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

// queueJobTask resposible for processing a request job step, and queueing it in river for processing
func (a *App) queueJobTask(step Step) {
	ctx := context.Background()
	tx, err := a.dbPool.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction %v", err)
	}
	defer tx.Rollback(ctx)
	jobArgs := a.matchTypeToKind(step)
	_, err = a.riverClient.InsertTx(ctx, tx, jobArgs, nil)
	if err != nil {
		log.Fatalf("Failed to insert task %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("Failed to commit transaction %v", err)
	}
}

// / matchTypeToKind converts the request job step, into a river queue job type to be processed
func (a *App) matchTypeToKind(step Step) river.JobArgs {
	switch step.Type {
	case job.SortArgs{}.Kind():
		// create a new job task from the step
		jobArgs := job.SortArgs{
			Strings: step.Args["strings"].([]string),
		}
		return jobArgs
	default:
		log.Fatalf("Unknown job type %s", step.Type)
		return nil
	}
}
