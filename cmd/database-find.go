package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/Arapak/sio-tool/util"

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
	task, selected := getTask(tasks)
	if selected {
		task.Display()
		deleteTask := askForDeletion()
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

func getTaskById(tasks []database_client.Task, id int) (database_client.Task, error) {
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return database_client.Task{}, errors.New("task with this id not found")
}

func getTask(tasks []database_client.Task) (task database_client.Task, selected bool) {
	color.Green(`Select task (by ID): `)
	for {
		value := util.ScanlineTrim()
		if value == "" {
			return database_client.Task{}, false
		}
		id, err := strconv.Atoi(value)
		if err != nil || id < 0 {
			color.Red(`this is not a positive integer`)
			continue
		}
		task, err = getTaskById(tasks, id)
		if err != nil {
			color.Red(err.Error())
			continue
		}
		return task, true
	}
}

func askForDeletion() bool {
	color.Green(`Do you want to delete this task (y/N): `)
	value := util.ScanlineTrim()
	return value == "y" || value == "Y"
}
