# portal-lb

## HTTP

Map all incoming HTTP requests to the following model and forward to `/api/v1/client/relay`:

```json
{
  "data": "", // Request body to be forwarded to nodes
  "method": "", // e.g. GET, POST, PUT, DELETE, etc.
  "path": "", // REST support
  "headers": [][] // HTTP headers
}
```