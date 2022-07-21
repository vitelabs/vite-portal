import { describe } from "mocha"
import { expect } from "chai"
import * as vite from "@vite/vuilder"
import config from "./vite.config.json"
import { startRelayer } from "../src/vite"
import { RpcClient } from "../src/client"
import { Relayer } from "../src/relayer"

const relayerUrl = "http://127.0.0.1:56331"
const providerUrl = relayerUrl + "/api/v1/client/relay"
const nodeUrl = "http://127.0.0.1:23456"

let relayer: Relayer
let provider: any
let deployer: any
let client: RpcClient

describe('test HelloWorld', () => {
  before(async function () {
    relayer = await startRelayer(relayerUrl)
    provider = vite.newProvider(providerUrl)
    deployer = vite.newAccount(config.networks.local.mnemonic, 0, provider)
    client = new RpcClient()
    // console.log('deployer', deployer.address)
  })

  after(async function () {
    await relayer.stop()
  })

  it('test height', async function () {
    const method = "ledger_getSnapshotChainHeight"
    const promises: Promise<any>[] = [
      client.send(nodeUrl, method),
      client.send(providerUrl, method)
    ]
    const result = await Promise.all(promises)
    console.log("original:", result[0], "relayed:", result[1])
    expect(result[0].result).to.be.equal(result[1].result)
    const height = await provider.request(method)
    expect(Number(height)).to.be.greaterThan(0)
    expect(Number(height)).to.be.greaterThanOrEqual(Number(result[0].result))
    expect(Number(height)).to.be.greaterThanOrEqual(Number(result[1].result))
  })

  it('test contract', async () => {
    // compile
    const compiledContracts = await vite.compile('HelloWorld.solpp')
    expect(compiledContracts).to.have.property('HelloWorld')

    // deploy
    let helloWorld = compiledContracts.HelloWorld
    helloWorld.setDeployer(deployer).setProvider(provider)
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