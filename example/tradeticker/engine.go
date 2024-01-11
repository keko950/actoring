package tradeticker

import (
	"sync"

	"github.com/keko950/actoring/actors"
	"github.com/sirupsen/logrus"
)

type TradeEngine struct {
	wg      *sync.WaitGroup
	manager actors.Actor
}

func NewEngine(numA, numT int, wg *sync.WaitGroup) *TradeEngine {
	return &TradeEngine{
		wg:      &sync.WaitGroup{},
		manager: NewActorManager(numA, numT, wg),
	}
}

func (c *TradeEngine) Run() {
	logrus.Debug("Starting Engine")
	go c.manager.Start()
}

func (c *TradeEngine) SubmitTask(task actors.Task) {
	logrus.Debug("Submitting Task")
	c.manager.AddTask(task)
}

func (c *TradeEngine) Shutdown() {
	c.manager.Stop()
	c.wg.Wait()
	logrus.Debug("Engine shutdown sucessfully")
}
