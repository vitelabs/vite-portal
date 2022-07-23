import { it } from "mocha"
import { expect } from "chai"
import { NodeEntity } from "../src/relayer"
import { TestCommon } from "./common"

export function testNodes(common: TestCommon) {
  it('test insert invalid node', async function () {
    const node = {}
    const result = await common.relayer.putNode(node as NodeEntity)
    expect(result.status).to.be.equal(400)
    expect(result.data.code).to.be.equal(400)
    expect(result.data.message).to.be.equal("node is invalid")
  })

  it('test get nonexistent node', async function () {
    const getNode1 = await common.relayer.getNode("1234")
    expect(getNode1.status).to.be.equal(404)
    expect(getNode1.data.code).to.be.equal(404)
    expect(getNode1.data.message).to.be.equal("node does not exist")
  })

  it('test insert node', async function () {
    const nodesBefore = await common.relayer.getNodes()
    expect(nodesBefore.total).to.be.greaterThanOrEqual(0)

    const nodeId = "1234"
    const getNodeBefore = await common.relayer.getNode(nodeId)
    expect(getNodeBefore.status).to.be.equal(404)

    const node = createNode(nodeId, "chain1", "0.0.0.0")
    const putResult = await common.relayer.putNode(node)
    expect(putResult.status).to.be.equal(200)

    const nodesAfter = await common.relayer.getNodes()
    expect(nodesAfter.total).to.be.greaterThan(nodesBefore.total)

    const getNodeAfter = await common.relayer.getNode(nodeId)
    expect(getNodeAfter.status).to.be.equal(200)
    const nodeAfter: NodeEntity = getNodeAfter.data
    expect(node.id).to.be.equal(nodeAfter.id)
    expect(node.chain).to.be.equal(nodeAfter.chain)
    expect(node.ipAddress).to.be.equal(nodeAfter.ipAddress)
    expect(node.rewardAddress).to.be.equal(nodeAfter.rewardAddress)
  })
};

const createNode = (id: string, chain: string, ipAddress: string): NodeEntity => {
  return {
    id,
    chain,
    ipAddress
  }
}