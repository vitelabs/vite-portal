import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"

export function testHeight(common: TestCommon) {
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