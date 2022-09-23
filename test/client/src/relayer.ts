import axios, { AxiosInstance } from "axios"
import { NodeEntity, GenericPage, RelayerConfig, AppInfo, JsonRpcResponse, Jwt } from "./types"
import { RpcHttpClient } from "./client"
import { TestConstants } from "./constants"
import { BaseProcess } from "./process"

export class Relayer extends BaseProcess {
  config: RelayerConfig
  rpcClient: RpcHttpClient
  rpcAuthClient: RpcHttpClient
  axiosClient: AxiosInstance

  constructor(config: RelayerConfig, timeout: number, clientIp: string) {
    super(timeout)
    this.config = config
    this.rpcClient = new RpcHttpClient(timeout, clientIp)
    const jwt: Jwt = {
      secret: config.jwtSecret
    }
    this.rpcAuthClient = new RpcHttpClient(timeout, clientIp, jwt)
    this.axiosClient = axios.create({
      baseURL: config.rpcRelayHttpUrl,
      timeout: timeout,
      headers: {
        [TestConstants.DefaultHeaderTrueClientIp]: clientIp
      },
      validateStatus: function () {
        return true
      }
    })
  }

  name(): string {
    return "relayer"
  }

  startCommand(): string {
    return "./start_relayer.sh"
  }

  killCommand(): string {
    return "./stop_relayer.sh"
  }

  startArgs(): string[] {
    const overrides = {
      rpcPort: this.extractPort(this.config.rpcUrl),
      rpcAuthPort: this.extractPort(this.config.rpcAuthUrl),
      rpcRelayHttpPort: this.extractPort(this.config.rpcRelayHttpUrl),
      rpcRelayWsPort: this.extractPort(this.config.rpcRelayWsUrl)
    }
    const args = [
      JSON.stringify(overrides)
    ]
    return args
  }

  init = async (): Promise<void> => {
    return Promise.resolve()
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const response = await this.rpcClient.send(this.config.rpcUrl, "core_getAppInfo")
    return response.data?.result?.name === "vite-portal-relayer"
  }

  getAppInfo1 = async (): Promise<AppInfo> => {
    const response = await this.axiosClient.get(`/`)
    return response.data
  }

  getAppInfo2 = async (): Promise<AppInfo> => {
    const response = await this.rpcClient.send(this.config.rpcUrl, "core_getAppInfo")
    return response.data.result
  }

  getNodes = async (chain: string, offset?: number, limit?: number): Promise<JsonRpcResponse<GenericPage<NodeEntity>>> => {
    const params = [
      chain,
      !!offset ? offset : 0,
      !!limit ? limit : 0
    ]
    const response = await this.rpcAuthClient.send(this.config.rpcAuthUrl, "admin_getNodes", params)
    return response.data
  }

  getNode = async (id: string): Promise<JsonRpcResponse<NodeEntity>> => {
    const response = await this.rpcAuthClient.send(this.config.rpcAuthUrl, "admin_getNode", [id])
    return response.data
  }

  putNode = async (node: NodeEntity): Promise<JsonRpcResponse<any>> => {
    const response = await this.rpcAuthClient.send(this.config.rpcAuthUrl, "admin_putNode", [node])
    return response.data
  }

  deleteNode = async (id: string): Promise<JsonRpcResponse<any>> => {
    const response = await this.rpcAuthClient.send(this.config.rpcAuthUrl, "admin_deleteNode", [id])
    return response.data
  }
}