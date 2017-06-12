package main

type Dispatcher struct {
        // A pool of workers channels that are registered with the dispatcher
        WorkerPool chan chan Job
        maxWorkers int
}


func NewDispatcher(maxWorkers int) *Dispatcher {
        pool := make(chan chan Job, maxWorkers)
        return &Dispatcher{WorkerPool: pool, maxWorkers:maxWorkers}
}

func (d *Dispatcher) Run() {
        // starting n number of workers
        for i := 0; i < d.maxWorkers; i++ {
                // id starts from 1
                worker := NewWorker(i + 1, d.WorkerPool)
                worker.Start()
        }
        go d.dispatch()
}

func (d *Dispatcher) dispatch() {
        for {
                select {
                case job := <-JobQueue:
                // a job request has been received
                        go func() {
                                // try to obtain a worker job channel that is available.
                                // this will block until a worker is idle
                                workerJobQueue := <-d.WorkerPool
                                // dispatch the job to the worker job channel
                                workerJobQueue <- job
                        }()
                }
        }
}
