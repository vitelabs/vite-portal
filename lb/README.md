# portal-lb

The load balancer accepts incoming traffic from clients and routes requests to its registered relayers e.g. based on round-robin.

## HTTP

Currently the following mapping is not implemented but in a future release all incoming HTTP requests could be mapped to the following model before being forwarded to a relayer:

```json
{
  "data": "", // Request body to be forwarded to nodes
  "method": "", // e.g. GET, POST, PUT, DELETE, etc.
  "path": "", // REST support
  "headers": [][] // HTTP headers
}
```