package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-final-project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не поддерживается"})

		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})

		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})

			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("ошибка выборки: %v", err)})

		return
	}

	rep := strings.TrimSpace(task.Repeat)
	if rep == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})

			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})

		return
	}

	now := time.Now()

	next, err := NextDate(now, task.Date, rep)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("некорректное правило повторения: %v", err)})

		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})

		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
