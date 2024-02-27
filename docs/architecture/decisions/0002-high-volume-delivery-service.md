# 2. High Volume Delivery Service

Date: 2024-02-25

## Status

In Review

## Context

We need to trigger alerts/messages to users using our webhook endpoint. Requirement from the support team is to
monitor a specific directory for new files, parse the user-ids and send messages to each one of them.

## Approaches
This document captures 3 approaches to handle the file processing and delivery to users via webhook api.

### Approach 1 - Synchronous Processing
In this approach, the system will process files sequentially, reading each line and making webhook API calls synchronously for each user ID.

### Approach 2 - Background Job Processing (Preferred)
In this approach, a pool of worker threads is utilized to process files asynchronously. Each worker thread reads a file, extracts user IDs, and makes API calls concurrently, allowing for parallel processing.
We will be leveraging GoCraft library with underlying standalone redis node to enqueue jobs for each user-id. The worker setup can be scaled up independently, 
leading to total decoupling of the file processing and webhook delivery while maintaining a shared logic inside the same codebase. GoCraft offers support for deadset incase the webhook delivery fails even after
all the retry attempts. This way we can make sure no information for user notification is missed.

### Approach 3 - Event Driven Processing
This approach employs an event-driven architecture where file upload events trigger processing tasks. User IDs extracted from files are pushed onto a message queue, which is then consumed by worker processes
responsible for making API calls. This again gives us the flexibility to decouple the file processing with webhook delivery, and lets us handle the failed deliveries by not discarding the message till its processed.
Message queue are a first class citizen in distributed systems, which warrants all the cost that comes with it. They must be deployed and monitored like the rest of your infrastructure; they must be reliable and highly available.
message queueing is a generally considered a superset of background processing. All the benefits we get from background processing can be realized through message queueing as well.

### Comparison

| **Approach**                  | **Pros**                                                                                                                                                                                                                                                                                                     | **Cons**                                                                                                                                                                                                                                                                                                                         |
|-------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Sequential Processing**     | - Simple to implement <br> - Straightforward error handling and flow control due to synchronous processing                                                                                                                                                                                                   | - Not suitable for high-volume processing as it processes files sequentially, tightly coupled with delivery service response <br> - Longer processing times, especially for large files, leading to potential delays in alert delivery <br> - No out of the box way to handle failed message deliveries. Messages might get lost |
| **Background Job Processing** | - Better resource utilization by leveraging concurrency, making it suitable for high-volume processing <br> - Gocraft offers deadset for failed job processing, lets us trigger alert to users for which system might have errored out. <br> - Decouples the file processing with message delivery operation | - Additional infra cost (redis) for managing state for workers                                                                                                                                                                                                                                                                   |
| **Event Driven Processing**   | - Highly scalable and resilient to varying workloads since message queue acts as a distribution system component. Decouples file processing from API calls <br> - Messages are not lost even if the delivery fails. Message queue maintains it in the workflow till the processing is complete               | - Message queues are required when we usually need 2 different system components to communicate. Not necessarily the case here. No need to over complicated system till we dont have a requirement <br> - Cost of maintaining a RabbitMQ cluster is usually higher than a worker architecture                                    |
