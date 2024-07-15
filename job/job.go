package job

import (
	"context"
	"fmt"
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
