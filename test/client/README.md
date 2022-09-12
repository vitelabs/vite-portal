# Test client

- Run `npm run test`

## Docker

1. docker-compose up
2. docker-compose kill
3. docker-compose down (otherwise volumes are still in use and can't be removed)
4. ./docker_remove_volumes.sh

### Get chain height

    curl -i -X POST http://localhost:48132 \
    -H 'Content-Type: application/json; charset=UTF-8' \
    --data-raw '
    {
      "jsonrpc": "2.0",
      "id": 1,
      "method": "ledger_getSnapshotChainHeight",
      "params": null
    }'