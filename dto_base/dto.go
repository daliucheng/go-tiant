package dto_base

type TaskCommonReq struct {
	SessionId string `json:"session_id,optional"`
	TaskType  string `json:"task_type"`
	Input     string `json:"input,optional"`
	TimeoutMs int    `json:"timeout_ms,optional,default=10000"`
}

type TaskCommonResp struct {
	SessionId string `json:"session_id,optional"`
	TaskId    string `json:"task_id,optional"`
	TaskType  string `json:"task_type,optional"`
	Input     string `json:"input,optional"`
	Output    string `json:"output,optional"`
	Status    string `json:"status,optional"`
}
