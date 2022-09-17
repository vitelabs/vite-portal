import { BaseProcess } from "./process"
import { AppInfo, GenericPage, JsonRpcResponse, OrchestratorConfig, RelayerEntity } from "../src/types"
import { RpcHttpClient } from "./client"

export class Orchestrator extends BaseProcess {
  url: string
  authUrl: string
  rpcClient: RpcHttpClient
  rpcAuthClient: RpcHttpClient

  constructor(config: OrchestratorConfig, timeout: number, clientIp: string) {
    super(timeout)
    this.url = config.rpcUrl
    this.authUrl = config.rpcAuthUrl
    this.rpcClient = new RpcHttpClient(timeout, clientIp)
    this.rpcAuthClient = new RpcHttpClient(timeout, clientIp, undefined, config.jwtSecret)
  }

  name(): string {
    return "orchestrator"
  }

  startCommand(): string {
    return "./start_orchestrator.sh"
  }

  killCommand(): string {
    return "./stop_orchestrator.sh"
  }

  startArgs(): string[] {
    return []
  }

  init = async (): Promise<void> => {
    return Promise.resolve()
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const response = await this.rpcClient.send(this.url, "core_getAppInfo")
    return response.data?.result?.name === "vite-portal-orchestrator"
  }

  getAppInfo = async (): Promise<AppInfo> => {
    const response = await this.rpcAuthClient.send(this.authUrl, "core_getAppInfo")
    return response.data.result
  }

  getNodes = async (chain: string, offset?: number, limit?: number): Promise<JsonRpcResponse<GenericPage<RelayerEntity>>> => {
    const params = [
      chain,
      !!offset ? offset : 0,
      !!limit ? limit : 0
    ]
    const response = await this.rpcAuthClient.send(this.authUrl, "admin_getNodes", params)
    return response.data
  }

  getRelayers = async (offset?: number, limit?: number): Promise<JsonRpcResponse<GenericPage<RelayerEntity>>> => {
    const params = [
      !!offset ? offset : 0,
      !!limit ? limit : 0
    ]
    const response = await this.rpcAuthClient.send(this.authUrl, "admin_getRelayers", params)
    return response.data
  }
}