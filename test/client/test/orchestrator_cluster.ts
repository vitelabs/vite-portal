import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeCluster } from "../src/cluster"
import { TestConstants } from "../src/constants"
import { CommonUtil } from "../src/utils"

export function testOrchestratorCluster(common: TestCommon) {
  describe("testOrchestratorCluster", () => {
    let cluster: NodeCluster

    before(async function () {
      const nodes = await common.orchestrator.getNodes(TestConstants.SupportedChains.ViteBuidl)
      expect(nodes.result.total).to.be.equal(0)
      cluster = new NodeCluster(30000)
      cluster.url = "http://127.0.0.1:48132"
      await cluster.start()
    })

    after(async function () {
      await cluster.stop()
    })

    it('test local cluster', async function () {
      const chain = TestConstants.SupportedChains.ViteBuidl
      await CommonUtil.expectAsync(async () => {
        const nodes = await common.orchestrator.getNodes(chain)
        return nodes.result.total >= 3
      }, common.timeout)
      const nodes = await common.orchestrator.getNodes(chain)
      const node1 = nodes.result.entries.find(e => e.id === "d7e63ddca1e41d311db9668bb0f6ff549cab9c24de03ecaa1643940a8fdc3937")
      expect(node1).to.not.be.undefined
      expect(node1?.name).to.be.equal("s1")
      expect(node1?.chain).to.be.equal(chain)
      expect(node1?.version).to.not.be.undefined
      expect(node1?.commit).to.not.be.undefined
      expect(node1?.rewardAddress).to.be.equal("vite_xxxxxxxxxxxxxxxxxx")
      expect(node1?.transport).to.not.be.undefined
      expect(node1?.remoteAddress).to.not.be.undefined
      expect(node1?.clientIp).to.not.be.undefined
      expect(node1?.status).to.be.equal(1)
      const node2 = nodes.result.entries.find(e => e.id === "5037da9f811f390bcf4046ae70b3c9fd88d912ccbfacd74a352412a32c6166a8")
      expect(node2).to.not.be.undefined
      expect(node2?.name).to.be.equal("s2")
      expect(node2?.chain).to.be.equal(chain)
      expect(node2?.version).to.not.be.undefined
      expect(node2?.commit).to.not.be.undefined
      expect(node2?.rewardAddress).to.be.undefined
      expect(node2?.transport).to.not.be.undefined
      expect(node2?.remoteAddress).to.not.be.undefined
      expect(node2?.clientIp).to.not.be.undefined
      expect(node2?.status).to.be.equal(1)
      const node3 = nodes.result.entries.find(e => e.id === "1f5dcf96afb50a30f574f78fdac9b5da19d7c392de44e99826e0fca8cc5b83d3")
      expect(node3).to.not.be.undefined
      expect(node3?.name).to.be.equal("s3")
    })
  })
}