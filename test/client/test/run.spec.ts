import { describe } from "mocha"
import { TestCommon } from "./common"
import { testEmpty } from "./empty"
import { testOrchestratorCluster } from "./orchestrator_cluster"
import { testOrchestratorNode } from "./orchestrator_node"
import { testOrchestratorRelayer } from "./orchestrator_relayer"
import { testOrchestrator } from "./orchestrator"
import { testRelayer } from "./relayer"
import { testRelayerContract } from "./relayer_contract"
import { testRelayerHeight } from "./relayer_height"
import { testRelayerJwt } from "./relayer_jwt"
import { testRelayerNodes } from "./relayer_node"
import { testRelayerMockNodes } from "./relayer_mock_node"
import { testRelayerRaw } from "./relayer_raw"
import { testRelayerRelay } from "./relayer_relay"

describe('run tests', () => {
  let common = new TestCommon()

  before(async function () {
    await common.startAsync()
  })

  beforeEach(function () {
    common.httpMockCollector.clear()
    common.defaultMockNode.clear()
    common.timeoutMockNode.clear()
  })

  after(async function () {
    await common.stopAsync()
  })

  afterEach(function () {
    common.httpMockCollector.clear()
    common.defaultMockNode.clear()
    common.timeoutMockNode.clear()
  })

  testEmpty(common)
  testRelayer(common)
  testOrchestratorCluster(common)
  testOrchestratorNode(common)
  testOrchestratorRelayer(common)
  testOrchestrator(common)
  testRelayerHeight(common)
  testRelayerJwt(common)
  testRelayerContract(common)
  testRelayerNodes(common)
  testRelayerMockNodes(common)
  testRelayerRaw(common)
  testRelayerRelay(common)
})