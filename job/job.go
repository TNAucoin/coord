package job

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/riverqueue/river"
)

type SortArgs struct {
	// Strings is a slice of strings to sort.
	Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }

type SortWorker struct {
	river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
	// inputImage := fmt.Sprintf("%s/%s", job.Args.WorkDir, job.Args.ImageFileName)
	// outputImagePath := fmt.Sprintf("%s/%s", job.Args.WorkDir, job.Args.OuputPath)
	//
	// cmd := exec.Command("convert", inputImage, "-negate", outputImagePath)
	//
	// cmd.Stdout = os.Stdout
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// }
	// return nil
	sort.Strings(job.Args.Strings)
	fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
	return nil
}

// Inverse Color Job
type InverseArgs struct {
	JobDir     string `json:"job_dir"`
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`
}

func (InverseArgs) Kind() string { return "inverse" }

type InverseWorker struct {
	river.WorkerDefaults[InverseArgs]
}

func (w *InverseWorker) Work(ctx context.Context, job *river.Job[InverseArgs]) error {

	inputImage := fmt.Sprintf("%s/%s", job.Args.JobDir, job.Args.InputPath)
	outputImagePath := fmt.Sprintf("%s/%s", job.Args.JobDir, job.Args.OutputPath)

	cmd := exec.Command("convert", inputImage, "-negate", outputImagePath)

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return nil
}

// Create Job Dir
type CreateWorkspaceArgs struct {
	JobID         string `json:"job_id"`
	WorkspacePath string `json:"workspace_path"`
}

func (CreateWorkspaceArgs) Kind() string { return "create_workspace" }

type CreateWorkspaceWorker struct {
	river.WorkerDefaults[CreateWorkspaceArgs]
}

func (w *CreateWorkspaceWorker) Work(ctx context.Context, job *river.Job[CreateWorkspaceArgs]) error {
	_, err := os.Stat(job.Args.WorkspacePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("workspace path is not a valid dir: %v", err)
	}

	jobDir := fmt.Sprintf("%s/%s", job.Args.WorkspacePath, job.Args.JobID)
	if err := os.Mkdir(jobDir, 0755); err != nil {
		return fmt.Errorf("failed to create job dir: %v", err)
	}

	return nil
}
