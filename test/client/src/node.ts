import { ChildProcessWithoutNullStreams, spawn } from "child_process"
import { TestContants } from "./constants"

export class VuilderNode {
  process?: ChildProcessWithoutNullStreams
  binPath: string
  stopped: boolean

  constructor() {
    this.binPath = TestContants.DefaultBinPath
    this.stopped = false
  }

  private name(): string {
    return "node"
  }

  private handleProcessOutput = (): void => {
    this.process?.on('error', (error) => {
      console.error(`[${this.name()}] error: ${error}`)
    })

    this.process?.stdout.on('data', (data) => {
      console.log(`[${this.name()}] stdout: ${data}`)
    });

    this.process?.stderr.on('data', (data) => {
      console.error(`[${this.name()}] stderr: ${data}`)
    });
  }

  async start() {
    console.log(`[${this.name()}] Starting...`)

    console.log("Binary:", this.binPath)
    this.process = spawn(
      `npx`,
      [
        "vuilder",
        "--version"
      ],
      {
        cwd: this.binPath,
        detached: true
      },
    )
    this.handleProcessOutput()

    console.log(`[${this.name()}] Started.`)
  }

  async stop() {
    if (this.stopped) return
    console.log(`[${this.name()}] Stopping.`)
    if (this.process?.pid) {
      process.kill(-this.process.pid)
    }
    this.stopped = true
    console.log(`[${this.name()}] Stopped.`)
  }
}