import { ChildProcess, ChildProcessWithoutNullStreams, exec, spawn } from "child_process"
import { CommonUtil } from "./utils"
import { TestContants } from "./constants"

export abstract class BaseProcess {
  process?: ChildProcessWithoutNullStreams
  binPath: string
  stopped: boolean

  constructor() {
    this.binPath = TestContants.DefaultBinPath
    this.stopped = false
  }

  abstract name(): string
  abstract command(): string
  abstract killCommand(): string
  abstract args(): string[]
  abstract initAsync(): Promise<void>
  abstract isUp(): Promise<boolean>

  protected extractPort = (url: string): number => {
    const temp = new URL(url)
    return parseInt(temp.port)
  }

  private handleProcessOutput = (process: ChildProcess): void => {
    process?.on('error', (error) => {
      console.error(`[${this.name()}] error: ${error}`)
    })

    process?.stdout?.on('data', (data) => {
      console.log(`[${this.name()}] stdout: ${data}`)
    });

    process?.stderr?.on('data', (data) => {
      console.error(`[${this.name()}] stderr: ${data}`)
    });
  }

  async start() {
    console.log(`[${this.name()}] Starting...`)
    await this.initAsync()

    console.log("Binary:", this.binPath)
    this.process = spawn(
      this.command(),
      this.args(),
      {
        cwd: this.binPath,
        detached: true
      },
    )
    this.handleProcessOutput(this.process)

    await CommonUtil.retry(this.isUp, `Wait for [${this.name()}]`)
    console.log(`[${this.name()}] Started.`)
  }

  async stop() {
    if (this.stopped) return
    console.log(`[${this.name()}] Stopping.`)
    if (this.process?.pid) {
      /* The - in front of the PID instructs process.kill 
         to kill the process group the PID belongs to 
         instead of just the process the PID belongs to. */
      try {
        process.kill(-this.process.pid)
      } catch (error) {
        // If this happens the process most likely did not start properly -> old process still running
        // Possible solution: pgrep <name> | xargs kill -9
        console.log(error)
        await this.killAsync()
      }
    }
    this.stopped = true
    console.log(`[${this.name()}] Stopped.`)
  }

  async killAsync() {
    const command = this.killCommand()
    if (CommonUtil.isNullOrWhitespace(command)) {
      return
    }
    const process = exec(command, {
      cwd: this.binPath,
    })
    this.handleProcessOutput(process)
  }
}