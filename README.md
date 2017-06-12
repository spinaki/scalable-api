# Scalable Asynchronous API
This is a library for scalable REST API Endpoint for expensive requests -- designed to handle millions of requests. 

Imaging a common use case. You are building an API service that can ingest data (image, video, text, etc) and run expensive jobs on them.
For instance, running face / object detection in images, scene detection in videos, NLP in text. Clearly -- it is non trivial to
scale such an API service. This library provides a set of tools to do that.

Asynchronous queue based library for handling millions of time intensive requests per minute (e.g., network call, data processing).
This library also enables specifying topological relationship between jobs -- similar to a tree.

* Goal was to be able to handle a large amount of POST requests from millions of endpoints.
* The REST API would receive a JSON or binary data -- like images etc, and the goal is to do a bunch of intensive tasks, 
like uploading to a cloud storage provider like S3, running expensive tasks like machine learning / computer vision algorithms.
* The naive approach would be to create a buffered channel and queue up all the jobs pertaining to a request. 
However, in most real world scenarios the number of requests coming in is orders of magnitudes higher than the 
single processor's capacity to handle them.  Hence this might not be useful when there are millions of expensive requests coming in 
 -- the synchronous processor  would still handle one job at a time and eventually the queue will be backed up and throughput will decrease.
* Better approach is to use a common pattern of thread pools from Java -- however utilizing channels.
* In this library we use a bunch of workers which spawn their goroutines and handle jobs using channels.
* The library can also handle topological dependencies between jobs.

### Details
* Dispatcher is responsible for pulling work requests off of the queue, and distributing them to the next available worker. 
* The worker is a buffered channel of channels. The channels that go into this channel, are what the workers use to receive the work request.

##### References
* [Marcio.io Blog](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
* [Writing worker queues, in Go](http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html)