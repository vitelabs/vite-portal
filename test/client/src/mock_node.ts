import express from "express"
import * as http from "http"
import { JsonRpcResponse, NodeEntity } from "./types";
import { CommonUtil } from "./utils";

export abstract class MockNode {
  app: express.Express
  name?: string
  chain: string
  port: number
  server?: http.Server
  requests: any[]
  entity: NodeEntity

  constructor(chain: string, port: number) {
    this.app = express()
    this.chain = chain
    this.port = port
    this.requests = []
    this.entity = {
      id: CommonUtil.uuid(),
      chain: chain,
      rpcHttpUrl: `http://127.0.0.1:${this.port}`,
      rpcWsUrl: `ws://127.0.0.1:${this.port}`
    }
  }

  abstract start(): void

  stop = () => {
    this.server?.close(() => {
      console.log(`[${this.name}] on port ${this.port} has been closed.`)
    })
  }

  clear = () => {
    this.requests = []
  }
}

export class DefaultMockNode extends MockNode {
  static DEFAULT_RESPONSE: JsonRpcResponse = {
    jsonrpc: "2.0",
    id: 1,
    result: "This is a test response!"
  }

  constructor(chain: string, port: number) {
    super(chain, port)
    this.name = "DefaultMockNode"
  }

  start = () => {
    if (this.server) {
      return
    }

    this.app.get('/', (req, res) => {
      this.requests.push(req)
      res.json(DefaultMockNode.DEFAULT_RESPONSE)
    })

    this.app.post('/', (req, res) => {
      this.requests.push(req)
      res.json(DefaultMockNode.DEFAULT_RESPONSE)
    })

    this.server = this.app.listen(this.port, () => {
      console.log(`[${this.name}] is listening on port ${this.port}.`)
    })
  }
}

export class TimeoutMockNode extends MockNode {
  constructor(chain: string, port: number) {
    super(chain, port)
    this.name = "TimeoutMockNode"
  }

  start = () => {
    if (this.server) {
      return
    }

    this.app.get('/', (req, res) => {
      this.requests.push(req)
    })

    this.app.post('/', (req, res) => {
      this.requests.push(req)
    })

    this.server = this.app.listen(this.port, () => {
      console.log(`[${this.name}] is listening on port ${this.port}.`)
    })
  }
}