# scalable-api
This is a library for scalable REST API Endpoint for expensive requests. 
Asynchronous queue based library for handling millions of time intensive requests per minute (e.g., network call, data processing).
Enable Specifying relationship between jobs -- similar to a tree.
* Goal was to be able to handle a large amount of POST requests from millions of endpoints.
* The REST API would receive a JSON / binary data -- like images etc, and the goal is to do a bunch of intensive tasks, 
like uploading to a cloud storage provider like S3, running expensive tasks like machine learning / computer vision algorithms.
*  create a buffered channel where we could queue up some jobs and upload them to S3, and since we could control the maximum 
number of items in our queue and we had plenty of RAM to queue up jobs in memory, we thought it would be okay to just 
buffer jobs in the channel queue.

* This approach didn’t buy us anything, we have traded flawed concurrency with a buffered queue that was simply postponing the problem. Our synchronous processor was only uploading one payload at a time to S3, and since the rate of incoming requests were much larger than the ability of the single processor to upload to S3, our buffered channel was quickly reaching its limit and blocking the request handler ability to queue more items.
We were simply avoiding the problem and started a count-down to the death of our system eventually. Our latency rates kept increasing in a constant rate minutes after we deployed this flawed version.
* Better approach is to use a common pattern of thread pools from Java -- however utilizing channels instead.

### Details
* Dispatcher is responsible for pulling work requests off of the queue, and distributing them to the next available worker. 
* The worker queue is the weirdest part about all of this (if you are not familiar with Go, and maybe, even if you are). 
It is a buffered channel of channels. The channels that go into this channel, are what the workers use to receive the work request.

##### References
* [Marcio.io Blog](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
* [Writing worker queues, in Go](http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html)