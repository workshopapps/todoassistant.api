package reminderEntity

type GetPendingTasks struct {
	TaskId      string `json:"task_id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndTime     string `json:"end_time"`
	DeviceId    string `json:"device_id"`
	// request for searched task
}
