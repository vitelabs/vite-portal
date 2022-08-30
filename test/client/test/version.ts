import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { FileUtil, getLocalFileUtil } from "../src/utils"

export function testVersion(common: TestCommon) {
  let fileUtil: FileUtil

  before(async function () {
    fileUtil = getLocalFileUtil()
  })

  describe("testVersion", () => {
    it('test getVersion', async function () {
      const expected = await fileUtil.readFileAsync("../../shared/pkg/version/buildversion")
      const actual = await common.relayer.getVersion()
      expect("text/plain; charset=UTF-8").to.be.equal(actual.headers["content-type"])
      expect(expected.trim()).to.be.equal(actual.data)
    })

    it('test getName', async function () {
      const actual = await common.relayer.getName()
      expect("text/plain; charset=UTF-8").to.be.equal(actual.headers["content-type"])
      expect("vite-portal-relayer").to.be.equal(actual.data)
    })
  })
}