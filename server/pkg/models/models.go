package models

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// todo extend models to include request time and execution time
type Task struct {
	Type    string    `json:"type" query:"type" db:"type"`
	Id      uuid.UUID `param:"id" json:"id" query:"id" db:"id"`
	AgentId uuid.UUID `param:"agent_id" json:"agent_id" query:"agent_id" db:"agent_id"`
}

type Result struct {
	Contents string    `json:"contents" db:"contents"`
	Id       uuid.UUID `param:"id" json:"id" query:"id" db:"id"`
	TaskId   uuid.UUID `param:"task_id" json:"task_id" query:"task_id" db:"task_id"`
	AgentId  uuid.UUID `param:"agent_id" json:"agent_id" query:"agent_id" db:"agent_id"`
}

type TaskHistory struct {
	TaskObject Task     `json:"task"`
	Options    []string `json:"options"`
	TaskResult Result   `json:"results"`
}

type Implant struct {
	C2Host                   string
	C2Port                   string
	C2TasksURI               string
	C2ResultsURI             string
	Results                  *[]Result `param:"results" json:"results" query:"results" db:"results"`
	Tasks                    *[]Task   `param:"tasks" json:"tasks" query:"tasks" db:"tasks"`
	Client                   http.Client
	Info                     MachineInfo `param:"info" json:"info" query:"info" db:"info"`
	Id                       uuid.UUID   `param:"id" json:"id" query:"id" db:"id"`
	DwellDistributionSeconds float64
	Running                  bool `param:"running" json:"running" query:"running" db:"running"`
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

type User struct {
	Name     string    `json:"username" db:"username" validate:"required"`
	Email    string    `json:"email,omitempty" db:"email,omitempty" validate:"omitempty,email"`
	Password string    `json:"password" db:"password" validate:"required,min=8,max=300"`
	Id       uuid.UUID `json:"id" db:"id"`
}

type UserValidator struct {
	Validator *validator.Validate
}

func (u *UserValidator) Validate(i interface{}) error {
	return u.Validator.Struct(i)
}

type ServiceAccount struct{}

// TODO this will be the general interface for anything that needs to authenticate against the API
type Principal interface{}
