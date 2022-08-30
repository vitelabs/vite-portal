import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeEntity } from "../src/types"
import { CommonUtil } from "../src/utils"

export function testRaw(common: TestCommon) {
  describe("testRaw", () => {
    let node: NodeEntity

    before(async function () {
      node = await common.insertNodeAsync(CommonUtil.uuid())
    })

    after(async function () {
      await common.relayer.deleteNode(node.id)
    })

    it('test getSnapshotChainHeight raw', async function () {
      const method = "ledger_getSnapshotChainHeight"
      const batch = [
        common.client.createJsonRpcRequest(method),
        common.client.createJsonRpcRequest(method)
      ]
      const promises: Promise<any>[] = [
        common.client.provider.post(common.nodeHttpUrl, batch),
        common.client.provider.post(common.providerUrl, batch)
      ]
      const results = await Promise.all(promises)

      expect(results.length).to.be.equal(2)
      expect(results[0].data.length).to.be.equal(2)
      expect(results[1].data.length).to.be.equal(2)

      const data0_0 = results[0].data[0]
      const data1_0 = results[1].data[0]

      expect(data0_0.jsonrpc).to.be.equal("2.0")
      expect(Number(data0_0.id)).to.be.greaterThan(0)
      expect(Number(data0_0.result)).to.be.greaterThan(0)

      expect(data0_0.jsonrpc).to.be.equal(data1_0.jsonrpc)
      expect(data0_0.id).to.be.equal(data1_0.id)
      expect(data0_0.result).to.be.equal(data1_0.result)

      const data0_1 = results[0].data[1]
      const data1_1 = results[1].data[1]

      expect(data0_1.jsonrpc).to.be.equal("2.0")
      expect(Number(data0_1.id)).to.be.greaterThan(0)
      expect(Number(data0_1.result)).to.be.greaterThan(0)

      expect(data0_1.jsonrpc).to.be.equal(data1_1.jsonrpc)
      expect(data0_1.id).to.be.equal(data1_1.id)
      expect(data0_1.result).to.be.equal(data1_1.result)
    })
  })
}