package environment

import (
	"chapter-4/shared/database"
	"chapter-4/shared/tasks"
)

type Env struct {
	Db    *database.DB
	Tasks *tasks.Tasks
}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) InitDB() error {
	db, err := database.NewClient()
	if err != nil {
		return err
	}
	e.Db = db
	return nil
}

func (e *Env) InitTasks() error {
	tasks, err := tasks.NewClient()
	if err != nil {
		return err
	}
	e.Tasks = tasks
	return nil
}
