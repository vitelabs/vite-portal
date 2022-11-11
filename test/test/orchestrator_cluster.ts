import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { NodeCluster } from "../src/cluster"
import { TestConstants } from "../src/constants"
import { VitePortal } from "../src/portal"
import { NodeExtendedEntity, RelayerConfig } from "../src/types"
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

    it('test node status update', async function () {
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

    it('test relayer nodes', async function () {
      // make sure ports are not used by another process
      const relayerConfig: RelayerConfig = {
        rpcUrl: "http://127.0.0.1:55331",
        rpcAuthUrl: "http://127.0.0.1:55332",
        rpcRelayHttpUrl: "http://127.0.0.1:55333",
        rpcRelayWsUrl: "http://127.0.0.1:55334",
        jwtSecret: TestConstants.DefaultJwtSecret
      }
      const relayer = await VitePortal.startRelayer(relayerConfig, common.timeout)
      const chain1 = TestConstants.SupportedChains.ViteBuidl
      let response1 = await relayer.getNodes(chain1)
      expect(response1.result.limit).to.be.equal(1000)
      expect(response1.result.offset).to.be.equal(0)
      expect(response1.result.total).to.be.equal(3)
      expect(response1.result.entries.length).to.be.equal(3)
      expect(response1.result.entries[0].id).to.not.be.empty
      expect(response1.result.entries[0].chain).to.equal(chain1)
      expect(response1.result.entries[0].rpcHttpUrl).to.satisfy((e: string) => e.startsWith("http://") && e.endsWith("48132"))
      expect(response1.result.entries[0].rpcWsUrl).to.satisfy((e: string) => e.startsWith("ws://") && e.endsWith("41420"))
      const chain2 = TestConstants.SupportedChains.ViteMain
      let response2 = await relayer.getNodes(chain2)
      expect(response2.result.limit).to.be.equal(1000)
      expect(response2.result.offset).to.be.equal(0)
      expect(response2.result.total).to.be.equal(0)
      expect(response2.result.entries.length).to.be.equal(0)
      await relayer.stop()
    })

    it('test node status dispatch', async function () {
      const limit = 1000
      const response1 = await common.orchestrator.getKafkaDefaultMessages(0, limit)
      await common.orchestrator.dispatchNodeStatus()
      await CommonUtil.sleep(2000)
      const response2 = await common.orchestrator.getKafkaDefaultMessages(response1.result?.length ?? 0, limit, 500)
      if (!response2.result || response2.result.length !== 3) {
        console.log("response1", response1)
        console.log("response2", response2)
      }
      const events = response2.result.map(e => JSON.parse(e))
      const event1 = events.find(e => e.nodeName == "s1")
      expect(event1).to.not.be.undefined
      expect(event1.eventId).to.not.be.undefined
      expect(event1.timestamp).to.not.be.undefined
      expect(event1.round).to.not.be.undefined
      expect(event1.ip).to.not.be.undefined
      expect(event1.successTime).to.be.equal(1)
      expect(event1.viteAddress).to.be.equal("vite_xxxxxxxxxxxxxxxxxx")
      expect(event1.chain).to.be.equal("vite_buidl")
      const event2 = events.find(e => e.nodeName == "s2")
      expect(event2).to.not.be.undefined
      expect(event2.eventId).to.not.be.undefined
      expect(event2.timestamp).to.not.be.undefined
      expect(event2.round).to.not.be.undefined
      expect(event2.ip).to.not.be.undefined
      expect(event2.successTime).to.be.equal(1)
      expect(event2.viteAddress).to.be.empty
      expect(event2.chain).to.be.equal("vite_buidl")
      const event3 = events.find(e => e.nodeName == "s3")
      expect(event3).to.not.be.undefined
      expect(event3.eventId).to.not.be.undefined
      expect(event3.timestamp).to.not.be.undefined
      expect(event3.round).to.not.be.undefined
      expect(event3.ip).to.not.be.undefined
      expect(event3.successTime).to.be.equal(1)
      expect(event3.viteAddress).to.be.empty
      expect(event3.chain).to.be.equal("vite_buidl")
    })
  })
}