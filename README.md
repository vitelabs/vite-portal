# VitePortal

VitePortal is a scaling solution to help process the increasing amount of Remote Procedure Calls (RPCs). This is achieved by introducing a load balancer responsible for spawning relayers as needed. A relayer is a standalone application which forwards every RPC request to multiple full nodes and handles the responses. By determining the majority result (consensus) it is possible to reward honest or punish malicious full nodes and thus incentivize them to partake in the process.

<h1 align="center">
	<img src="assets/images/overview.jpg" alt="VitePortal overview">
</h1>

From a user perspective the basic flow is:

1. Send request to `portal.vite.net` (e.g. AWS Elastic Load Balancer)
2. AWS ELB forwards the request to one of the available relayers
3. Relayer forwards the request to multiple, randomly selected full nodes
4. Relayer returns the fastest response to the user via AWS ELB
5. Relayer collects all responses and pushes a summary to e.g. Apache Kafka for further processing by the orchestrator and worker

From a full node perspective the basic flow is:

1. Full node establishes a WebSocket connection with the orchestrator based on the `DashboardTargetURL` configured in `node_config.json`
2. Orchestrator gets the HTTP + WS ports and verifies if the full node is "legitimate"
3. Orchestrator broadcasts the new full node (public ip address + ports) to all relayers
4. Relayer can use the new full node to serve requests of users

This monorepo is organized as follows:

- [relayer](./relayer) - the relayer forwards every RPC request to multiple full nodes and handles the responses
- [orchestrator](./orchestrator) - the orchestrator keeps track of the global state such as participating full nodes
- [test](./test) - contains all integration tests and the logic to start a cluster of nodes locally

## Session handling

<h1 align="center">
	<img src="assets/images/session.jpg" alt="Session handling">
</h1>

## Docker compose

Depending on the operation system it might be needed to install `gettext` for envsubst first.

```
apt-get update && apt-get install gettext
```

### Build

```
docker-compose build
```

### Start

```
docker-compose up -d
```

### Stop

```
docker-compose down
```

### Inspect r1

```
docker exec -it vite-portal_r1_1 /bin/bash
```

## Experimental deployment

1. docker-compose build
2. docker-compose up -d
3. [Insert node(s)](./relayer#put_node) with the curl command
4. Test [relay request](./relayer#post_relay) with the curl command
5. Setup test AWS Load Balancer which serves requests from e.g. https://portal-buidl.vite.net
6. Point AWS Load Balancer to the [relay request](./relayer#post_relay) endpoint

Step 3 is optional. All nodes connect to the orchestrator automatically and the orchestrator broadcasts newly added or updated nodes to all relayers.