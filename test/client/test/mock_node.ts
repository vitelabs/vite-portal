import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { DefaultMockNode } from "../src/mock_node"
import { CommonUtil } from "../src/utils"

export function testMockNodes(common: TestCommon) {
  it('test DefaultMockNode client.send', async function () {
    const url = common.defaultMockNode.entity.rpcHttpUrl
    const result = await common.client.send(url, "test_method")
    expect(result.status).to.be.equal(200)
    expect(result.data.result).to.be.equal(DefaultMockNode.DEFAULT_RESPONSE.result)
  })

  it('test DefaultMockNode provider.request', async function () {
    const url = common.defaultMockNode.entity.rpcHttpUrl
    const node = common.createRandomNode("mockchain")
    node.rpcHttpUrl = url
    const putResult = await common.relayer.putNode(node)
    expect(putResult.status).to.be.equal(200)
    const result = await common.provider.request("test_method")
    expect(result).to.be.equal(DefaultMockNode.DEFAULT_RESPONSE.result)
    const deleteResult = await common.relayer.deleteNode(node.id)
    expect(deleteResult.status).to.be.equal(200)
  })

  it('test TimeoutMockNode client.send', async function () {
    const url = common.timeoutMockNode.entity.rpcHttpUrl
    await CommonUtil.expectThrowsAsync(() => common.client.send(url, "test_method"), "timeout of 2100ms exceeded")
  })

  it('test TimeoutMockNode provider.request', async function () {
    const url = common.timeoutMockNode.entity.rpcHttpUrl
    const node = common.createRandomNode("mockchain")
    node.rpcHttpUrl = url
    const putResult = await common.relayer.putNode(node)
    expect(putResult.status).to.be.equal(200)
    const result = await CommonUtil.expectThrowsAsync(() => common.provider.request(url, "test_method"))
    expect(result.error).to.be.equal("error executing the http request: relay timed out")
    const deleteResult = await common.relayer.deleteNode(node.id)
    expect(deleteResult.status).to.be.equal(200)
  })
}