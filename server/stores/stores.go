package stores

import (
	"encoding/json"
	"os"

	"github.com/thecloudnativehacker/kagec2/server/pkg/models"
)

type (
	TaskStore interface {
		Save(ts *[]models.Task) error
		Load(ts *[]models.Task) error
	}
	ResultStore interface {
		Save(rs *[]models.Result) error
		Load(rs *[]models.Result) error
	}
	//agents and taskhistory not fleshed out yet
	TaskHistoryStore interface {
		Save(th *[]models.TaskHistory) error
		Load(th *[]models.TaskHistory) error
	}
	AgentStore interface {
		Save(a *[]models.Implant) error
		Load(a *[]models.Implant) error
	}

	taskStore struct {
		file  string
		tasks *[]models.Task
	}

	resultStore struct {
		file    string
		results *[]models.Result
	}

	taskHistoryStore struct {
		file  string
		tasks *[]models.TaskHistory
	}
)

const (
	tasksFile       = "tasks.json"
	agentsFile      = "agents.json"
	resultsFile     = "results.json"
	taskHistoryFile = "taskhistory.json"
)

func (t *taskStore) Save(ts *[]models.Task) error {
	tasks, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	err = os.WriteFile(t.file, tasks, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskStore) Load(ts *[]models.Task) error {
	content, err := os.ReadFile(t.file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, ts)
	if err != nil {
		return err
	}
	return nil
}

func NewTaskStore() TaskStore {
	t := taskStore{file: tasksFile}
	return &t
}

func (r *resultStore) Save(rs *[]models.Result) error {
	tasks, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	err = os.WriteFile(r.file, tasks, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (r *resultStore) Load(rs *[]models.Result) error {
	content, err := os.ReadFile(r.file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, rs)
	if err != nil {
		return err
	}
	return nil
}

func NewResultStore() ResultStore {
	r := resultStore{file: resultsFile}
	return &r
}

func (t *taskHistoryStore) Save(ts *[]models.TaskHistory) error {
	taskHistory, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	err = os.WriteFile(t.file, taskHistory, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskHistoryStore) Load(ts *[]models.TaskHistory) error {
	content, err := os.ReadFile(t.file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, ts)
	if err != nil {
		return err
	}
	return nil
}

func NewTaskHistoryStore() TaskHistoryStore {
	t := taskHistoryStore{file: taskHistoryFile}
	return &t
}
