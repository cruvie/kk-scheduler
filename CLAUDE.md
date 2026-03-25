# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

kk-scheduler is a job scheduling system that uses cron and gRPC. Users register services (gRPC servers implementing `KKScheduleTriggerServer`), add jobs with cron specs, and the scheduler triggers jobs at their scheduled times via gRPC calls.

## Commands

### Backend (Go)
```bash
# Format code
just fmt

# Run linter
just lint

# Run tests
just test

# Generate protobuf code (requires buf)
just proto-gen

# Check for dead code
just deadcode

# Check for vulnerabilities
just govulncheck

# Run the server
go run ./internal/main/main.go
```

### Frontend (UI)
```bash
cd ui
npm install
npm run dev      # Development server on http://localhost:3000
npm run build    # Production build
```

## Architecture

### Directory Structure
- `internal/main/` - Entry point, gRPC/HTTP server setup
- `internal/schedule/` - Core scheduling logic: `Client` manages cron jobs, `StoreDriver` interface for persistence
- `internal/api_handlers/` - Handler implementations per RPC method (job/, service/)
- `internal/api_impl/` - gRPC server registration and method routing
- `internal/g_config/` - YAML config loading (ports, etcd settings)
- `kk_scheduler/` - Protobuf definitions and generated Go code

### Key Patterns

**API Handler Pattern**: Each RPC has a dedicated handler struct in `internal/api_handlers/`:
```go
type ApiJobEnable struct {
    *kk_grpc.DefaultApi[kk_scheduler.JobEnable_Input]
}
```
Handlers implement `CheckInput()`, `Service()`, and `Handler()` methods. The `kk_grpc.GrpcHandler` function orchestrates the call flow.

**Global Client**: `schedule.GClient` is the singleton managing the cron scheduler and storage. All job/service operations go through this client.

**Store Interface**: `StoreDriver` interface allows pluggable storage backends. Default is Etcd (`store_etcd.go`). To add a new store, implement the interface and update config.

**Trigger Mechanism**: Jobs trigger gRPC calls to registered services via `triggerFunc()`, which dials the service target and calls `Trigger(funcName)`.

### Protobuf
- Proto files in `kk_scheduler/*.proto`
- Generated code: `*_grpc.pb.go` (gRPC), `*.pb.go` (messages)
- Uses buf for generation with `buf.gen.yaml`
- Edition 2023 with `API_OPAQUE` level

## Configuration

`config.toml` at project root:
```toml
GrpcPort = 8666
HttpPort = 8667
WebPort = 8668

[StoreEtcd]
UserName = "root"
Password = "root"
Endpoints = ["http://127.0.0.1:2379"]
```

Environment: Set `KK_Schedule` env var for environment mode (handled by `kk_go_kit`).

## Dependencies

- `gitee.com/cruvie/kk_go_kit` - Shared library for gRPC utilities, logging, server management. Local copy exists in `kk_go_kit/` directory.
- `github.com/robfig/cron/v3` - Cron scheduling engine
- `go.etcd.io/etcd/client/v3` - Etcd client for job/service persistence