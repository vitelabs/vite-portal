import { NodeEntity, GenericPage, RelayerConfig, AppInfo } from "./types"
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
    const response = await this.rpcClient.send(this.config.rpcAuthUrl, "db_getNodes", params)
    return response.data.result
  }

  getNode = (id: string) => {
    return this.axiosClient.get(`/api/v1/db/nodes/${id}`)
  }

  putNode = (node: NodeEntity) => {
    return this.axiosClient.put(`/api/v1/db/nodes`, node)
  }

  deleteNode = (id: string) => {
    return this.axiosClient.delete(`/api/v1/db/nodes/${id}`)
  }
}