package database_client

import (
	"database/sql"
	"fmt"
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

func AddTask(db *sql.DB, t Task) error {
	sql := `
        INSERT INTO tasks(name, source, path, shortname, link, contest_id, contest_stage_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(sql, t.Name, t.Source, t.Path, t.ShortName, t.Link, t.ContestID, t.ContestStageID)
	if err != nil {
		return fmt.Errorf("failed to add task to database: %v", err)
	}
	return nil
}

func FindTasks(db *sql.DB, t Task) ([]Task, error) {
	var tasks []Task
	sql := `
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
	rows, err := db.Query(sql, t.Name, t.Source, t.Path, t.ShortName, t.Link, t.ContestID, t.ContestStageID)
	if err != nil {
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
	sql := `
        DELETE FROM tasks
        WHERE id = ?
    `
	_, err := db.Exec(sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete task from database: %v", err)
	}
	return nil
}
