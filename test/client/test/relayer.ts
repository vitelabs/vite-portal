import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { FileUtil, getLocalFileUtil } from "../src/utils"

export function testRelayer(common: TestCommon) {
  let fileUtil: FileUtil

  before(async function () {
    fileUtil = getLocalFileUtil()
  })

  describe("testRelayer", () => {
    it('test getAppInfo', async function () {
      const expectedVersion = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      const actual1 = await common.relayer.getAppInfo1()
      expect(actual1.id).to.not.be.empty
      expect(actual1.version).to.be.equal(expectedVersion.trim())
      expect(actual1.name).to.be.equal("vite-portal-relayer")
      const actual2 = await common.relayer.getAppInfo2()
      expect(actual2.id).to.be.equal(actual1.id)
      expect(actual2.version).to.be.equal(actual1.version)
      expect(actual2.name).to.be.equal(actual1.name)
    })
  })
}