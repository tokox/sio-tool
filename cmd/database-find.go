package cmd

import (
	"database/sql"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"strconv"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/fatih/color"
	_ "modernc.org/sqlite"
)

func DatabaseFind() (err error) {
	cfg := config.Instance
	t := GetTaskFromArgs()
	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		fmt.Printf("failed to open database connection: %v\n", err)
		return
	}
	defer db.Close()
	tasks, err := database_client.FindTasks(db, t)
	if err != nil {
		return
	}
	if len(tasks) == 0 {
		color.Red(`no task found matching criteria`)
		return
	}
	database_client.Display(tasks)
	task := getTask(tasks)
	if task != nil {
		task.Display()
		deleteTask := false
		if err = survey.AskOne(&survey.Confirm{Message: `Do you want to delete this task?`, Default: false}, &deleteTask); err != nil {
			return
		}
		if deleteTask {
			err = database_client.DeleteTask(db, task.ID)
			if err != nil {
				return err
			}
			color.Green(`Task successfully deleted`)
		}
	}
	return
}

func getTaskById(tasks []database_client.Task, id int) *database_client.Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func getTask(tasks []database_client.Task) *database_client.Task {
	taskID := ""
	if err := survey.AskOne(&survey.Input{Message: `Select task (by ID)`}, &taskID); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return nil
	}
	return getTaskById(tasks, id)
}
