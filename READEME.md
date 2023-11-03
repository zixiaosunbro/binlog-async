# Example Async Deal Mysql Binlog
This project provides an example of dealing with MySQL binlog asynchronously. Mysql is widely used as a storage tool, and we often encounter the following problems:
1. Redis is used as cache, and when the data changes, the cache needs to be updated promptly.
2. For complex queries, data is also saved in Elasticsearch(ES). When data changes, needs to be synchronized with ES, and so on.

To address these issues, I frequently initiate a service that acts as a MySQL slave, receiving binlogs from the master MySQL server. This service parses binlog data to update target cache and synchronize data. However, another issue may arise when dealing with high binlog query per second(QPS), which can lead to increased service workload. To mitigate this, I handle the binlogs asynchronously.

In this project, I use [go-mysql-org/go-mysql](https://github.com/go-mysql-org/go-mysql) to receive and parse MySQL binlogs, and [hibiken/asynq](https://github.com/hibiken/asynq) to implement asynchronous actions.


# problem
TODO: ERROR 1236 (HY000): Could not open log file