package scheduler

type Job interface {
	Do()
	Cancel()
}
