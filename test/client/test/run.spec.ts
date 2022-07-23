import { describe } from "mocha"
import { TestCommon } from "./common"
import { testHeight } from "./height"
import { testHelloWorld } from "./HelloWorld"
import { testNodes } from "./node"

describe('run tests', () => {
  let common = new TestCommon()

  before(async function () {
    await common.startAsync()
  })

  after(async function () {
    await common.stopAsync()
  })

  testHeight(common)
  testHelloWorld(common)
  testNodes(common)
})