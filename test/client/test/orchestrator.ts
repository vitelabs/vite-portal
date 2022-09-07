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
    it('test getAppInfo', async function () {
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      const actual = await common.orchestrator.getAppInfo()
      expect(actual.id).to.not.be.empty
      expect(actual.version).to.be.equal(expectedVersion.trim())
      expect(actual.name).to.be.equal("vite-portal-orchestrator")
    })
  })
}
