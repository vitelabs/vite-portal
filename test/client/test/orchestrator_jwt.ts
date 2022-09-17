import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestConstants } from "../src/constants"
import { Orchestrator } from "../src/orchestrator"
import { CommonUtil } from "../src/utils"

export function testOrchestratorJwt(common: TestCommon) {
  let orchestrator1: Orchestrator
  let orchestrator2: Orchestrator

  before(async function () {
    orchestrator1 = new Orchestrator({
      rpcUrl: "",
      rpcAuthUrl: "http://127.0.0.1:57331",
      jwtSecret: TestConstants.DefaultJwtSecret
    }, common.timeout, TestConstants.DefaultIpAddress)
    orchestrator2 = new Orchestrator({
      rpcUrl: "",
      rpcAuthUrl: "http://127.0.0.1:57332",
      jwtSecret: "invalid_secret"
    }, common.timeout, TestConstants.DefaultIpAddress)
  })

  describe("testOrchestratorJwt", () => {
    it('test get nodes', async function () {
      const chain = CommonUtil.uuid()
      const response = await orchestrator1.getNodes(chain)
      expect(response.error).to.not.be.undefined
      expect(response.error?.code).to.be.equal(-32601)

      const result = await CommonUtil.expectThrowsAsync(() => orchestrator2.getNodes(chain))
      expect(result.response).to.not.be.undefined
      expect(result.response?.status).to.be.equal(403)
    })

    it('test get relayers', async function () {
      const response = await orchestrator1.getRelayers()
      expect(response.error).to.not.be.undefined
      expect(response.error?.code).to.be.equal(-32601)

      const result = await CommonUtil.expectThrowsAsync(() => orchestrator2.getRelayers())
      expect(result.response).to.not.be.undefined
      expect(result.response?.status).to.be.equal(403)
    })
  })
}