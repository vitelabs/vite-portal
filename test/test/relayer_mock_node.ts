import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { DefaultMockNode } from "../src/mock_node"
import { CommonUtil } from "../src/utils"

export function testRelayerMockNodes(common: TestCommon) {
  it('test DefaultMockNode client.send', async function () {
    const url = common.defaultMockNode.entity.rpcHttpUrl
    const response = await common.client.send(url, "test_method")
    expect(response.status).to.be.equal(200)
    expect(response.data.result).to.be.equal(DefaultMockNode.DEFAULT_RESPONSE.result)
  })

  it('test DefaultMockNode provider.request', async function () {
    const url = common.defaultMockNode.entity.rpcHttpUrl
    const node = common.createRandomNode("mockchain1")
    node.rpcHttpUrl = url
    const putResponse = await common.relayer.putNode(node)
    expect(putResponse.error).to.be.undefined
    const response = await common.provider.request("test_method")
    expect(response).to.be.equal(DefaultMockNode.DEFAULT_RESPONSE.result)
    const deleteResponse = await common.relayer.deleteNode(node.id)
    expect(deleteResponse.error).to.be.undefined
  })

  it('test TimeoutMockNode client.send', async function () {
    const url = common.timeoutMockNode.entity.rpcHttpUrl
    await CommonUtil.expectThrowsAsync(() => common.client.send(url, "test_method"), "timeout of 2100ms exceeded")
  })

  it('test TimeoutMockNode provider.request', async function () {
    const url = common.timeoutMockNode.entity.rpcHttpUrl
    const node = common.createRandomNode("mockchain2")
    node.rpcHttpUrl = url
    const putResponse = await common.relayer.putNode(node)
    expect(putResponse.error).to.be.undefined
    const response = await CommonUtil.expectThrowsAsync(() => common.provider.request("test_method"))
    expect(response.error).to.be.equal("error executing the http request: relay timed out")
    const deleteResponse = await common.relayer.deleteNode(node.id)
    expect(deleteResponse.error).to.be.undefined
  })
}