package db

import "fmt"

type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task) (int64, error) {
	const q = `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?);
	`

	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC, id ASC
		LIMIT ?;
	`

	rows, err := DB.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		t := new(Task)
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	if out == nil {
		out = make([]*Task, 0)
	}
	return out, nil
}

func GetTask(id string) (*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?;
	`
	var t Task
	err := DB.QueryRow(q, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(t *Task) error {
	const q = `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?;
	`
	res, err := DB.Exec(q, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}
