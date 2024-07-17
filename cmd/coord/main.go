package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		j, err := a.createJobSubmissionItem(job.Steps)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func (a *App) createJobSubmissionItem(steps []Step) (*JobItem, error) {
	id := uuid.New().String()
	jobItem := &JobItem{ID: id, Steps: steps, CurrentStep: 0}
	if err := a.setupWorkspace(jobItem); err != nil {
		return nil, err
	}
	// TODO: we will probably want some validation checking here prior to starting queue
	if err := a.queueJobTask(jobItem); err != nil {
		log.Fatal("failed to queue jobItem step")
		return nil, err
	}

	return jobItem, nil
}

func (a *App) setupWorkspace(jobItem *JobItem) error {
	var workspaceRootDir = "/workspace"
	_, err := os.Stat(workspaceRootDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("workspace dir not found: %v", err)
	}
	if err := os.Mkdir(fmt.Sprintf("%s/%s", workspaceRootDir, jobItem.ID), 0755); err != nil {
		return fmt.Errorf("failed to create job workspace folder: %v", err)
	}
	return nil
}

// queueJobTask resposible for processing a request job step, and queueing it in river for processing
func (a *App) queueJobTask(jobItem *JobItem) error {
	ctx := context.Background()
	tx, err := a.dbPool.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction %v", err)
		return err
	}
	defer tx.Rollback(ctx)
	// TODO: do this in a better way, for now we just pull the first step
	jobArgs, err := a.convertReqArgsToJobArgs(jobItem.Steps[0])
	if err != nil {
		fmt.Printf("error matching kind to type: %v", err)
		return err
	}
	_, err = a.riverClient.InsertTx(ctx, tx, jobArgs, nil)
	if err != nil {
		log.Fatalf("Failed to insert task %v", err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("Failed to commit transaction %v", err)
		return err
	}
	return nil
}

// / matchTypeToKind converts the request job step, into a river queue job type to be processed
func (a *App) convertReqArgsToJobArgs(step Step) (river.JobArgs, error) {
	argBytes, err := json.Marshal(step.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal step args: %v", err)
	}
	switch step.Type {

	case job.InverseArgs{}.Kind():
		log.Printf("got inverse job step..")
		var inverseArgs job.InverseArgs

		if err := json.Unmarshal(argBytes, &inverseArgs); err != nil {
			log.Printf("error unmarshalling args")
			return nil, fmt.Errorf("error unmarshalling args for inverse: %v", err)
		}
		log.Printf("got inverse job: %v", inverseArgs)
		return inverseArgs, nil

	case job.SortArgs{}.Kind():
		var sortArgs job.SortArgs

		if err := json.Unmarshal(argBytes, &sortArgs); err != nil {
			return nil, fmt.Errorf("error unmarshalling args for sort: %v", err)
		}
		fmt.Printf("got sort job: %v", sortArgs)
		return sortArgs, nil

	default:
		return nil, fmt.Errorf("invalid job type")
	}
}
