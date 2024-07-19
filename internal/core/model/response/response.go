package response

type Response struct {
	Data    any    `json:"data"`
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SubmitJobResponse struct {
	JobID string `json:"job_id"`
}
