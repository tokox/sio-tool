package database_client

import (
	"database/sql"
	"fmt"
	"strings"
)

type Task struct {
	ID             int
	Name           string
	Source         string
	Path           string
	ShortName      string
	Link           string
	ContestID      string
	ContestStageID string
}

func createTableIfNotExist(db *sql.DB) error {
	sqlStatement := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			source TEXT NOT NULL,
			path TEXT NOT NULL,
			shortname TEXT,
			link TEXT,
			contest_id TEXT,
			contest_stage_id TEXT,
			UNIQUE(name,source,path,shortname,link,contest_id,contest_stage_id)
		);
  `
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

func AddTask(db *sql.DB, t Task) error {
	sqlStatement := `
        INSERT INTO tasks(name, source, path, shortname, link, contest_id, contest_stage_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(sqlStatement, t.Name, t.Source, t.Path, t.ShortName, t.Link, t.ContestID, t.ContestStageID)
	if err != nil {
		if strings.Contains(err.Error(), `no such table: tasks`) {
			err = createTableIfNotExist(db)
			if err != nil {
				return err
			}
			return AddTask(db, t)
		} else if strings.Contains(err.Error(), `constraint failed: UNIQUE constraint failed`) {
			return fmt.Errorf("this problem already exists in database")
		}
		return fmt.Errorf("failed to add task to database: %v", err)
	}
	return nil
}

func FindTasks(db *sql.DB, t Task) ([]Task, error) {
	var tasks []Task
	sqlStatement := `
	    SELECT id, name, source, path, shortname, link, contest_id, contest_stage_id
	    FROM tasks
	    WHERE LOWER(name) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(source) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(path) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(shortname) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(link) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(contest_id) LIKE '%' || LOWER(?) || '%'
	    AND LOWER(contest_stage_id) LIKE '%' || LOWER(?) || '%'
	`
	rows, err := db.Query(sqlStatement, t.Name, t.Source, t.Path, t.ShortName, t.Link, t.ContestID, t.ContestStageID)
	if err != nil {
		if strings.Contains(err.Error(), `no such table: tasks`) {
			err = createTableIfNotExist(db)
			if err != nil {
				return nil, err
			}
			return FindTasks(db, t)
		}
		return nil, fmt.Errorf("failed to find tasks in database: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Source, &task.Path, &task.ShortName, &task.Link, &task.ContestID, &task.ContestStageID); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read task rows: %v", err)
	}
	return tasks, nil
}

func DeleteTask(db *sql.DB, id int) error {
	sqlStatement := `
        DELETE FROM tasks
        WHERE id = ?
    `
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("failed to delete task from database: %v", err)
	}
	return nil
}
