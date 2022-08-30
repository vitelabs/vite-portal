import * as vite from "@vite/vuilder"
import config from "./vite.config.json"
import { TestContants } from "./constants"
import { VitePortal } from "../src/portal"
import { RpcClient } from "../src/client"
import { HttpMockCollector } from "../src/mock_collector"
import { DefaultMockNode, MockNode, TimeoutMockNode } from "../src/mock_node"
import { Relayer } from "../src/relayer"
import { NodeEntity } from "../src/types"
import { CommonUtil } from "../src/utils"
import { Orchestrator } from "../src/orchestrator"

export class TestCommon {
  orchestratorUrl: string
  relayerUrl: string
  providerUrl: string
  nodeHttpUrl: string
  orchestrator!: Orchestrator
  relayer!: Relayer
  provider: any
  deployer: any
  client!: RpcClient
  httpMockCollector: HttpMockCollector
  defaultMockNode: MockNode
  timeoutMockNode: MockNode

  constructor() {
    this.orchestratorUrl = "http://127.0.0.1:57331"
    this.relayerUrl = "http://127.0.0.1:56333"
    this.providerUrl = this.relayerUrl + "/api/v1/client/relay"
    this.nodeHttpUrl = config.networks.local.http
    this.httpMockCollector = new HttpMockCollector(23460)
    this.defaultMockNode = new DefaultMockNode(TestContants.DefaultChain, 23470)
    this.timeoutMockNode = new TimeoutMockNode(TestContants.DefaultChain, 23471)
  }

  startAsync = async () => {
    this.orchestrator = await VitePortal.startOrchestrator(this.orchestratorUrl)
    this.relayer = await VitePortal.startRelayer(this.relayerUrl)
    this.provider = vite.newProvider(this.providerUrl)
    this.deployer = vite.newAccount(config.networks.local.mnemonic, 0, this.provider)
    this.client = new RpcClient()
    this.httpMockCollector.start()
    this.defaultMockNode.start()
    this.timeoutMockNode.start()
  }

  stopAsync = async () => {
    await this.orchestrator.stop()
    await this.relayer.stop()
    this.httpMockCollector.stop()
    this.defaultMockNode.stop()
    this.timeoutMockNode.stop()
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