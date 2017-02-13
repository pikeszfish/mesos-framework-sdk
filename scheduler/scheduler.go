package scheduler

import (
	"fmt"
	"log"
	"mesos-framework-sdk/client"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/scheduler/events"
	"time"
)

const (
	subscribeRetry = 2
)

type Scheduler struct {
	frameworkID mesos.FrameworkID
	client      *client.Client
	events      chan *sched.Event
}

func NewScheduler(c *client.Client) *Scheduler {
	return &Scheduler{
		client: c,
		events: make(chan *sched.Event),
	}
}

// Create a Subscription to mesos.
func (c *Scheduler) Subscribe(frameworkInfo *mesos.FrameworkInfo) (<-chan *sched.Event, error) {
	// We really want the ID after the call...
	call := &sched.Call{
		FrameworkId: frameworkInfo.GetId(),
		Type:        sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: frameworkInfo,
		},
	}

	for {
		resp, err := c.client.Request(call)
		if err != nil {
			log.Println(err.Error())
		} else {
			// Once connected the client should set our framework ID for all outgoing calls after successful subscribe.
			go events.Loop(resp.Body, c.events)
			break
		}

		time.Sleep(time.Duration(subscribeRetry) * time.Second)
	}

	return c.events, nil
}

// Send a teardown request to mesos master.
func (c *Scheduler) Teardown() {
	teardown := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_TEARDOWN.Enum(),
	}
	resp, err := c.client.Request(teardown)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

// Skeleton funcs for the rest of the calls.

// Accepts offers from mesos master
func (c *Scheduler) Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters) {
	accept := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      &sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.client.Request(accept)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *Scheduler) Decline(offerIds []*mesos.OfferID, filters *mesos.Filters) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     &sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.client.Request(decline)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *Scheduler) Revive() {

	revive := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.client.Request(revive)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *Scheduler) Kill(taskId *mesos.TaskID, agentid *mesos.AgentID) {
	// Probably want some validation that this is a valid task and valid agentid.
	kill := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_KILL.Enum(),
		Kill:        &sched.Call_Kill{TaskId: taskId, AgentId: agentid},
	}

	resp, err := c.client.Request(kill)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *Scheduler) Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID) {
	shutdown := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_SHUTDOWN.Enum(),
		Shutdown: &sched.Call_Shutdown{
			ExecutorId: execId,
			AgentId:    agentId,
		},
	}
	resp, err := c.client.Request(shutdown)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// UUID should be a type
// TODO import extras uuid funcs.
func (c *Scheduler) Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte) {
	acknowledge := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &sched.Call_Acknowledge{AgentId: agentId, TaskId: taskId, Uuid: uuid},
	}
	resp, err := c.client.Request(acknowledge)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *Scheduler) Reconcile(tasks []*mesos.Task) {
	reconcile := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_RECONCILE.Enum(),
	}
	resp, err := c.client.Request(reconcile)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *Scheduler) Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte) {
	message := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_MESSAGE.Enum(),
		Message: &sched.Call_Message{
			AgentId:    agentId,
			ExecutorId: executorId,
			Data:       data,
		},
	}
	resp, err := c.client.Request(message)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)

}

// Sent by the scheduler to request resources from the master/allocator.
// The built-in hierarchical allocator simply ignores this request but other allocators (modules) can interpret this in
// a customizable fashion.
func (c *Scheduler) SchedRequest(resources []*mesos.Request) {
	request := &sched.Call{
		FrameworkId: &c.frameworkID,
		Type:        sched.Call_REQUEST.Enum(),
		Request: &sched.Call_Request{
			Requests: resources,
		},
	}
	resp, err := c.client.Request(request)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}