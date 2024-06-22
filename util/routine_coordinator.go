package util

import "sync"

type RoutineTask func(...any)

// Need to be initialized by initialize method on creation before use
// coordinator := util.RoutineCoordinator{}
// coordinator.Initialize(8)
//
//	for i := 0; i < 100; i++ {
//		workFunc := func(vals ...any) {
//			id := vals[0].(int)
//			time.Sleep(time.Second * 1)
//			fmt.Printf("Worker Id %d \n", id)
//			fmt.Printf("")
//		}
//		coordinator.Add(workFunc, i)
//	}
//
// coordinator.WaitForCompletion()

type RoutineCoordinator struct {
	wg    sync.WaitGroup
	guard chan struct{}
}

// initialize the coordinator
// @Param maxCoroutine int need to pass maximum concurrent running go routines. Minimum allowed value is 1.
func (c *RoutineCoordinator) Initialize(maxCoroutine int) {
	if maxCoroutine < 1 {
		maxCoroutine = 1
	}
	c.guard = make(chan struct{}, maxCoroutine)
	c.wg = sync.WaitGroup{}
}

// @warning dont use variable defined outside task function directly in task func. result may be inconsistent and buggy
func (c *RoutineCoordinator) Add(task RoutineTask, values ...any) {
	c.wg.Add(1)
	c.guard <- struct{}{}

	workFunc := func(wg *sync.WaitGroup, task RoutineTask, guard *chan struct{}) {
		defer wg.Done()
		task(values...)
		<-*guard
	}
	go workFunc(&c.wg, task, &c.guard)
}

func (c *RoutineCoordinator) WaitForCompletion() {
	c.wg.Wait()
}
