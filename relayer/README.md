# portal-relayer

Relayer is the core component of VitePortal and responsible for relaying data requests and responses to and from Vite full nodes.

## Getting Started
### Example usage

```
Usage:
  vite-portal-relayer [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Starts vite-portal-relayer daemon
  version     Get current version

Flags:
      --debug   sets log level to debug
  -h, --help    help for vite-portal-relayer

Use "vite-portal-relayer [command] --help" for more information about a command.
```

### Build

All build options will be located in `build/cmd/relayer` after running the following command:

`make all`

### Run

`go run cmd/relayer/main.go`

or

`go run cmd/relayer/main.go start`

The latter command makes use of `relayer_config.json` in the current directory or creates a new file. If the schema version does not match, the existing configuration file will be backed up and replaced with default values.

A description of the different configuration options can be found [here](./internal/types/config.go).

### Test

By running the following command all unit tests will be executed:

`go test ./...`

### Debug

1. Modify args in `.vscode/launch.json`
2. Set breakpoint(s)
3. Open `cmd/relayer/main.go`
4. Press F5

Note: Consider deleting `relayer_config.json` and `logs` in `cmd/relayer` before debugging.

## Docker

### Build image

See root README.md

### Run image

Before running the image you can create and modify the configuration file to be used by the relayer: `$HOME/.relayer/relayer_config.json`

A description of the different configuration options can be found [here](./internal/types/config.go).

```
docker run -v $HOME/.relayer/:/var/relayer/ -p 56331:56331 -p 56332:56332 -p 56333:56333 -p 56334:56334 --name portal-relayer --detach vitelabs/portal-relayer:test start --config /var/relayer/relayer_config.json
```

### Inspect container

```
docker exec -it portal-relayer /bin/bash
```

### Stop/remove container

```
docker rm $(docker stop $(docker ps -a -q --filter ancestor=vitelabs/portal-relayer:test --format="{{.ID}}")) || docker container prune --force
```

# API

* [Get list of nodes](#get_nodes)
* [Create or update a node](#put_node)
* [Get a node by identifier](#get_node)
* [Delete a node](#delete_node)
* [Get list of chains](#get_chains)
* [Relay request](#post_relay) ![](https://img.shields.io/static/v1?label=&message=important&color=yellow)

## Get list of nodes <a name="get_nodes"></a>

Those nodes are managed by the orchestrator (TODO) and used to serve relays.

### Request

```http
GET /api/v1/db/nodes
```

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `chain` | `string` | **Required**. The identifier of the chain |
| `offset` | `number` | The pagination offset |
| `limit` | `number` | The pagination limit |

    curl -i -X GET http://localhost:56333/api/v1/db/nodes?chain=vite_testnet

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 08:10:50 GMT
    Content-Length: 169

    {
      "entries":[
          {
            "id":"n1",
            "chain":"vite_testnet",
            "rpcHttpUrl":"https://buidl.vite.net/gvite",
            "rpcWsUrl":"wss://buidl.vite.net/gvite/ws"
          }
      ],
      "limit":1000,
      "offset":0,
      "total":1
    }

## Create or update a node <a name="put_node"></a>

TODO: add authorization to restrict access to orchestrator

### Request

```http
PUT /api/v1/db/nodes
```

    curl -i -X PUT http://localhost:56333/api/v1/db/nodes \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
        "id": "n1",
        "chain": "vite_testnet",
        "rpcHttpUrl": "https://buidl.vite.net/gvite",
        "rpcWsUrl": "wss://buidl.vite.net/gvite/ws"
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 08:55:01 GMT
    Content-Length: 4

    null

## Get a node by identifier <a name="get_node"></a>

### Request

```http
GET /api/v1/db/nodes/{id}
```

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `id` | `string` | **Required**. The unique identifier of the node |

    curl -i -X GET http://localhost:56333/api/v1/db/nodes/n1

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 09:00:17 GMT
    Content-Length: 121

    {
      "id":"n1",
      "chain":"vite_testnet",
      "rpcHttpUrl":"https://buidl.vite.net/gvite",
      "rpcWsUrl":"wss://buidl.vite.net/gvite/ws"
    }

## Delete a node <a name="delete_node"></a>

TODO: add authorization to restrict access to orchestrator

### Request

```http
DELETE /api/v1/db/nodes/{id}
```

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `id` | `string` | **Required**. The unique identifier of the node |

    curl -i -X DELETE http://localhost:56333/api/v1/db/nodes/n1

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 08:58:01 GMT
    Content-Length: 4

    null

## Get list of chains <a name="get_chains"></a>

### Request

    curl -i -X POST http://localhost:56331/ -d '{"jsonrpc": "2.0", "id": 1, "method": "db_getChains", "params": []}'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 09:02:26 GMT
    Content-Length: 16

    ["vite_testnet"]

## Relay request <a name="post_relay"></a> ![](https://img.shields.io/static/v1?label=&message=important&color=yellow)

The load balancer should forward all incoming HTTP requests to this endpoint given that the relayer can be deployed multiple times and thus scales horizontally. It is assumed that the load balancer routes traffic to relayers based on a pre-defined routing algorithm (e.g. round robin).

### Request

```http
POST /api/v1/client/relay
```

    curl -i -X POST http://localhost:56333/api/v1/client/relay \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
      "jsonrpc": "2.0",
      "id": 1,
      "method": "ledger_getSnapshotChainHeight",
      "params": null
    }'

### Response

    HTTP/1.1 200 OK
    Access-Control-Allow-Methods: POST
    Access-Control-Allow-Origin: *
    Content-Type: application/json; charset=UTF-8
    Date: Fri, 12 Aug 2022 09:52:44 GMT
    Content-Length: 44

    {"id":1,"jsonrpc":"2.0","result":"22293675"}