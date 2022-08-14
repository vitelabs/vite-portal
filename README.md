# VitePortal

VitePortal is a scaling solution to help process the increasing amount of Remote Procedure Calls (RPCs). This is achieved by introducing a load balancer responsible for spawning relayers as needed. A relayer is a standalone application which forwards every RPC request to multiple full nodes and handles the responses. By determining the majority result (consensus) it is possible to reward honest or punish malicious full nodes and thus incentivize them to partake in the process.

<h1 align="center">
	<img src="assets/images/overview.jpg" alt="VitePortal overview">
</h1>

This monorepo is organized as follows:

- [lb](./lb) - the load balancer accepts incoming traffic from clients and routes requests to its registered relayers
- [relayer](./relayer) - the relayer forwards every RPC request to multiple full nodes and handles the responses
- [orchestrator](./orchestrator) - the orchestrator keeps track of the global state such as participating full nodes
- [worker](./worker) - the worker is responsible to send out rewards to full nodes on a daily basis

## Session handling

<h1 align="center">
	<img src="assets/images/session.jpg" alt="Session handling">
</h1>

## Docker compose

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