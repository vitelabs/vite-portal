import { it } from "mocha"
import { expect } from "chai"
import { TestCommon } from "./common"
import { TestConstants } from "../src/constants"
import { Relayer } from "../src/relayer"
import { CommonUtil } from "../src/utils"

export function testRelayerJwt(common: TestCommon) {
  let relayer1: Relayer
  let relayer2: Relayer

  before(async function () {
    relayer1 = new Relayer({
      rpcUrl: "",
      rpcAuthUrl: "http://127.0.0.1:56331",
      rpcRelayHttpUrl: "",
      rpcRelayWsUrl: "",
      jwtSecret: TestConstants.DefaultJwtSecret
    }, common.timeout, CommonUtil.uuid())
    relayer2 = new Relayer({
      rpcUrl: "",
      rpcAuthUrl: "http://127.0.0.1:56332",
      rpcRelayHttpUrl: "",
      rpcRelayWsUrl: "",
      jwtSecret: "invalid_secret"
    }, common.timeout, CommonUtil.uuid())
  })

  it('test get nonexistent node', async function () {
    const response = await relayer1.getNode(CommonUtil.uuid())
    expect(response.result).to.be.undefined
    expect(response.error).to.not.be.undefined
    expect(response.error?.code).to.be.equal(-32601)

    const result = await CommonUtil.expectThrowsAsync(() => relayer2.getNode(CommonUtil.uuid()))
    expect(result.response).to.not.be.undefined
    expect(result.response?.status).to.be.equal(403)
  })

  it('test insert node', async function () {
    const chain = CommonUtil.uuid()
    const node = common.createRandomNode(chain)
    const response = await relayer1.putNode(node)
    expect(response.error).to.not.be.undefined
    expect(response.error?.code).to.be.equal(-32601)

    const result = await CommonUtil.expectThrowsAsync(() => relayer2.putNode(node))
    expect(result.response).to.not.be.undefined
    expect(result.response?.status).to.be.equal(403)
  })

  it('test delete nonexistent node', async function () {
    const response = await relayer1.deleteNode(CommonUtil.uuid())
    expect(response.error).to.not.be.undefined
    expect(response.error?.code).to.be.equal(-32601)

    const result = await CommonUtil.expectThrowsAsync(() => relayer2.deleteNode(CommonUtil.uuid()))
    expect(result.response).to.not.be.undefined
    expect(result.response?.status).to.be.equal(403)
  })

  it('test get nodes', async function () {
    const chain = CommonUtil.uuid()
    const response = await relayer1.getNodes(chain)
    expect(response.error).to.not.be.undefined
    expect(response.error?.code).to.be.equal(-32601)

    const result = await CommonUtil.expectThrowsAsync(() => relayer2.getNodes(chain))
    expect(result.response).to.not.be.undefined
    expect(result.response?.status).to.be.equal(403)
  })
}