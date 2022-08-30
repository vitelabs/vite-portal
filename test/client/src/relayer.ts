import { NodeEntity, GenericPage, RelayerConfig } from "./types"
import { BaseApp } from "./app"

export class Relayer extends BaseApp {
  config: RelayerConfig

  constructor(config: RelayerConfig) {
    super(config.rpcRelayHttpUrl)
    this.config = config
  }

  name(): string {
    return "relayer"
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const response = await this.axiosClient.get("/api")
    return response.data === "vite-portal-relayer"
  }

  getNodes = async (chain: string, offset?: number, limit?: number): Promise<GenericPage<NodeEntity>> => {
    const params = new URLSearchParams({
      chain
    })
    !!offset && params.append("offset", offset.toString())
    !!limit && params.append("limit", limit.toString())
    const response = await this.axiosClient.get(`/api/v1/db/nodes?${params.toString()}`)
    return response.data
  }

  getName = () => {
    return this.axiosClient.get("/")
  }

  getVersion = () => {
    return this.axiosClient.get("/api/v1")
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