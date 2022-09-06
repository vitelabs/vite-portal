import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { JsonRpcRequest } from "../src/types"
import { RpcWsClient } from "../src/client"
import { CommonUtil } from "../src/utils"

export function testOrchestratorNode(common: TestCommon) {
  describe("testOrchestratorNode", () => {
    it('test connect/disconnect', async function () {
      let connected = false
      let requests: JsonRpcRequest[] = []
      const client = new RpcWsClient(common.timeout, "ws://127.0.0.1:57331")
      client.ws.on('open', function open() {
        console.log('connected');
        connected = true
      });

      client.ws.on('close', function close(code: number, reason: Buffer) {
        console.log(`disconnected: ${code} ${reason}`);
        connected = false
      });

      client.ws.on('message', function message(data: Buffer) {
        const message = JSON.parse(data.toString())
        console.log(message)
        if (message.method) {
          requests.push(message)
        }
      });
      await CommonUtil.expectAsync(() => connected === true, common.timeout)
      await CommonUtil.expectAsync(() => requests.length === 1, common.timeout)
      expect(requests[0].method).to.be.equal("net_nodeInfo")
      await CommonUtil.expectAsync(() => connected === false, 6000)
    })
  })
}