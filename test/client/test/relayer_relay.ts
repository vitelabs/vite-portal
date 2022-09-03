import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestContants } from "./constants"
import { DefaultMockNode } from "../src/mock_node"
import { NodeEntity, NodeResponse, RelayResult } from "../src/types"
import { CommonUtil } from "../src/utils"

export function testRelayerRelay(common: TestCommon) {
  describe("testRelayerRelay", () => {
    let nodes: NodeEntity[]

    before(async function () {
      const chain = common.defaultMockNode.chain
      nodes = [
        common.defaultMockNode.entity,
        common.timeoutMockNode.entity,
        common.createRandomNode(chain),
        common.createRandomNode(chain),
      ]
      for (const node of nodes) {
        await common.relayer.putNode(node)
      }
    })

    after(async function () {
      for (const node of nodes) {
        await common.relayer.deleteNode(node.id)
      }
    })

    it('test getSnapshotChainHeight', async function () {
      const method = "ledger_getSnapshotChainHeight"
      const expected = await common.client.send(common.nodeHttpUrl, method)
      const result = await common.provider.request(method)
      expect(result).to.not.be.undefined
      // check if all mock nodes received a request
      const timeout = TestContants.DefaultRpcNodeTimeout + 200
      await CommonUtil.expectAsync(() => common.defaultMockNode.requests.length == 1, timeout)
      await CommonUtil.expectAsync(() => common.timeoutMockNode.requests.length == 1, timeout)
      const relayResults = common.httpMockCollector.results
      await CommonUtil.expectAsync(() => relayResults.length == 1, timeout)
      const relayResult: RelayResult = relayResults[0]
      expect(relayResult.sessionKey).to.be.equal("cb8b30cecd1857c59530f8bda15fab91")
      expect(relayResult.relay.host).to.be.equal("127.0.0.1:56333")
      expect(relayResult.relay.chain).to.be.equal(TestContants.DefaultChain)
      expect(relayResult.relay.clientIp).to.be.equal(TestContants.DefaultIpAddress)
      expect(relayResult.relay.payload.data).to.contain(method)
      expect(relayResult.relay.payload.method).to.be.equal("POST")
      expect(relayResult.relay.payload.path).to.be.equal("")
      expect(relayResult.relay.payload.headers).to.not.be.undefined
      // check if DeadlineExceeded and Cancelled are set correctly in dispatched relay result
      const response1 = getByNodeId(nodes[0].id, relayResult)
      expect(response1.cancelled).to.be.false
      expect(response1.deadlineExceeded).to.be.false
      expect(response1.error).to.be.empty
      expect(response1.response).to.be.equal(JSON.stringify(DefaultMockNode.DEFAULT_RESPONSE))
      expect(response1.responseTime).to.be.greaterThanOrEqual(0)
      const response2 = getByNodeId(nodes[1].id, relayResult)
      expect(response2.cancelled).to.be.false
      expect(response2.deadlineExceeded).to.be.true
      expect(response2.error).to.be.equal('Post "http://127.0.0.1:23471": context deadline exceeded')
      expect(response2.response).to.be.empty
      expect(response2.responseTime).to.be.greaterThanOrEqual(TestContants.DefaultRpcNodeTimeout)
      const response3 = getByNodeId(nodes[2].id, relayResult)
      expect(response3.cancelled).to.be.false
      expect(response3.deadlineExceeded).to.be.false
      expect(response3.error).to.be.empty
      expect(JSON.parse(response3.response).result).to.be.equal(expected.data.result)
      expect(response3.responseTime).to.be.greaterThanOrEqual(0)
      const response4 = getByNodeId(nodes[3].id, relayResult)
      expect(response4.cancelled).to.be.false
      expect(response4.deadlineExceeded).to.be.false
      expect(response4.error).to.be.empty
      expect(JSON.parse(response4.response).result).to.be.equal(expected.data.result)
      expect(response4.responseTime).to.be.greaterThanOrEqual(0)
    })

    function getByNodeId(nodeId: string, result: RelayResult): NodeResponse {
      for (const response of result.responses) {
        if (response.nodeId === nodeId) {
          return response
        }
      }
      throw new Error(`result does not contain node with id: ${nodeId}`)
    }
  })
}