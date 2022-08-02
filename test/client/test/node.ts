import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestContants } from "./constants"
import { NodeEntity } from "../src/relayer"
import { CommonUtil } from "../src/utils"

export function testNodes(common: TestCommon) {
  it('test insert invalid node', async function () {
    const node = {}
    const result = await common.relayer.putNode(node as NodeEntity)
    expect(result.status).to.be.equal(400)
    expect(result.data.code).to.be.equal(400)
    expect(result.data.message).to.be.equal("node is invalid")
  })

  it('test get nonexistent node', async function () {
    const result = await common.relayer.getNode(CommonUtil.uuid())
    expect(result.status).to.be.equal(404)
    expect(result.data.code).to.be.equal(404)
    expect(result.data.message).to.be.equal("node does not exist")
  })

  it('test insert and delete node', async function () {
    const chain = CommonUtil.uuid()
    const nodesBefore = await common.relayer.getNodes(chain)
    expect(nodesBefore.total).to.be.equal(0)

    const node = common.createRandomNode(chain)
    const getNodeBefore = await common.relayer.getNode(node.id)
    expect(getNodeBefore.status).to.be.equal(404)

    const putResult = await common.relayer.putNode(node)
    expect(putResult.status).to.be.equal(200)

    const nodesAfter = await common.relayer.getNodes(chain)
    expect(nodesAfter.total).to.be.equal(nodesBefore.total + 1)

    const getNodeAfter = await common.relayer.getNode(node.id)
    expect(getNodeAfter.status).to.be.equal(200)
    const nodeAfter: NodeEntity = getNodeAfter.data
    expect(node.id).to.be.equal(nodeAfter.id)
    expect(node.chain).to.be.equal(nodeAfter.chain)
    expect(node.rpcHttpUrl).to.be.equal(nodeAfter.rpcHttpUrl)
    expect(node.rpcWsUrl).to.be.equal(nodeAfter.rpcWsUrl)

    const deleteResult = await common.relayer.deleteNode(node.id)
    expect(deleteResult.status).to.be.equal(200)

    const nodesAfterDelete = await common.relayer.getNodes(chain)
    expect(nodesAfterDelete.total).to.be.equal(nodesAfter.total - 1)

    const getNodeAfterDelete = await common.relayer.getNode(node.id)
    expect(getNodeAfterDelete.status).to.be.equal(404)
  })

  it('test delete nonexistent node', async function () {
    const result = await common.relayer.deleteNode(CommonUtil.uuid())
    expect(result.status).to.be.equal(200)
    expect(result.data).to.be.null
  })

  it('test get paginated nodes', async function () {
    const chain = CommonUtil.uuid()
    const nodesBefore = await common.relayer.getNodes(chain)
    expect(nodesBefore.total).to.be.equal(0)

    const nodes: NodeEntity[] = []
    for (let index = 0; index < 10; index++) {
      const node = common.createRandomNode(chain)
      nodes.push(node)
      const putResult = await common.relayer.putNode(node)
      expect(putResult.status).to.be.equal(200)
    }

    const nodesAfter = await common.relayer.getNodes(chain)
    expect(nodesAfter.total).to.be.equal(nodes.length)
    expect(nodesAfter.entries.length).to.be.equal(nodes.length)
    expect(nodesAfter.limit).to.be.equal(TestContants.DefaultPageLimit)
    expect(nodesAfter.offset).to.be.equal(0)

    const page1 = await common.relayer.getNodes(chain, 0, 6)
    expect(page1.total).to.be.equal(nodes.length)
    expect(page1.entries.length).to.be.equal(6)
    expect(page1.limit).to.be.equal(6)
    expect(page1.offset).to.be.equal(0)
    expect(page1.entries[0].id).to.equal(nodes[0].id)

    const page2 = await common.relayer.getNodes(chain, page1.entries.length)
    expect(page2.total).to.be.equal(nodes.length)
    expect(page2.entries.length).to.be.equal(4)
    expect(page2.limit).to.be.equal(TestContants.DefaultPageLimit)
    expect(page2.offset).to.be.equal(page1.entries.length)
    expect(page2.entries[0].id).to.equal(nodes[6].id)
  })
};