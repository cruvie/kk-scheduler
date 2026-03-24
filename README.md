# kk-scheduler

A job scheduling system based on cron and grpc

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/cruvie/kk-schedule)

# Screenshot

![service-list](https://github.com/cruvie/kk-scheduler/blob/main/readme/service-list.png?raw=true)

![job-list](https://github.com/cruvie/kk-scheduler/blob/main/readme/job-list.png?raw=true)

# System Design


```mermaid
graph TB
  U[SDK RPCClient]
  W[Web UI]

  subgraph kk-scheduler
    C[API Handlers]
    G[Global Client]
    F[Cron Scheduler]
    C -- Update Job --> F
  end

  subgraph App
    A[kk-scheduler Client server]
    A2[User Jobs]
    A -- Trigger --> A2
  end

  subgraph Storage Layer
    J[Store Interface]
    I[Default Etcd Store]
  end

  U -- Get/Put/Enable/Trigger Job/Service --> C
  W -- Get/Put/Enable/Trigger Job/Service --> C
  C --> G
  F --> G
  G -- RPC Trigger --> A
  G <-- Get/Put --> J
  J --> I
```


# Deploy

## Docker

[docker-compose](https://github.com/cruvie/kk-scheduler/tree/main/deploy-docker)

visit http://localhost:8668

# Usage

- install

```shell
go get github.com/cruvie/server@latest
```

- Run a grpc server that implemented `kk_schedule.UnimplementedKKScheduleTriggerServer`
  see [client_server_test.go](https://github.com/cruvie/kk-scheduler/blob/main/server/internal/schedule_test/client_server_test.go)
- Put a service and job into kk-schedule and enable job
  see [readme_test.go](https://github.com/cruvie/kk-scheduler/blob/main/server/internal/schedule_test/readme_test.go)

# Contribute

## provide more test case

any test case PR is welcome

## support move storge engine

kk-schedule use Etcd as default storage engine, but any storage engine
implement [StoreDriver](https://github.com/cruvie/kk-scheduler/blob/main/server/internal/schedule/store.go)
can be used

Step1 create a `store_xxxx.go`
like [StoreDriver](https://github.com/cruvie/kk-scheduler/blob/main/server/internal/schedule/store_etcd.go)

Step2 test it

Step3 add config filed
in [config.go](https://github.com/cruvie/kk-scheduler/blob/main/server/internal/g_config/config.go)
and [config.yml](https://github.com/cruvie/kk-scheduler/blob/main/server/config.yml)

## improve readme doc and code comment

## Web UI improve