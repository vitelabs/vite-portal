import express from "express"
import * as http from "http"

export class HttpMockCollector {
  app: express.Express
  name: string
  port: number
  server?: http.Server
  results: any[]

  constructor(port: number) {
    this.app = express()
    this.app.use(express.json())
    this.name = "HttpMockCollector"
    this.port = port
    this.results = []
  }

  start = () => {
    if (this.server) {
      return
    }

    this.app.post('/', (req, res) => {
      const result = req.body
      this.results.push(result)
      res.sendStatus(200)
    })

    this.server = this.app.listen(this.port, () => {
      console.log(`[${this.name}] is listening on port ${this.port}.`)
    })
  }

  stop = () => {
    this.server?.close(() => {
      console.log(`[${this.name}] on port ${this.port} has been closed.`)
    })
  }

  clear = () => {
    this.results = []
  }
}