package services

import (
	"encoding/json"
	"fmt"
	"os"
)

type Tasks struct {
	TaskID string `json:"task_id"`
}

type TaskResponse struct {
	Status string `json:"status"`
}

func (task Tasks) MonitorTask(token string, taskID string, request Request) (status string) {
	task_status := "NOT_STARTED"
	exitLoop := false

	for !exitLoop {
		url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + TASKS + "/" + taskID
		body, _ := processRequest(token, url, "GET", nil)

		taskResponse := TaskResponse{}
		err := json.Unmarshal(body, &taskResponse)
		if err != nil {
			fmt.Println("Failed to parse the response body.\n[ERROR] -", err)
			os.Exit(1)
		}

		task_status = taskResponse.Status

		fmt.Println("Task status is: ", task_status)

		if task_status != "IN_PROGRESS" && task_status != "NOT_STARTED" {
			exitLoop = true
		}

	}

	return task_status
}
