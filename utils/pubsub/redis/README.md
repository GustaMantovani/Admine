## Redis Insight — quick note

Redis Insight is a lightweight GUI and toolset for inspecting Redis data, monitoring performance, and running commands. For pub/sub development it helps you visually subscribe to channels, inspect published messages, run test publishes from the built-in CLI, and monitor latency and throughput while exercising your producers and consumers.

How to use in this project:
- Point Redis Insight at the Redis instance used by this repo (see `redis_pubsub/docker-compose.yaml`).
- Use the Pub/Sub or CLI panels to subscribe to the channels your code publishes to and verify message payloads and patterns.
- Publish test messages to validate handlers and watch metrics (latency, ops/sec, memory) while running the project.

Redis Insight documentation and installation methods are available here: https://redis.io/docs/latest/operate/redisinsight/
