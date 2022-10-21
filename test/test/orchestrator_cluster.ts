import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeCluster } from "../src/cluster"
import { TestConstants } from "../src/constants"
import { NodeExtendedEntity } from "../src/types"
import { CommonUtil } from "../src/utils"

export function testOrchestratorCluster(common: TestCommon) {
  describe("testOrchestratorCluster", () => {
    let cluster: NodeCluster
    let nodeId1 = "d7e63ddca1e41d311db9668bb0f6ff549cab9c24de03ecaa1643940a8fdc3937"
    let nodeId2 = "5037da9f811f390bcf4046ae70b3c9fd88d912ccbfacd74a352412a32c6166a8"
    let nodeId3 = "1f5dcf96afb50a30f574f78fdac9b5da19d7c392de44e99826e0fca8cc5b83d3"

    before(async function () {
      const nodes = await common.orchestrator.getNodes(TestConstants.SupportedChains.ViteBuidl)
      expect(nodes.result.total).to.be.equal(0)
      cluster = new NodeCluster(30000)
      cluster.orchestrator = common.orchestrator
      await cluster.start()
    })

    after(async function () {
      await cluster.stop()
    })

    it('test local cluster', async function () {
      const start = Date.now()
      const chain = TestConstants.SupportedChains.ViteBuidl
      let nodes = await common.orchestrator.getNodes(chain)
      const beforeUpdate = (node: NodeExtendedEntity | undefined, name: string, rewardAddress?: string) => {
        expect(node).to.not.be.undefined
        expect(node?.name).to.be.equal(name)
        expect(node?.chain).to.be.equal(chain)
        expect(node?.version).to.not.be.undefined
        expect(node?.commit).to.not.be.undefined
        if (CommonUtil.isNullOrWhitespace(rewardAddress)) {
          expect(node?.rewardAddress).to.be.undefined
        } else {
          expect(node?.rewardAddress).to.be.equal(rewardAddress)
        }
        expect(node?.transport).to.not.be.undefined
        expect(node?.remoteAddress).to.not.be.undefined
        expect(node?.clientIp).to.not.be.undefined
        expect(node?.status).to.be.equal(1)
        expect(node?.lastBlock.hash).to.be.undefined
        expect(node?.lastBlock.height).to.be.undefined
        expect(node?.lastBlock.time).to.be.undefined
        expect(node?.lastUpdate).to.be.equal("0")
        expect(node?.delayTime).to.be.equal("0")
      }
      let node1 = nodes.result.entries.find(e => e.id === nodeId1)
      beforeUpdate(node1, "s1", "vite_xxxxxxxxxxxxxxxxxx")
      let node2 = nodes.result.entries.find(e => e.id === nodeId2)
      beforeUpdate(node2, "s2")
      let node3 = nodes.result.entries.find(e => e.id === nodeId3)
      beforeUpdate(node3, "s3")

      // update node status
      await common.orchestrator.updateNodeStatus()
      nodes = await common.orchestrator.getNodes(chain)
      const afterUpdate = (node: NodeExtendedEntity | undefined, name: string) => {
        expect(node).to.not.be.undefined
        expect(node?.name).to.be.equal(name)
        expect(node?.status).to.be.equal(1)
        expect(node?.lastBlock.hash).to.not.be.undefined
        expect(node?.lastBlock.height).to.be.greaterThan(0)
        expect(node?.lastBlock.time).to.be.greaterThan(0)
        expect(parseInt(node?.lastUpdate ?? "")).to.be.greaterThanOrEqual(start)
        expect(parseInt(node?.delayTime ?? "")).to.be.greaterThanOrEqual(0)
      }
      node1 = nodes.result.entries.find(e => e.id === nodeId1)
      afterUpdate(node1, "s1")
      node2 = nodes.result.entries.find(e => e.id === nodeId2)
      afterUpdate(node2, "s2")
      node3 = nodes.result.entries.find(e => e.id === nodeId3)
      afterUpdate(node3, "s3")
    })
  })
}