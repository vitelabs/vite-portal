import { expect } from "chai"
import { TestCommon } from "./common"

export function testOrchestrator(common: TestCommon) {
  describe("testOrchestrator", () => {
    it('test getPaginated relayers', async function () {
      const relayers = await common.orchestrator.getRelayers()
      expect(relayers.total).to.be.equal(1)
    })
  })
}
