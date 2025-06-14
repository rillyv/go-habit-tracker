package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rillyv/habit-tracker/graph/model"
)

func GetHabitByID(id string) (*model.Habit, error) {
	var habit model.Habit

	query := `SELECT id, title, description, created_at FROM habits WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&habit.ID, &habit.Title, &habit.Description, &habit.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("Failed to get habit: %w", err)
	}

	return &habit, nil
}

func ListHabits() ([]*model.Habit, error) {
	query := `SELECT id, title, description, created_at FROM habits ORDER BY created_at DESC`
	rows, err := DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("Failed to query habits: %w", err)
	}

	defer rows.Close()

	var habits []*model.Habit
	for rows.Next() {
		var habit model.Habit
		if err := rows.Scan(&habit.ID, &habit.Title, &habit.Description, &habit.CreatedAt); err != nil {
			return nil, fmt.Errorf("Failed to scan habit: %w", err)
		}
		habits = append(habits, &habit)
	}

	return habits, nil
}

func InsertHabit(input model.CreateHabitInput) (*model.Habit, error) {
	var id int
	createdAt := time.Now()

	query := `INSERT INTO habits (title, description, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := DB.QueryRow(query, input.Title, input.Description, createdAt).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("Failed to insert habit: %w", err)
	}

	return &model.Habit{
		ID:          strconv.FormatInt(int64(id), 10),
		Title:       input.Title,
		Description: input.Description,
		CreatedAt:   createdAt.Format(time.RFC1123Z),
	}, nil
}

func UpdateHabit(input model.UpdateHabitInput) (*model.Habit, error) {
	var habit model.Habit
	query := `UPDATE habits SET title = COALESCE($2, title), description = COALESCE($3, description) WHERE id = $1 RETURNING id, title, description, created_at`
	err := DB.QueryRow(query, input.ID, input.Title, input.Description).Scan(&habit.ID, &habit.Title, &habit.Description, &habit.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("Failed to update habit: %w", err)
	}

	return &habit, nil
}

func DeleteHabit(input model.DeleteHabitInput) (*model.Habit, error) {
	var habit model.Habit

	query := `DELETE FROM habits WHERE id = $1 RETURNING id, title, description, created_at`
	err := DB.QueryRow(query, input.ID).Scan(&habit.ID, &habit.Title, &habit.Description, &habit.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("Failed to delete habit: %w", err)
	}

	return &habit, nil
}
