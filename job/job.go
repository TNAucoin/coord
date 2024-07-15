package job

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("got inverse job args: %v", job.Args)
	fmt.Printf("got args: %v", job.Args)
	return nil
}
