package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thecloudnativehacker/kagec2/server/pkg/models"
)

var (
	taskJSON           = `{"type":"ping", "agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4"}`
	nilTaskJSON        = `{"type":"ping"}`
	resultJSON         = `{"contents":"asdf", "agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4"}`
	nilAgentResultJSON = `{"contents":"asdf"}`
	nilTaskResultJSON  = `{"contents":"asdf", "agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4"}`
	taskHistoryJSON    = `{"task":{"id":"11111111-1111-1111-1111-111111111111","agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4","type":"ping"},"options":["option1","option2"],"results":{"id":"22222222-2222-2222-2222-222222222222","task_id":"11111111-1111-1111-1111-111111111111","agent_id":"49c10782-cfe9-472a-9c86-d66d641e9ca4","contents":"pong"}}`
	userJSON           = `{"username":"test_user", "email":"test_user@protonmail.com", "password":"mypassword"}`
	userNoEmailJSON    = `{"username":"test_user", "password":"mypassword"}`
	userNoNameJSON     = `{"email":"test_user@protonmail.com", "password":"mypassword"}`
	userNoPassJSON     = `{"username":"test_user", "email":"test_user@protonmail.com"}`
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
	err = json.Unmarshal(rec.Body.Bytes(), &respTask)
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

func TestAddNilAgentIdTask(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(nilTaskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = AddTask(c)
	if c.Response().Status != 400 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
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
	// need to add a task first get that task id and use it to construct the result json
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
	err = json.Unmarshal(rec.Body.Bytes(), &respTask)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}

	// need to get the task id here.
	res := models.Result{}
	res.TaskId = respTask.Id
	res.AgentId = respTask.AgentId
	res.Contents = "test contents"
	res.Id, err = uuid.NewRandom()
	if err != nil {
		t.Errorf("Error generating uuid for result: %v", err)
	}
	mRes, err := json.Marshal(res)
	if err != nil {
		t.Errorf("Error Marshalling Response: %v", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/api/results", strings.NewReader(string(mRes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = AddResult(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}

	// will need to flesh this out with more fully defined results and tasks
	respResult := models.Result{}
	err = json.Unmarshal(rec.Body.Bytes(), &respResult)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if respResult.Contents != "test contents" {
		t.Errorf("Unexpected contents, expected test contents got: %v", respResult.Contents)
	}
	if respResult.AgentId.String() != "49c10782-cfe9-472a-9c86-d66d641e9ca4" {
		t.Errorf("Expected agent UUID to be 49c10782-cfe9-472a-9c86-d66d641e9ca4 got: %v instead", respResult.AgentId)
	}
}

// need to fix the checking to enusre the right error is being returned
func TestAddNilAgentIdResult(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/results", strings.NewReader(nilAgentResultJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = AddResult(c)
	if c.Response().Status != 400 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}
}

// need to fix the checking to enusre the right error is being returned
func TestAddNilTaskIdResult(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/results", strings.NewReader(nilTaskResultJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = AddResult(c)
	if c.Response().Status != 400 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}
}

// TODO need to implement check for valid agent first
func TestAddBadAgentIdResult(t *testing.T) {
}

func TestAddTaskHistory(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/taskhistory", strings.NewReader(taskHistoryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddTaskHistory(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}

	respHistory := models.TaskHistory{}
	err = json.Unmarshal(rec.Body.Bytes(), &respHistory)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if strings.Join(respHistory.Options, ",") != "option1,option2" {
		t.Errorf("Unexpected options, expected option1,option2 got: %v", strings.Join(respHistory.Options, ","))
	}
	if respHistory.TaskObject.AgentId.String() != "49c10782-cfe9-472a-9c86-d66d641e9ca4" {
		t.Errorf("Expected task agent UUID to be 49c10782-cfe9-472a-9c86-d66d641e9ca4 got: %v instead", respHistory.TaskObject.AgentId.String())
	}
	if respHistory.TaskObject.Id.String() != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("Expected  task UUID to be 11111111-1111-1111-1111-111111111111 got: %v instead", respHistory.TaskObject.Id.String())
	}
	if respHistory.TaskObject.Type != "ping" {
		t.Errorf("Expected Task Type to be ping, got %v", respHistory.TaskObject.Type)
	}
	if respHistory.TaskResult.AgentId.String() != "49c10782-cfe9-472a-9c86-d66d641e9ca4" {
		t.Errorf("Expected result agent UUID to be 49c10782-cfe9-472a-9c86-d66d641e9ca4 got: %v instead", respHistory.TaskResult.AgentId.String())
	}
	if respHistory.TaskResult.Id.String() != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("Expected result task UUID to be 22222222-2222-2222-2222-222222222222 got: %v instead", respHistory.TaskResult.Id.String())
	}
	if respHistory.TaskResult.TaskId.String() != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("Expected result task UUID to be 11111111-1111-1111-1111-111111111111 got: %v instead", respHistory.TaskResult.TaskId.String())
	}
	if respHistory.TaskResult.Contents != "pong" {
		t.Errorf("Expected Task Result Contents to be pong, got %v", respHistory.TaskObject.Type)
	}
}

func TestAddUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddUser(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}
	respUser := models.User{}
	err = json.Unmarshal(rec.Body.Bytes(), &respUser)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if respUser.Name != "test_user" {
		t.Errorf("Expected Username: test_user got: %v", respUser.Name)
	}
	if respUser.Email != "test_user@protonmail.com" {
		t.Errorf("Expected Email: test_user@protonmail.com got: %v", respUser.Email)
	}
}

// should pass
func TestAddUserWithoutEmail(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userNoEmailJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddUser(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 200 {
		t.Errorf("Expected 200 response got %v", c.Response().Status)
	}
	respUser := models.User{}
	err = json.Unmarshal(rec.Body.Bytes(), &respUser)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %v", err)
	}
	if respUser.Name != "test_user" {
		t.Errorf("Expected Username: test_user got: %v", respUser.Name)
	}
}

// should fail
func TestAddUserWithoutName(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userNoNameJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddUser(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 400 {
		t.Errorf("Expected 400 response got %v", c.Response().Status)
	}
	if rec.Body.String() != "Could not validate request." {
		t.Errorf("Expected error validating request, received none.")
	}
}

// should fail
func TestAddUserWithoutPass(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userNoPassJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := AddUser(c)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if c.Response().Status != 400 {
		t.Errorf("Expected 400 response got %v", c.Response().Status)
	}
	if rec.Body.String() != "Could not validate request." {
		t.Errorf("Expected error validating request, received none.")
	}
}

// should fail
func TestAddUserWithShortPass(t *testing.T) {
}

// should fail
func TestAddUserWithLongPass(t *testing.T) {
}
