package tradeticker

import (
	"sync"

	"github.com/keko950/actoring/actors"
	"github.com/sirupsen/logrus"
)

type ActorManager struct {
	actors       []actors.Actor
	index        int
	maxQueueSize int
	tasks        chan actors.Task
	stop         chan bool
	wg           *sync.WaitGroup
}

func NewActorManager(numA, numT int, wg *sync.WaitGroup) *ActorManager {
	actorsArray := make([]actors.Actor, numA)
	for i := range actorsArray {
		actorsArray[i] = NewTradeActor(numT, wg, make(chan bool))
	}
	manager := &ActorManager{
		actors:       actorsArray,
		maxQueueSize: numA,
		tasks:        make(chan actors.Task, numT),
		stop:         make(chan bool),
		wg:           wg,
	}
	go manager.Start()
	return manager

}

func (e *ActorManager) Start() {
	logrus.Debug("Started manager")
	e.wg.Add(1)
	e.index = 0

	for t := range e.tasks {
		for {

			e.index = e.index % e.maxQueueSize
			err := e.actors[e.index].AddTask(t)
			if err != nil {
				e.index++
				continue
			}

			e.index++
			break

		}
	}

	e.stop <- true

	logrus.Debug("Manager started")

}

func (e *ActorManager) Stop() {
	defer e.wg.Done()
	for _, actor := range e.actors {
		actor.Stop()
	}
	close(e.tasks)
	<-e.stop
	logrus.Debug("Manager stopped")
}

func (e *ActorManager) AddTask(task actors.Task) error {
	if len(e.tasks) >= e.maxQueueSize {
		logrus.Debug("Manager queue is full")
	} else {
		e.tasks <- task
		return nil
	}

	return nil

}
