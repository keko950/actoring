package main

import (
	"fmt"
	"sync"

	"github.com/keko950/actoring/example/tradeticker"
	"github.com/sirupsen/logrus"
)

type PrintTask struct {
}

func (p *PrintTask) Execute() {
	fmt.Println("Hello World from an Actor!")
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	engine := tradeticker.NewEngine(5, 3, &sync.WaitGroup{})
	engine.Run()

	for i := 1; i <= 10; i++ {
		t := &PrintTask{}
		engine.SubmitTask(t)
		fmt.Println("enviado")
	}

	engine.Shutdown()

}
