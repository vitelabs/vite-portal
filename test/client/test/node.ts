import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { RpcWsClient } from "../src/client"

export function testNode(common: TestCommon) {
  describe("testNode", () => {
    it('test connect', async function () {
      const client = new RpcWsClient(common.timeout, "ws://127.0.0.1:57331")
      client.ws.on('open', function open() {
        console.log('connected');
        client.ws.send(Date.now());
      });

      client.ws.on('close', function close(code: number, reason: Buffer) {
        console.log(`disconnected: ${code} ${reason}`);
      });

      client.ws.on('message', function message(data: Buffer) {
        console.log(JSON.parse(data.toString()))
        //console.log(`Round-trip time: ${Date.now() - data} ms`);
        // setTimeout(function timeout() {
        //   client.ws.send(Date.now());
        // }, 500);
      });
      expect(true).to.be.true
    })
  })
}