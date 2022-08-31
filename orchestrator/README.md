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
websocat ws://localhost:57332/ -E
```

```
{"jsonrpc": "2.0", "id": 1, "result": "1234"}
{"jsonrpc": "2.0", "id": 2, "method": "core_getAppInfo", "params": []}
{"jsonrpc": "2.0", "id": 2, "method": "admin_getSecret", "params": []}
```

# API

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

    {"jsonrpc":"2.0","id":1,"result":"v0.1"}
