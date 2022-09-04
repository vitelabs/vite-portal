import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { Relayer } from "../src/relayer"
import { CommonUtil, FileUtil, getLocalFileUtil } from "../src/utils"
import { RelayerConfig } from "../src/types"

export function testOrchestrator(common: TestCommon) {
  let fileUtil: FileUtil

  before(async function () {
    fileUtil = getLocalFileUtil()
  })

  describe("testOrchestrator1", () => {
    it('test getAppInfo', async function () {
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      const actual = await common.orchestrator.getAppInfo()
      expect(actual.id).to.not.be.empty
      expect(actual.version).to.be.equal(expectedVersion.trim())
      expect(actual.name).to.be.equal("vite-portal-orchestrator")
    })

    it('test getPaginated relayers', async function () {
      const relayers = await common.orchestrator.getRelayers()
      expect(relayers.total).to.be.equal(1)
      const relayer = relayers.entries[0]
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      expect(relayer.version).to.be.equal(expectedVersion.trim())
      expect(relayer.id).to.not.be.empty
      expect(relayer.transport).to.be.equal("ws")
      expect(relayer.remoteAddress).to.not.be.empty
      expect(relayer.httpInfo.version).to.be.empty
      expect(relayer.httpInfo.userAgent).to.be.equal("Go-http-client/1.1")
      expect(relayer.httpInfo.origin).to.be.empty
      expect(relayer.httpInfo.host).to.be.equal("127.0.0.1:57331")
      expect(relayer.httpInfo.auth).to.be.empty
    })

    describe("testOrchestrator1.1", () => {
      let relayer: Relayer

      after(async function () {
        await relayer.stop()
      })

      it('test spawn relayer', async function () {
        const relayersBefore = await common.orchestrator.getRelayers()
        expect(relayersBefore.total).to.be.equal(1)
        const config: RelayerConfig = {
          rpcUrl: "http://127.0.0.1:56341",
          rpcAuthUrl: "http://127.0.0.1:56342",
          rpcRelayHttpUrl: "http://127.0.0.1:56343",
          rpcRelayWsUrl: "http://127.0.0.1:56344",
        }
        relayer = new Relayer(config, common.timeout)
        await relayer.start()
        const relayersAfter = await common.orchestrator.getRelayers()
        expect(relayersAfter.total).to.be.equal(2)
        expect(relayersAfter.entries[0].id).to.be.equal(relayersBefore.entries[0].id)
        expect(relayersAfter.entries[0].id).to.not.be.equal(relayersAfter.entries[1].id)
      })
    })
  })
}
