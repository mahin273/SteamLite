# StreamLite

A production-grade Netflix-style video streaming platform built with Go microservices, gRPC, Redpanda, PostgreSQL, Redis, SeaweedFS, and Docker.

## Architecture Overview
StreamLite comprises 7 autonomous microservices behind a Traefik API Gateway and an Nginx media edge cache:
- `auth-svc`: Authentication, Argon2id, EdDSA JWT with JWKS rotation, sessions.
- `catalog-svc`: Title and video metadata management, full-text search.
- `upload-svc`: Direct-to-storage resumable S3 multipart upload manager.
- `transcoder`: FFmpeg adaptive bitrate ladder workers with Postgres SKIP LOCKED job queue.
- `streaming-svc`: Playback authorization, concurrent stream limits, and signed URLs.
- `watch-history-svc`: Watch progress buffer, continue-watching row, and ratings.
- `recommendation-svc`: Event-driven trending and collaborative filtering recommendations.

## Quickstart
```bash
make test
```
