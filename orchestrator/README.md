# portal-orchestrator

Orchestrator is the heart of VitePortal and responsible for managing relayers and collecting uptime data from Vite full nodes.

## Getting Started
### Example usage

```
Usage:
  vite-portal-orchestrator [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Starts vite-portal-orchestrator daemon
  version     Get current version

Flags:
      --debug   sets log level to debug
  -h, --help    help for vite-portal-orchestrator

Use "vite-portal-orchestrator [command] --help" for more information about a command.
```

### Build

All build options will be located in `build/cmd/orchestrator` after running the following command:

`make all`

### Run

`go run cmd/orchestrator/main.go`

or

`go run cmd/orchestrator/main.go start`

The latter command makes use of `orchestrator_config.json` in the current directory or creates a new file. If the schema version does not match, the existing configuration file will be backed up and replaced with default values.

A description of the different configuration options can be found [here](./internal/types/config.go).

### Test

By running the following command all unit tests will be executed:

`go test ./...`

### Debug

1. Modify args in `.vscode/launch.json`
2. Set breakpoint(s)
3. Open `cmd/orchestrator/main.go`
4. Press F5

Note: Consider deleting `orchestrator_config.json` and `logs` in `cmd/orchestrator` before debugging.

#### Install websocat

Download the latest `websocat` executable file from the releases page in GitHub repository:

```
sudo wget -qO /usr/local/bin/websocat https://github.com/vi/websocat/releases/latest/download/websocat.aarch64-unknown-linux-musl
```

Set execute permission:

```
sudo chmod a+x /usr/local/bin/websocat
```

Now `websocat` will be available for all users as a system-wide command:

```
websocat ws://localhost:57331/ -E
```

Example response:

```
{"jsonrpc": "2.0", "id": 1, "result": {"id": "1234", "version": 0, "netId": 1}}
```

Create `Bearer` token with `relayer/internal/orchestrator/client/client_test.go`:

```
websocat ws://localhost:57332/ -E -H='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjQ2MTQ2MzksImlzcyI6InZpdGUtcG9ydGFsLXJlbGF5ZXIiLCJzdWIiOiJ0ZXN0MTIzNCJ9.e3dbqQ9RG656Pk4UaKL1IgIVi9IFqk05u_9orBvx1AA'
```

Example response:

```
{"jsonrpc": "2.0", "id": 1, "result": {"id": "test1234", "name": "vite-portal-relayer"}}
```

Example requests:

```
{"jsonrpc": "2.0", "id": 2, "method": "core_getAppInfo", "params": []}
{"jsonrpc": "2.0", "id": 2, "method": "admin_getRelayers", "params": [0, 0]}
```

## Docker

### Build image

```
docker build -f orchestrator.Dockerfile --tag vitelabs/portal-orchestrator:test .
```

### Run image

Before running the image you can create and modify the configuration file to be used by the orchestrator: `$HOME/.orchestrator/orchestrator_config.json`

A description of the different configuration options can be found [here](./internal/types/config.go).

```
docker run -v $HOME/.orchestrator/:/var/orchestrator/ -p 57331:57331 -p 57332:57332 --name portal-orchestrator --detach vitelabs/portal-orchestrator:test start --config /var/orchestrator/orchestrator_config.json
```

### Inspect container

```
docker exec -it portal-orchestrator /bin/bash
```

### Stop/remove container

```
docker rm $(docker stop $(docker ps -a -q --filter ancestor=vitelabs/portal-orchestrator:test --format="{{.ID}}")) || docker container prune --force
```

# API

* [Get version](#get_version)
* [Get list of nodes](#get_nodes)
* [Get list of relayers](#get_relayers)

## Get version <a name="get_version"></a>

### Request

    curl -i -X POST http://localhost:57331/ \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
      "jsonrpc": "2.0",
      "id": 1,
      "method": "core_getAppInfo",
      "params": null
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Vary: Origin
    Date: Thu, 25 Aug 2022 14:21:14 GMT
    Content-Length: 41

    {"jsonrpc":"2.0","id":1,"result":{"id":"ccb6cf52-5a74-4846-a7f5-a579f4cb49c6","version":"v0.0.1-alpha.6","name":"vite-portal-orchestrator"}}

## Get list of nodes <a name="get_nodes"></a>

Nodes listed here have established a connection with the orchestrator automatically based on the `DashboardTargetURL` configured in `node_config.json`

### Request

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `chain` | `string` | **Required**. The identifier of the chain |
| `offset` | `number` | The pagination offset |
| `limit` | `number` | The pagination limit |

    curl -i -X POST http://localhost:57332/ \
    -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjQ2MTQ2MzksImlzcyI6InZpdGUtcG9ydGFsLXJlbGF5ZXIiLCJzdWIiOiJ0ZXN0MTIzNCJ9.e3dbqQ9RG656Pk4UaKL1IgIVi9IFqk05u_9orBvx1AA' \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
        "jsonrpc": "2.0", 
        "id": 1, 
        "method": "admin_getNodes", 
        "params": ["vite_buidl",0,0]
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Vary: Origin
    Date: Sun, 04 Sep 2022 07:02:04 GMT
    Content-Length: 83

    {
      "jsonrpc":"2.0",
      "id":1,
      "result":{
          "entries":[
          ],
          "limit":1000,
          "offset":0,
          "total":1
      }
    }

## Get list of relayers <a name="get_relayers"></a>

Relayers listed here have established a connection with the orchestrator automatically based on the `orchestratorWsUrl` configured in `relayer_config.json`

### Request

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `offset` | `number` | The pagination offset |
| `limit` | `number` | The pagination limit |

    curl -i -X POST http://localhost:57332/ \
    -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjQ2MTQ2MzksImlzcyI6InZpdGUtcG9ydGFsLXJlbGF5ZXIiLCJzdWIiOiJ0ZXN0MTIzNCJ9.e3dbqQ9RG656Pk4UaKL1IgIVi9IFqk05u_9orBvx1AA' \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
        "jsonrpc": "2.0", 
        "id": 1, 
        "method": "admin_getRelayers", 
        "params": [0,0]
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Vary: Origin
    Date: Sun, 04 Sep 2022 07:02:04 GMT
    Content-Length: 271

    {
      "jsonrpc":"2.0",
      "id":1,
      "result":{
          "entries":[
            {
                "id":"4b75732d-0e10-46f8-964c-4e79e3a88674",
                "version":"v0.0.1-alpha.6",
                "transport":"ws",
                "remoteAddress":"172.20.0.3:46474",
                "httpInfo":{
                  "userAgent":"Go-http-client/1.1",
                  "host":"o1:57332"
                }
            }
          ],
          "limit":1000,
          "offset":0,
          "total":1
      }
    }

## Docker

### Build image

```
docker build -f orchestrator.Dockerfile --tag vitelabs/portal-orchestrator:test .
```