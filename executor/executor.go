package executor

/*
Executor interface and default executor implementation is defined here.
*/
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/executor/events"
	exec "mesos-framework-sdk/include/executor"
	"mesos-framework-sdk/include/mesos"
)

type Executor interface {
	FrameworkID() *mesos_v1.FrameworkID
	ExecutorID() *mesos_v1.ExecutorID
	Client() *client.Client
	Events() chan *exec.Event
	Subscribe()
	Update(taskStatus *mesos_v1.TaskStatus)
	Message(data []byte)
}

type DefaultExecutor struct {
	frameworkID *mesos_v1.FrameworkID
	executorID  *mesos_v1.ExecutorID
	client      *client.Client
	events      chan *exec.Event
	handlers    events.ExecutorEvents
}

// Creates a new default executor
func NewDefaultExecutor(c *client.Client) *DefaultExecutor {
	return &DefaultExecutor{
		frameworkID: &mesos_v1.FrameworkID{Value: proto.String("")},
		executorID:  &mesos_v1.ExecutorID{Value: proto.String("")},
		client:      c,
		events:      make(chan *exec.Event),
	}

}

// Default listening method on the
func (d *DefaultExecutor) listen() {
	for {
		switch t := <-d.events; t.GetType() {
		case exec.Event_SUBSCRIBED:
			go d.handlers.Subscribed()
			break
		case exec.Event_ACKNOWLEDGED:
			go d.handlers.Acknowledged()
			break
		case exec.Event_MESSAGE:
			go d.handlers.Message()
			break
		case exec.Event_KILL:
			go d.handlers.Kill()
			break
		case exec.Event_LAUNCH:
			go d.handlers.Launch()
			break
		case exec.Event_LAUNCH_GROUP:
			go d.handlers.LaunchGroup()
			break
		case exec.Event_SHUTDOWN:
			go d.handlers.Shutdown()
			break
		case exec.Event_ERROR:
			go d.handlers.Error()
			break
		case exec.Event_UNKNOWN:
			fmt.Println("Unknown event caught.")
			break
		}
	}
}

func (d *DefaultExecutor) FrameworkID() {
	return d.frameworkID
}

func (d *DefaultExecutor) ExecutorID() {
	return d.executorID
}

func (d *DefaultExecutor) Subscribe() {
	// Both id's for framework and executor will be empty here.
	subscribe := &exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_SUBSCRIBE.Enum(),
	}
	d.client.Request(subscribe)
}

func (d *DefaultExecutor) Update(taskStatus *mesos_v1.TaskStatus) {
	update := exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_UPDATE.Enum(),
		Update: &exec.Call_Update{
			Status: taskStatus,
		},
	}
	d.client.Request(update)

}
func (d *DefaultExecutor) Message(data []byte) {
	message := exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_MESSAGE.Enum(),
		Message: &exec.Call_Message{
			Data: data,
		},
	}
	d.client.Request(message)

}
