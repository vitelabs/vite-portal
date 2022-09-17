import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { RpcWsClient, RpcWsClientWrapper } from "../src/client"
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

    var runs1 = [
      { id: '1', options: { jwtSubject: undefined, jwtSecret: undefined } },
      { id: '2', options: { jwtSubject: undefined, jwtSecret: "invalid_secret" } },
    ]

    runs1.forEach(function (run) {
      it('test connect/disconnect error ' + run.id, async function () {
        const url = "ws://127.0.0.1:57332"
        const client = new RpcWsClient(common.timeout, url, CommonUtil.uuid(), run.options.jwtSubject, run.options.jwtSecret)
        expect(client.error).to.be.undefined
        await CommonUtil.sleep(10)
        expect(client.error).to.not.be.undefined
        expect(client.error?.message).to.be.equal("Unexpected server response: 403")
      })
    })

    var runs2 = [
      {
        id: '1', options: {
          expected1: "net_nodeInfo",
          expected2: "failed to call \'net_nodeInfo\'",
          jwtSubject: undefined,
          jwtSecret: TestConstants.DefaultJwtSecret
        }
      },
      {
        id: '2', options: {
          expected1: "core_getAppInfo",
          expected2: "context deadline exceeded",
          jwtSubject: TestConstants.DefaultJwtRelayerSubject,
          jwtSecret: TestConstants.DefaultJwtSecret
        }
      },
    ]

    runs2.forEach(function (run) {
      it('test connect/disconnect success ' + run.id, async function () {
        const url = "ws://127.0.0.1:57332"
        const client = new RpcWsClient(common.timeout, url, CommonUtil.uuid(), run.options.jwtSubject, run.options.jwtSecret)
        const wrapper = new RpcWsClientWrapper(client)

        await CommonUtil.expectAsync(async () => wrapper.connected === true, common.timeout)
        await CommonUtil.expectAsync(async () => wrapper.requests.length === 1, common.timeout)
        expect(wrapper.requests[0].method).to.be.equal(run.options.expected1)
        await CommonUtil.expectAsync(async () => wrapper.connected === false, 6000)
        expect(wrapper.errors.length).to.be.equal(1)
        expect(wrapper.errors[0].error?.code).to.be.equal(-32000)
        expect(wrapper.errors[0].error?.message).to.be.equal(run.options.expected2)
      })
    })
  })
}