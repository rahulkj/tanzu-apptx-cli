package services

import (
	"encoding/json"
	"log"
	"os"
)

type Tasks struct {
	TaskId string `json:"task_id"`
}

type TaskResponse struct {
	Status string `json:"status"`
}

func (task Tasks) MonitorTask(token string, taskId string, request Request) (status string) {
	task_status := "UNKNOWN"

	for task_status != "SUCCESS" {
		url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + TASKS + "/" + taskId
		body, _ := processRequest(token, url, "GET", nil)

		taskResponse := TaskResponse{}
		err := json.Unmarshal(body, &taskResponse)
		if err != nil {
			log.Println("Failed to parse the response body.\n[ERROR] -", err)
			os.Exit(1)
		}
		task_status = taskResponse.Status
	}

	return task_status
}
