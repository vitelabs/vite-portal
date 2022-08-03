import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeEntity } from "../src/relayer"
import { CommonUtil } from "../src/utils"

export function testHeight(common: TestCommon) {
  let node: NodeEntity

  before(async function () {
    node = await common.insertNodeAsync(CommonUtil.uuid())
  })

  after(async function () {
    await common.relayer.deleteNode(node.id)
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
};