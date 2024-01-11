package actors

type Engine interface {
	Run()
	SubmitTask(task Task)
	Shutdown()
}
