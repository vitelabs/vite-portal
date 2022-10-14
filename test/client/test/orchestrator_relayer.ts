import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestConstants } from "../src/constants"
import { Relayer } from "../src/relayer"
import { CommonUtil, FileUtil, getLocalFileUtil } from "../src/utils"
import { RelayerConfig } from "../src/types"

export function testOrchestratorRelayer(common: TestCommon) {
  let fileUtil: FileUtil

  before(async function () {
    fileUtil = getLocalFileUtil()
  })

  describe("testOrchestratorRelayer1", () => {
    it('test getPaginated relayers', async function () {
      const relayers = await common.orchestrator.getRelayers()
      expect(relayers.result.total).to.be.equal(1)
      const relayer = relayers.result.entries[0]
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      expect(relayer.version).to.be.equal(expectedVersion.trim())
      expect(relayer.id).to.not.be.empty
      expect(relayer.transport).to.be.equal("ws")
      expect(relayer.remoteAddress).to.not.be.empty
      expect(relayer.httpInfo.version).to.be.undefined
      expect(relayer.httpInfo.userAgent).to.be.equal("Go-http-client/1.1")
      expect(relayer.httpInfo.origin).to.be.undefined
      expect(relayer.httpInfo.host).to.be.equal("127.0.0.1:57332")
    })

    describe("testOrchestratorRelayer1.1", () => {
      let relayer: Relayer

      before(async function () {
        const config: RelayerConfig = {
          rpcUrl: "http://127.0.0.1:56341",
          rpcAuthUrl: "http://127.0.0.1:56342",
          rpcRelayHttpUrl: "http://127.0.0.1:56343",
          rpcRelayWsUrl: "http://127.0.0.1:56344",
          jwtSecret: TestConstants.DefaultJwtSecret
        }
        relayer = new Relayer(config, common.timeout, CommonUtil.uuid())
      })

      after(async function () {
        await relayer.stop()
      })

      it('test spawn/despawn relayer', async function () {
        const relayersBefore = await common.orchestrator.getRelayers()
        expect(relayersBefore.result.total).to.be.equal(1)
        await relayer.start()
        const relayersAfter1 = await common.orchestrator.getRelayers()
        expect(relayersAfter1.result.total).to.be.equal(2)
        expect(relayersAfter1.result.entries[0].id).to.be.equal(relayersBefore.result.entries[0].id)
        expect(relayersAfter1.result.entries[0].id).to.not.be.equal(relayersAfter1.result.entries[1].id)
        await relayer.stop()
        await CommonUtil.sleep(100)
        const relayersAfter2 = await common.orchestrator.getRelayers()
        expect(relayersAfter2.result.total).to.be.equal(1)
        expect(relayersAfter2.result.entries[0].id).to.be.equal(relayersBefore.result.entries[0].id)
      })
    })
  })
}
