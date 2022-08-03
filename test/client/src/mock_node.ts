import express from "express"
import * as http from "http"
import { JsonRpcResponse } from "./types";

export abstract class MockNode {
  app: express.Express
  name?: string
  port?: number
  server?: http.Server

  constructor() {
    this.app = express()
  }

  url = () => {
    return `http://127.0.0.1:${this.port}`
  }

  abstract start(port: number): void

  stop = () => {
    this.server?.close(() => {
      console.log(`[${this.name}] on port ${this.port} has been closed.`)
    })
  }
}

export class DefaultMockNode extends MockNode {
  static DEFAULT_RESPONSE: JsonRpcResponse = {
    jsonrpc: "2.0",
    id: 1,
    result: "This is a test response!"
  }

  constructor() {
    super()
    this.name = "DefaultMockNode"
  }

  start = (port: number) => {
    if (this.server) {
      return
    }

    this.port = port

    this.app.get('/', (req, res) => {
      res.json(DefaultMockNode.DEFAULT_RESPONSE)
    })

    this.app.post('/', (req, res) => {
      res.json(DefaultMockNode.DEFAULT_RESPONSE)
    })

    this.server = this.app.listen(port, () => {
      console.log(`[${this.name}] is listening on port ${port}.`)
    })
  }
}
