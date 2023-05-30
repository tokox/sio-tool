package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Arapak/sio-tool/config"
	"github.com/Arapak/sio-tool/database_client"
	"github.com/Arapak/sio-tool/util"

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

func readValue(valueName string, nilable bool) (value string) {
	for {
		fmt.Printf("%v: ", valueName)
		value = util.ScanlineTrim()
		if !nilable && value == "" {
			color.Red(`value cannot be empty`)
		} else {
			return
		}
	}
}

func ReadTask(task database_client.Task) database_client.Task {
	if task.Source == "" {
		task.Source = readValue("source", false)
	}
	props := getValuesProperties(task.Source)
	return ExtendTaskInfo(task, props)
}

func ExtendTaskInfo(task database_client.Task, props valuesProperties) database_client.Task {
	if task.Source == "" {
		task.Source = readValue("source", true)
	}
	if task.Name == "" && props.Name.needed {
		task.Name = readValue("name", props.Name.nilable)
	}
	if props.Path.needed {
		for {
			if task.Path == "" {
				task.Path = readValue("path", props.Path.nilable)
			}
			if task.Path == "" && props.Path.nilable {
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
		task.ShortName = readValue("shortname", props.ShortName.nilable)
	}
	if task.Link == "" && props.Link.needed {
		task.Link = readValue("link", props.Link.nilable)
	}
	if task.ContestID == "" && props.ContestID.needed {
		task.ContestID = readValue("contest", props.ContestID.nilable)
	}
	if task.ContestStageID == "" && props.ContestStageID.needed {
		task.ContestStageID = readValue("stage", props.ContestStageID.nilable)
	}
	return task
}

type properties struct {
	needed  bool
	nilable bool
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
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{false, true},
		}
	case "sio", "sio2":
		return valuesProperties{
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, true},
			properties{false, true},
		}
	case "oi":
		return valuesProperties{
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
			properties{true, false},
		}
	default:
		return valuesProperties{
			properties{true, false},
			properties{true, false},
			properties{true, true},
			properties{true, true},
			properties{true, true},
			properties{true, true},
		}
	}
}
