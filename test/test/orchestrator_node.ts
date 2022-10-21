import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { RpcWsClient, RpcWsClientWrapper } from "../src/client"
import { TestConstants } from "../src/constants"
import { CommonUtil } from "../src/utils"

export function testOrchestratorNode(common: TestCommon) {
  describe("testOrchestratorNode", () => {
    xit('test local node', async function () {
      // TODO: try starting a node with `npx vuilder node --config <config.json>`
      // Set "Single": false in the config otherwise net_nodeInfo returns mock data (invalid netId, node id, etc.)
      await CommonUtil.expectAsync(async () => {
        const nodes = await common.orchestrator.getNodes(TestConstants.SupportedChains.ViteBuidl)
        return nodes.result.total === 1
      }, common.timeout)
      const nodes = await common.orchestrator.getNodes(TestConstants.SupportedChains.ViteBuidl)
      const node = nodes.result.entries[0]
      expect(node.id).to.not.be.empty
    })

    it('test connect/disconnect', async function () {
      const url = "ws://127.0.0.1:57331/ws/gvite/1@0000000000000000000000000000000000000000000000000000000000000000"
      const client = new RpcWsClient(common.timeout, url, CommonUtil.uuid())
      const wrapper = new RpcWsClientWrapper(client)

      await CommonUtil.expectAsync(async () => wrapper.connected === true, common.timeout)
      await CommonUtil.expectAsync(async () => wrapper.requests.length === 1, common.timeout)
      expect(wrapper.requests[0].method).to.be.equal("net_nodeInfo")
      await CommonUtil.expectAsync(async () => wrapper.connected === false, 6000)
      expect(wrapper.errors.length).to.be.equal(1)
      expect(wrapper.errors[0].error?.code).to.be.equal(-32000)
      expect(wrapper.errors[0].error?.message).to.be.equal("failed to call 'net_nodeInfo'")
    })
  })
}