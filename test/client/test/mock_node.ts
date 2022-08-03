import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"

export function testMockNodes(common: TestCommon) {
  it('test mock node send', async function () {
    const url = common.defaultMockNode.url()
    const result = await common.client.send(url, "test_method")
    expect(result.status).to.be.equal(200)
    expect(result.data).to.be.equal("This is a test response!")
  })
}