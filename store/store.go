package store

import (
	"errors"

	"github.com/mickamy/go_todo_app/entity"
)

var (
	Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}

	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (s *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
	s.LastID++
	t.ID = s.LastID
	s.Tasks[t.ID] = t
	return t.ID, nil
}

func (s *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, 0, len(s.Tasks))
	for i, task := range s.Tasks {
		tasks[i-1] = task
	}
	return tasks
}
