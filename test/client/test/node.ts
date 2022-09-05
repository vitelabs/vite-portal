import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"

export function testNode(common: TestCommon) {
  describe("testNode", () => {
    it('test connect', async function () {
      expect(true).to.be.true
    })
  })
}