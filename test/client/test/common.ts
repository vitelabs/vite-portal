import * as vite from "@vite/vuilder"
import config from "./vite.config.json"
import { startRelayer } from "../src/vite"
import { RpcClient } from "../src/client"
import { Relayer } from "../src/relayer"
import { NodeEntity } from "../src/relayer"
import { CommonUtil } from "../src/utils"

export class TestCommon {
  relayerUrl: string
  providerUrl: string
  nodeUrl: string
  relayer!: Relayer
  provider: any
  deployer: any
  client!: RpcClient

  constructor() {
    this.relayerUrl = "http://127.0.0.1:56331"
    this.providerUrl = this.relayerUrl + "/api/v1/client/relay"
    this.nodeUrl = "http://127.0.0.1:23456"
  }

  startAsync = async () => {
    this.relayer = await startRelayer(this.relayerUrl)
    this.provider = vite.newProvider(this.providerUrl)
    this.deployer = vite.newAccount(config.networks.local.mnemonic, 0, this.provider)
    this.client = new RpcClient()
  }

  stopAsync = async () => {
    await this.relayer.stop()
  }

  createRandomNode = (chain: string): NodeEntity => {
    return {
      id: CommonUtil.uuid(),
      chain,
      rpcHttpUrl: "http://127.0.0.1:23456",
      rpcWsUrl: "ws://127.0.0.1:23457"
    }
  }

  insertNodeAsync = async (chain: string): Promise<NodeEntity> => {
    const node = this.createRandomNode(chain)
    const result = await this.relayer.putNode(node)
    if (result.status != 200) {
      throw new Error("node could not be inserted")
    }
    return node
  }
}