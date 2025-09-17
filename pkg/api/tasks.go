package api

import (
	"fmt"
	"net/http"
	"strconv"

	"go-final-project/pkg/db"
)

type taskDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type tasksResp struct {
	Tasks []taskDTO `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не поддерживается"})

		return
	}

	const limit = 50

	rows, err := db.Tasks(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("ошибка выборки: %v", err)})

		return
	}

	resp := tasksResp{Tasks: make([]taskDTO, 0, len(rows))}
	for _, t := range rows {
		resp.Tasks = append(resp.Tasks, taskDTO{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
