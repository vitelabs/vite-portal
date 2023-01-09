import * as vite from "@vite/vuilder"
import config from "./vite.config.json"
import { RpcHttpClient } from "../src/client"
import { TestConstants } from "../src/constants"
import { Kafka } from "../src/kafka"
import { HttpMockCollector } from "../src/mock_collector"
import { DefaultMockNode, MockNode, TimeoutMockNode } from "../src/mock_node"
import { Orchestrator } from "../src/orchestrator"
import { VitePortal } from "../src/portal"
import { Relayer } from "../src/relayer"
import { NodeEntity, OrchestratorConfig, RelayerConfig } from "../src/types"
import { CommonUtil } from "../src/utils"

export class TestCommon {
  timeout: number
  kafka: Kafka
  orchestratorConfig: OrchestratorConfig
  relayerConfig: RelayerConfig
  providerUrl: string
  nodeHttpUrl: string
  orchestrator!: Orchestrator
  relayer!: Relayer
  provider: any
  deployer: any
  client!: RpcHttpClient
  httpMockCollector: HttpMockCollector
  defaultMockNode: MockNode
  timeoutMockNode: MockNode

  constructor() {
    this.timeout = 2100
    this.kafka = new Kafka(this.timeout)
    this.orchestratorConfig = {
      rpcUrl: "http://127.0.0.1:57331",
      rpcAuthUrl: "http://127.0.0.1:57332",
      jwtSecret: TestConstants.DefaultJwtSecret
    }
    this.relayerConfig = {
      rpcUrl: "http://127.0.0.1:56331",
      rpcAuthUrl: "http://127.0.0.1:56332",
      rpcRelayHttpUrl: "http://127.0.0.1:56333",
      rpcRelayWsUrl: "http://127.0.0.1:56334",
      jwtSecret: TestConstants.DefaultJwtSecret
    }
    this.providerUrl = this.relayerConfig.rpcRelayHttpUrl + "/relay"
    this.nodeHttpUrl = config.networks.local.http
    this.httpMockCollector = new HttpMockCollector(23460)
    this.defaultMockNode = new DefaultMockNode(TestConstants.DefaultChain, 23470)
    this.timeoutMockNode = new TimeoutMockNode(TestConstants.DefaultChain, 23471)
  }

  startAsync = async () => {
    VitePortal.startCleanup()
    await this.kafka.start()
    this.orchestrator = await VitePortal.startOrchestrator(this.orchestratorConfig, this.timeout)
    this.relayer = await VitePortal.startRelayer(this.relayerConfig, this.timeout)
    this.provider = vite.newProvider(this.providerUrl)
    this.deployer = vite.newAccount(config.networks.local.mnemonic, 0, this.provider)
    this.client = new RpcHttpClient(this.timeout, "1.2.3.4")
    this.httpMockCollector.start()
    this.defaultMockNode.start()
    this.timeoutMockNode.start()
  }

  stopAsync = async () => {
    await this.kafka.stop()
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
    const response = await this.relayer.putNode(node)
    if (response.error) {
      throw new Error("node could not be inserted")
    }
    return node
  }
}