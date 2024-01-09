package models

import (
	"net/http"

	"github.com/google/uuid"
)

type Task struct {
	Id      uuid.UUID `param:"id" json:"id" query:"id" db:"id"`
	AgentId uuid.UUID `param:"agent_id" json:"agent_id" query:"agent_id" db:"agent_id"`
	Type    string    `json:"type" query:"type" db:"type"`
}

type Result struct {
	Id       uuid.UUID `param:"id" json:"id" query:"id" db:"id"`
	TaskID   uuid.UUID `param:"task_id" json:"task_id" query:"task_id" db:"task_id"`
	AgentId  uuid.UUID `param:"agent_id" json:"agent_id" query:"agent_id" db:"agent_id"`
	Contents string    `json:"contents" db:"contents"`
}

type TaskHistory struct {
	TaskObject Task     `json:"task_object"`
	Options    []string `json:"task_options"`
	TaskResult Result   `json:"task_results"`
}

// type Agent struct {
// 	Id         uuid.UUID         `param:"id" json:"id" query:"id" db:"id"`
// 	Name       string            `json:"name" query:"name" db:"name"`
// 	Interfaces map[string]string `json:"interfaces" db:"interfaces"`
// }

type Implant struct {
	Id                       uuid.UUID `param:"id" json:"id" query:"id" db:"id"`
	C2Host                   string
	C2Port                   string
	C2TasksURI               string
	C2ResultsURI             string
	DwellDistributionSeconds float64
	Client                   http.Client
	Running                  bool        `param:"running" json:"running" query:"running" db:"running"`
	Results                  *[]Result   `param:"results" json:"results" query:"results" db:"results"`
	Tasks                    *[]Task     `param:"tasks" json:"tasks" query:"tasks" db:"tasks"`
	Info                     MachineInfo `param:"info" json:"info" query:"info" db:"info"`
}

type MachineInfo struct {
	Interfaces  map[string]map[string][]string `param:"interfaces" json:"interfaces" query:"interfaces" db:"interfaces"`
	Hostname    string                         `param:"hostname" json:"hostname" query:"hostname" db:"hostname"`
	CurrentUser string                         `param:"current_user" json:"current_user" query:"current_user" db:"current_user"`
	CurrentUId  string                         `param:"current_uid" json:"current_uid" query:"current_uid" db:"current_uid"`
	CurrentGId  string                         `param:"current_gid" json:"current_gid" query:"current_gid" db:"current_gid"`
	OS          string                         `param:"os" json:"os" query:"os" db:"os"`
	Arch        string                         `param:"arch" json:"arch" query:"arch" db:"arch"`
	PID         int                            `param:"pid" json:"pid" query:"pid" db:"pid"`
}
