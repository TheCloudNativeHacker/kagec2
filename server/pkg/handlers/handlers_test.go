package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/thecloudnativehacker/kagec2/server/pkg/models"
	_ "github.com/thecloudnativehacker/kagec2/server/pkg/models"
)

var (
	taskJSON   = `{"type":"ping", "agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4"}`
	resultJSON = `{"contents":"asdf", "agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4"}`
)

// will need to add test for a task id being put into the request
// want this to not overwrite the task id created, also want to make sure there is a valid agent id eventually
func TestAddTask(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(taskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddTask(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}

	respTask := models.Task{}
	err = json.Unmarshal([]byte(rec.Body.String()), &respTask)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if respTask.Type != "ping" {
		t.Errorf("Expected task type to be ping got %v", respTask.Type)
	}
	if respTask.Id.String() == "00000000-0000-0000-0000-000000000000" {
		t.Errorf("Expected random task UUID got default UUID instead: %v", respTask.Id)
	}
	if respTask.AgentId.String() != "49c10782-cfe9-472a-9c86-d66d641e9ca4" {
		t.Errorf("Expected agent UUID to be 49c10782-cfe9-472a-9c86-d66d641e9ca4 got: %v instead", respTask.AgentId)
	}
}

// func TestGetTask(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(taskJSON))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	err := AddTask(c)
// 	if err != nil {
// 		t.Errorf("Error adding task: %v", err)
// 	}

// 	//need to get the task and check, have to figure out the way to get the uuid
// }

// need to tesk that result matches with a task and agent
func TestAddResult(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/results", strings.NewReader(resultJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddResult(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}

	// will need to flesh this out with more fully defined results and tasks
	respResult := models.Result{}
	err = json.Unmarshal([]byte(rec.Body.String()), &respResult)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if respResult.Contents != "asdf" {
		t.Errorf("Unexpected contents, expected asdf got: %v", respResult.Contents)
	}
	if respResult.AgentId.String() != "49c10782-cfe9-472a-9c86-d66d641e9ca4" {
		t.Errorf("Expected agent UUID to be 49c10782-cfe9-472a-9c86-d66d641e9ca4 got: %v instead", respResult.AgentId)
	}
}
