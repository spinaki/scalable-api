package main
import (
        "time"
)

type PayLoad struct {
        Id    string
        Data  string
        Time  time.Time
        Email string
        RemoteUrl string
}

type Job struct {
        Name         string
        JobPayload      PayLoad
        JobProcessor func(PayLoad) (error, PayLoad)
        // optional sequential job to run when this job is completed
        NextParallelJobs[] *Job
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// init the service
// this should be invoked when the API service is initialized
func InitJobService() {
        // init the jobqueue
        JobQueue = make(chan Job, 1000)
        maxWorker := 10
        dispatcher := NewDispatcher(maxWorker)
        dispatcher.Run()
}


// create custom jobs when a rest API request comes and starts the process (by enqueing the jobs in the JobQueue
// this should be invoked when new API request comes in
func AddNewJobs(appConfig JobConfig) {
        payLoad := PayLoad{
                Id: appConfig.Id,
                Data:appConfig.Data,
                Email: appConfig.Email,
                Time:appConfig.Time,
        }
        enqueueDemoJobs(&payLoad)
}

// TODO: make this configurable
func enqueueDemoJobs(payLoad *PayLoad) {
        emailNotifier := Job{Name:"SendEmailJob", JobPayload:*payLoad, JobProcessor: SendEmailJob}

        updateDatabaseJob := Job{Name:"UpdateDatabaseJobStatus", JobPayload:*payLoad, JobProcessor: UpdateDatabaseJobStatus}

        workSpaceCleaner := Job{Name:"CleanUpWorkspace", JobPayload:*payLoad, JobProcessor: CleanUpWorkspaceJob}

        cloudUploader := Job{Name:"UploadJob", JobPayload:*payLoad, JobProcessor: S3UploadJob,
                NextParallelJobs: []*Job{&workSpaceCleaner, &updateDatabaseJob, &emailNotifier}}

        heavyComputer := Job{Name:"HeavyComputationJob", JobPayload:*payLoad, JobProcessor: HeavyComputationJob,
                NextParallelJobs: []*Job{&cloudUploader}}

        preProcessor := Job{Name:"PreProcessor", JobPayload:*payLoad,
                JobProcessor: PreprocessingJob, NextParallelJobs: []*Job{&heavyComputer} }

        JobQueue <- preProcessor
}
