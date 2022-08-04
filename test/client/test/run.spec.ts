import { describe } from "mocha"
import { TestCommon } from "./common"
import { testHeight } from "./height"
import { testHelloWorld } from "./HelloWorld"
import { testMockNodes } from "./mock_node"
import { testNodes } from "./node"
import { testRelay } from "./relay"

describe('run tests', () => {
  let common = new TestCommon()

  before(async function () {
    await common.startAsync()
  })

  beforeEach(function () {
    common.defaultMockNode.clear()
    common.timeoutMockNode.clear()
  })

  after(async function () {
    await common.stopAsync()
  })

  afterEach(function () {
    common.defaultMockNode.clear()
    common.timeoutMockNode.clear()
  })

  testHeight(common)
  testHelloWorld(common)
  testNodes(common)
  testMockNodes(common)
  testRelay(common)
})