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
      common.client.send(common.nodeUrl, method),
      common.client.send(common.providerUrl, method)
    ]
    const result = await Promise.all(promises)
    console.log("original:", result[0], "relayed:", result[1])
    expect(result[0].result).to.be.equal(result[1].result)
    const height = await common.provider.request(method)
    expect(Number(height)).to.be.greaterThan(0)
    expect(Number(height)).to.be.greaterThanOrEqual(Number(result[0].result))
    expect(Number(height)).to.be.greaterThanOrEqual(Number(result[1].result))
  })
};