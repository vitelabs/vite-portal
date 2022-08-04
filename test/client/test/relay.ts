import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeEntity } from "../src/types"
import { CommonUtil } from "../src/utils"

export function testRelay(common: TestCommon) {
  describe("testRelay", () => {
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
      const result = await common.provider.request(method)
      expect(result).to.not.be.undefined
      // check if all mock nodes received a request
      await CommonUtil.sleep(100)
      expect(common.defaultMockNode.requests.length).to.be.equal(1)
      expect(common.timeoutMockNode.requests.length).to.be.equal(1)
    })
  })
};