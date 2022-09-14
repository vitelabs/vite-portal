import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestContants } from "./constants"
import { NodeCluster } from "../src/cluster"
import { CommonUtil } from "../src/utils"

export function testOrchestratorCluster(common: TestCommon) {
  describe("testOrchestratorCluster", () => {
    let cluster: NodeCluster

    before(async function () {
      const nodes = await common.orchestrator.getNodes(TestContants.SupportedChains.ViteBuidl)
      expect(nodes.total).to.be.equal(0)
      cluster = new NodeCluster(30000)
      cluster.url = "http://127.0.0.1:48132"
      await cluster.start()
    })

    after(async function () {
      await cluster.stop()
    })

    it('test local cluster', async function () {
      await CommonUtil.expectAsync(async () => {
        const nodes = await common.orchestrator.getNodes(TestContants.SupportedChains.ViteBuidl)
        return nodes.total === 1
      }, common.timeout)
      const nodes = await common.orchestrator.getNodes(TestContants.SupportedChains.ViteBuidl)
      const node = nodes.entries[0]
      expect(node.id).to.not.be.empty
    })
  })
}