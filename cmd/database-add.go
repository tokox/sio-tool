package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Arapak/sio-tool/util"
	"os"
	"path/filepath"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/fatih/color"
	_ "modernc.org/sqlite"
)

func DatabaseAdd() (err error) {
	cfg := config.Instance
	t := GetTaskFromArgs()
	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		fmt.Printf("failed to open database connection: %v\n", err)
		return
	}
	defer db.Close()
	t = ReadTask(t)
	err = database_client.AddTask(db, t)
	if err == nil {
		color.Green(`Task successfully added`)
	}
	return
}

func GetTaskFromArgs() database_client.Task {
	return database_client.Task{Name: Args.Name, Source: Args.Source, Path: Args.Path,
		Link: Args.Link, ShortName: Args.Shortname, ContestID: Args.Contest, ContestStageID: Args.Stage}
}

func checkPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("this path does not exist")
	}
	return nil
}

func ReadTask(task database_client.Task) database_client.Task {
	if task.Source == "" {
		if err := survey.AskOne(&survey.Input{Message: "source"}, &task.Source, survey.WithValidator(survey.Required)); err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
	props := getValuesProperties(task.Source)
	return ExtendTaskInfo(task, props)
}

func ExtendTaskInfo(task database_client.Task, props valuesProperties) database_client.Task {
	if task.Source == "" {
		if err := survey.AskOne(&survey.Input{Message: "source"}, &task.Source); err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
	if task.Name == "" && props.Name.needed {
		util.GetValue("name", &task.Name, props.Name.required)
	}
	if props.Path.needed {
		for {
			if task.Path == "" {
				util.GetValue("path", &task.Path, props.Path.required)
			}
			if task.Path == "" && props.Path.required {
				break
			}
			err := checkPath(task.Path)
			if err == nil {
				break
			}
			color.Red(err.Error())
			task.Path = ""
		}
		if task.Path != "" {
			task.Path, _ = filepath.Abs(task.Path)
		}
	}
	if task.ShortName == "" && props.ShortName.needed {
		util.GetValue("shortname", &task.ShortName, props.ShortName.required)
	}
	if task.Link == "" && props.Link.needed {
		util.GetValue("link", &task.Link, props.Link.required)
	}
	if task.ContestID == "" && props.ContestID.needed {
		util.GetValue("contest", &task.ContestID, props.ContestID.required)
	}
	if task.ContestStageID == "" && props.ContestStageID.needed {
		util.GetValue("stage", &task.ContestStageID, props.ContestStageID.required)
	}
	return task
}

type properties struct {
	needed   bool
	required bool
}

type valuesProperties struct {
	Name           properties
	Path           properties
	ShortName      properties
	Link           properties
	ContestID      properties
	ContestStageID properties
}

func getValuesProperties(source string) valuesProperties {
	switch source {
	case "cf":
		return valuesProperties{
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{false, false},
		}
	case "sio", "sio2":
		return valuesProperties{
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, false},
			properties{false, false},
		}
	case "oi":
		return valuesProperties{
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
		}
	default:
		return valuesProperties{
			properties{true, true},
			properties{true, true},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
		}
	}
}
