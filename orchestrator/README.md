# portal-orchestrator

## Run

```
./mvnw spring-boot:run
```

```
curl localhost:8080
```

```
curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Host: localhost:8080" \
     --header "Origin: http://localhost:8080" \
     http://localhost:8080/ws/gvite/1234
```