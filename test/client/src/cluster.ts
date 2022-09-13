import { RpcHttpClient } from "./client"
import { BaseProcess } from "./process"
import { CommonUtil } from "./utils"

export class NodeCluster extends BaseProcess {
  url?: string
  rpcClient: RpcHttpClient

  constructor(timeout: number) {
    super(timeout)
    this.rpcClient = new RpcHttpClient(timeout)
  }

  name(): string {
    return "cluster"
  }

  startCommand(): string {
    return "./start_cluster.sh"
  }

  killCommand(): string {
    return "./stop_cluster.sh"
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
    if (CommonUtil.isNullOrWhitespace(this.url)) {
      console.log(`[${this.name()}] error: url is not initialized`)
      process.exit(1)
    }
    const response = await this.rpcClient.send(this.url!, "ledger_getSnapshotChainHeight")
    return response.data?.result > 1
  }

  async stop() {
    await super.kill()
    await super.stop()
  }
}