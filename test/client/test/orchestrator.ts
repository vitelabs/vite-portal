import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { FileUtil, getLocalFileUtil } from "../src/utils"

export function testOrchestrator(common: TestCommon) {
  let fileUtil: FileUtil

  before(async function () {
    fileUtil = getLocalFileUtil()
  })

  describe("testOrchestrator", () => {
    it('test getPaginated relayers', async function () {
      const relayers = await common.orchestrator.getRelayers()
      expect(relayers.total).to.be.equal(1)
      const relayer = relayers.entries[0]
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      console.log(relayer)
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
  })
}
