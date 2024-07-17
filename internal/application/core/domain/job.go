package domain

type Job struct {
	ID          int64     `json:"id"`
	JobSteps    []JobStep `json:"steps"`
	CurrentStep int       `json:"current_step"`
	Status      string    `json:"status"`
}

type JobStep struct {
	// Args map[string]interface{} `json:"args"`
	Type string `json:"type"`
}

func NewJob(jobID string, steps []JobStep) Job {
	return Job{
		JobSteps: steps,
		Status:   "pending",
	}
}
