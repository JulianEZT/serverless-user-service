# Async Worker (Lambda B) — Placeholder

This directory is reserved for **Lambda B**: an async worker that will consume messages from the SQS queue (EVENTS_QUEUE_URL) for background processing.

**Status:** Not implemented yet. It will be implemented in a future iteration.

The event contract (e.g. `user.created` envelope and payload) is defined in `pkg/events` and will be shared by both Lambda A (publisher) and Lambda B (consumer).
