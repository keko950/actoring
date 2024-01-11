package tradeticker

import (
	"errors"
	"sync"

	"github.com/keko950/actoring/actors"
	"github.com/sirupsen/logrus"
)

type TradeActor struct {
	tc   chan actors.Task
	stop chan bool
	wg   *sync.WaitGroup
}

func NewTradeActor(numT int, wg *sync.WaitGroup, stopc chan bool) *TradeActor {

	actor := &TradeActor{
		tc:   make(chan actors.Task, numT),
		stop: stopc,
		wg:   wg,
	}

	go actor.Start()
	return actor
}

func (e *TradeActor) AddTask(task actors.Task) error {

	select {
	case e.tc <- task:
		logrus.Debug("Sent task")
	default:
		return errors.New("channel is full")
	}

	return nil

}

func (e *TradeActor) Start() {
	e.wg.Add(1)
	logrus.Debug("Actor started")

	for t := range e.tc {
		t.Execute()
	}

	e.stop <- true

}

func (e *TradeActor) Stop() {
	close(e.tc)
	<-e.stop

	e.wg.Done()
	logrus.Debug("Actor terminated")

}
