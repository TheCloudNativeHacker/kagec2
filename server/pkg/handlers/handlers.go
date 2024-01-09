package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thecloudnativehacker/kagec2/server/pkg/models"
	"github.com/thecloudnativehacker/kagec2/server/stores"
)

var (
	lock             = sync.Mutex{}
	tasks            = []models.Task{}
	results          = []models.Result{}
	taskHistory      = []models.TaskHistory{}
	taskStore        = stores.NewTaskStore()
	resultStore      = stores.NewResultStore()
	taskHistoryStore = stores.NewTaskHistoryStore()
)

func init() {
	err := taskStore.Load(&tasks)
	if err != nil {
		log.Println("Could not load tasks.")
	}
	err = resultStore.Load(&results)
	if err != nil {
		log.Println("Could not load results")
	}
}
func RootHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "home.page.tmpl", map[string]interface{}{})
}

func GetLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.page.tmpl", map[string]interface{}{})
}

func Login(c echo.Context) error {
	return nil
}

// testing for tasks and results api things
// need to test all the endpoints, urls, with/without uuid, and with/without trash data for the id
func GetTasks(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer taskStore.Save(&tasks)
	if c.QueryParam("id") == "" {
		if c.QueryParam("agent_id") != "" {
			agent_id, _ := uuid.Parse(c.QueryParam("agent_id"))
			filteredTasks := []models.Task{}
			for i := range tasks {
				if tasks[i].AgentId == agent_id {
					filteredTasks = append(filteredTasks, tasks[i])
				}
			}
			return c.JSON(http.StatusOK, filteredTasks)
		}
		return c.JSON(http.StatusOK, tasks)
	}
	id, err := uuid.Parse(c.QueryParam("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Task ID format.")
	}
	for i := range tasks {
		if tasks[i].Id == id {
			log.Println(tasks[i])
			return c.JSON(http.StatusOK, tasks[i])
		}
	}
	return c.String(http.StatusNotFound, "Task not found.")
}

func GetTask(c echo.Context) error {
	// id := uuid.MustParse(c.QueryParam("id"))
	lock.Lock()
	defer lock.Unlock()
	defer taskStore.Save(&tasks)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Task ID format.")
	}
	for i := range tasks {
		if tasks[i].Id == id {
			return c.JSON(http.StatusOK, tasks[i])
		}
	}
	return c.String(http.StatusNotFound, "Task not found.")
}

// need to make sure to enforce certain IDs to be included
func AddTask(c echo.Context) error {
	//need to do additional request validation
	lock.Lock()
	defer lock.Unlock()
	defer taskStore.Save(&tasks)
	task := new(models.Task)
	err := c.Bind(&task)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	i, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Got error: ", err)
	}
	task.Id = i
	log.Println(task)
	tasks = append(tasks, *task)
	return c.JSON(http.StatusOK, task)
}

func GetResult(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer resultStore.Save(&results)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Result ID format.")
	}
	for i := range results {
		if results[i].Id == id {
			return c.JSON(http.StatusOK, results[i])
		}
	}
	return c.String(http.StatusNotFound, "Result not found.")
}

func GetResults(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer resultStore.Save(&results)
	if c.QueryParam("id") == "" {
		return c.JSON(http.StatusOK, results)
	}
	id, err := uuid.Parse(c.QueryParam("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Result ID format.")
	}
	for i := range results {
		if results[i].Id == id {
			log.Println(results[i])
			return c.JSON(http.StatusOK, results[i])
		}
	}
	return c.String(http.StatusNotFound, "Result not found.")
}

// need to make sure to enforce certain IDs to be included
func AddResult(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer resultStore.Save(&results)
	//need to do additional request validation
	result := new(models.Result)
	err := c.Bind(&result)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	i, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Got error: ", err)
	}
	result.Id = i
	log.Println(result)
	results = append(results, *result)
	return c.JSON(http.StatusOK, result)
}

func GetTaskHistory(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer taskHistoryStore.Save(&taskHistory)
	if c.QueryParam("id") == "" {
		return c.JSON(http.StatusOK, taskHistory)
	}
	id, err := uuid.Parse(c.QueryParam("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Result ID format.")
	}
	//get task history for specific task id
	for i := range taskHistory {
		if taskHistory[i].TaskObject.Id == id {
			log.Println(results[i])
			return c.JSON(http.StatusOK, results[i])
		}
	}
	return c.String(http.StatusNotFound, "Result not found.")
}

func GetTaskHistoryById(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer taskHistoryStore.Save(&taskHistory)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error: ", err)
		return c.String(http.StatusBadRequest, "Incorrect Result ID format.")
	}
	for i := range taskHistory {
		if taskHistory[i].TaskObject.Id == id {
			return c.JSON(http.StatusOK, taskHistory[i])
		}
	}
	return c.String(http.StatusNotFound, "Result not found.")
}

func AddTaskHistory(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	defer taskHistoryStore.Save(&taskHistory)
	//need to do additional request validation
	taskH := new(models.TaskHistory)
	err := c.Bind(&taskH)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	//need to get the Task object and the Results object
	log.Println(taskH)
	taskHistory = append(taskHistory, *taskH)
	return c.JSON(http.StatusOK, taskH)
}

// func GetImplants(c echo.Context) error {
// }

// func GetImplant(c echo.Context) error {
// }

// func AddImplant(c echo.Context) error {
// }
