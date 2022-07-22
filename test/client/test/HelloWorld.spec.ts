import { describe } from "mocha"
import { expect } from "chai"
import * as vite from "@vite/vuilder"
import { TestCommon } from "./common"

describe('test HelloWorld', () => {
  let common: TestCommon

  before(async function () {
    common = new TestCommon()
    await common.startAsync()
  })

  after(async function () {
    await common.stopAsync()
  })

  it('test contract', async () => {
    // compile
    const compiledContracts = await vite.compile('HelloWorld.solpp')
    expect(compiledContracts).to.have.property('HelloWorld')

    // deploy
    let helloWorld = compiledContracts.HelloWorld
    helloWorld.setDeployer(common.deployer).setProvider(common.provider)
    await helloWorld.deploy({})
    expect(helloWorld.address).to.be.a('string')
    console.log(helloWorld.address)

    // check default value of data
    let result = await helloWorld.query('data', [])
    console.log('return', result)
    expect(result).to.be.an('array').with.lengthOf(1)
    expect(result![0]).to.be.equal('123')

    // call HelloWorld.set(456)
    await helloWorld.call('set', ['456'], {})

    // check value of data
    result = await helloWorld.query('data', [])
    console.log('return', result)
    expect(result).to.be.an('array').with.lengthOf(1)
    expect(result![0]).to.be.equal('456')
  })
})