import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestContants } from "./constants"
import { RpcWsClient } from "../src/client"
import { JsonRpcRequest, JsonRpcResponse } from "../src/types"
import { CommonUtil } from "../src/utils"
import { VuilderNode } from "../src/node"

export function testOrchestratorNode(common: TestCommon) {
  describe("testOrchestratorNode", () => {
    it('test local node', async function () {
      const node = new VuilderNode()
      await node.start()
    })

    xit('test local node', async function () {
      // TODO: try starting a node with `npx vuilder node --config <config.json>`
      // Set "Single": false in the config otherwise net_nodeInfo returns mock data (invalid netId, node id, etc.)
      await CommonUtil.expectAsync(async () => {
        const nodes = await common.orchestrator.getNodes(TestContants.SupportedChains.ViteBuidl)
        return nodes.total === 1
      }, common.timeout)
      const nodes = await common.orchestrator.getNodes(TestContants.SupportedChains.ViteBuidl)
      const node = nodes.entries[0]
      expect(node.id).to.not.be.empty
    })

    it('test connect/disconnect', async function () {
      let connected = false
      let requests: JsonRpcRequest[] = []
      let errors: JsonRpcResponse<any>[] = []

      const client = new RpcWsClient(common.timeout, "ws://127.0.0.1:57331/ws/gvite/1@0000000000000000000000000000000000000000000000000000000000000000")
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
        } else if (message.error) {
          errors.push(message)
        }
      });

      await CommonUtil.expectAsync(async () => connected === true, common.timeout)
      await CommonUtil.expectAsync(async () => requests.length === 1, common.timeout)
      expect(requests[0].method).to.be.equal("net_nodeInfo")
      await CommonUtil.expectAsync(async () => connected === false, 6000)
      expect(errors.length).to.be.equal(1)
      expect(errors[0].error?.code).to.be.equal(-32000)
      expect(errors[0].error?.message).to.be.equal("failed to call 'net_nodeInfo'")
    })
  })
}