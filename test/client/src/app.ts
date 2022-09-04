import { ChildProcessWithoutNullStreams, spawn } from "child_process"
import axios, { AxiosInstance } from "axios"
import { CommonUtil } from "./utils"
import { RpcClient } from "./client"
import { TestContants } from "./constants"

export abstract class BaseApp {
  process?: ChildProcessWithoutNullStreams
  url: string
  binPath: string
  stopped: boolean
  rpcClient: RpcClient
  axiosClient: AxiosInstance

  constructor(url: string, timeout: number) {
    this.url = url
    this.binPath = TestContants.DefaultBinPath
    this.stopped = false
    this.rpcClient = new RpcClient(timeout)
    this.axiosClient = axios.create({
      baseURL: url,
      timeout: timeout,
      validateStatus: function () {
        return true
      }
    })
  }

  abstract name(): string
  abstract isUp(): Promise<boolean>
  abstract getConfigOverrides(): string

  protected extractPort = (url: string): number => {
    const temp = new URL(url)
    return parseInt(temp.port)
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
      `./start_${this.name()}.sh`,
      [
        this.getConfigOverrides()
      ],
      {
        cwd: this.binPath,
        detached: true
      },
    )
    this.handleProcessOutput()

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
      process.kill(-this.process.pid)
    }
    this.stopped = true
    console.log(`[${this.name()}] Stopped.`)
  }
}