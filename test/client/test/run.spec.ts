import { describe } from "mocha"
import { TestCommon } from "./common"
import { testEmpty } from "./empty"
import { testOrchestratorNode } from "./orchestrator_node"
import { testOrchestratorRelayer } from "./orchestrator_relayer"
import { testOrchestrator } from "./orchestrator"
import { testRelayer } from "./relayer"
import { testRelayerHeight } from "./relayer_height"
import { testRelayerContract } from "./relayer_contract"
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
  testOrchestratorNode(common)
  testOrchestratorRelayer(common)
  testOrchestrator(common)
  testRelayerHeight(common)
  testRelayerContract(common)
  testRelayerNodes(common)
  testRelayerMockNodes(common)
  testRelayerRaw(common)
  testRelayerRelay(common)
})