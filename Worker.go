package main

import "log"
// Worker represents the worker that executes the job
type Worker struct {
        id         int
        WorkerPool  chan chan Job
        JobChannel  chan Job
        quitChannel    	chan bool
}

// invoked from dispatcher
func NewWorker(id int, workerPool chan chan Job) Worker {
        return Worker{
                id: id,
                WorkerPool: workerPool,
                JobChannel: make(chan Job),
                quitChannel:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in case we need to stop it
func (w Worker) Start() {
        go func() {
                for {
                        // register the current worker into the worker queue.
                        w.WorkerPool <- w.JobChannel
                        select {
                        case job := <-w.JobChannel:
                        // we have received a work request.
                                if err, payLoad := job.JobProcessor(job.JobPayload); err == nil {
                                        if job.NextParallelJobs != nil {
                                                for _, nextJob := range job.NextParallelJobs {
                                                        nextJob.JobPayload = payLoad // set the payload since it might have been modified
                                                        JobQueue <- *nextJob
                                                }
                                        }
                                } else {
                                        log.Println("err not nil: ", err.Error())
                                        SendFailureEmail("custom-sender@email.com", job, err)
                                }
                        case <-w.quitChannel:
                        // we have received a signal to stop
                                return
                        }
                }
        }()
}

func (w Worker) Stop() {
        go func() {
                w.quitChannel <- true
        }()
}
