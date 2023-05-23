package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"

	_ "github.com/mattn/go-sqlite3"
)

func DatabaseGoto() (err error) {
	cfg := config.Instance
	t := GetTaskFromArgs()
	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		fmt.Printf("failed to open database connection: %v\n", err)
		return
	}
	defer db.Close()
	tasks, err := database_client.FindTasks(db, t)
	if len(tasks) == 1 {
		fmt.Printf(tasks[0].Path)
	} else if len(tasks) == 0 {
		return errors.New("no tasks match given criteria")
	} else {
		database_client.Display(tasks)
		return errors.New("more than one task matches given criteria")
	}
	return
}
