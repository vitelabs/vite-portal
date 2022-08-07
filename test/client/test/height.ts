import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeEntity } from "../src/types"
import { CommonUtil } from "../src/utils"

export function testHeight(common: TestCommon) {
  describe("testHeight", () => {
    let node: NodeEntity

    before(async function () {
      node = await common.insertNodeAsync(CommonUtil.uuid())
    })

    after(async function () {
      await common.relayer.deleteNode(node.id)
    })

    it('test unsupported method', async function () {
      const method = "test_method_1234"
      const result = await CommonUtil.expectThrowsAsync(() => common.provider.request(method))
      expect(result.error.code).to.be.equal(-32601)
      expect(result.error.message).to.be.equal("The method test_method_1234_ does not exist/is not available")
    })

    it('test getSnapshotChainHeight', async function () {
      const method = "ledger_getSnapshotChainHeight"
      const promises: Promise<any>[] = [
        common.client.send(common.nodeHttpUrl, method),
        common.client.send(common.providerUrl, method)
      ]
      const results = await Promise.all(promises)
      console.log("original:", results[0].data, "relayed:", results[1].data)
      expect(results[0].data.result).to.be.equal(results[1].data.result)
      const height = await common.provider.request(method)
      expect(Number(height)).to.be.greaterThan(0)
      expect(Number(height)).to.be.greaterThanOrEqual(Number(results[0].data.result))
      expect(Number(height)).to.be.greaterThanOrEqual(Number(results[1].data.result))
    })

    it('test getSnapshotChainHeight batch', async function () {
      const method = "ledger_getSnapshotChainHeight"
      const expected = await common.client.send(common.nodeHttpUrl, method)
      expect(Number(expected.data.result)).to.be.greaterThan(0)
      const batch = await common.provider.batch([{
        methodName: method,
        params: []
      }, {
        methodName: method,
        params: []
      }])
      expect(batch.length).to.be.equal(2)
      expect(Number(batch[0].result)).to.be.greaterThanOrEqual(Number(expected.data.result))
      expect(batch[0].error).to.be.null
      expect(Number(batch[1].result)).to.be.greaterThanOrEqual(Number(expected.data.result))
      expect(batch[1].error).to.be.null
    })
  })
};