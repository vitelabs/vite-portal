import { NodeEntity, GenericPage, RelayerConfig, AppInfo, JsonRpcResponse } from "./types"
import { BaseApp } from "./app"

export class Relayer extends BaseApp {
  config: RelayerConfig

  constructor(config: RelayerConfig, timeout: number) {
    super(config.rpcRelayHttpUrl, timeout)
    this.config = config
  }

  name(): string {
    return "relayer"
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

  getNodes = async (chain: string, offset?: number, limit?: number): Promise<GenericPage<NodeEntity>> => {
    const params = [
      chain,
      !!offset ? offset : 0,
      !!limit ? limit : 0
    ]
    const response = await this.rpcClient.send(this.config.rpcAuthUrl, "admin_getNodes", params)
    return response.data.result
  }

  getNode = async (id: string): Promise<JsonRpcResponse<NodeEntity>> => {
    const response = await this.rpcClient.send(this.config.rpcAuthUrl, "admin_getNode", [id])
    return response.data
  }

  putNode = async (node: NodeEntity): Promise<JsonRpcResponse<any>> => {
    const response = await this.rpcClient.send(this.config.rpcAuthUrl, "admin_putNode", [node])
    return response.data
  }

  deleteNode = async (id: string): Promise<JsonRpcResponse<any>> => {
    const response = await this.rpcClient.send(this.config.rpcAuthUrl, "admin_deleteNode", [id])
    return response.data
  }
}