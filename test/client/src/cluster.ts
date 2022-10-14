import { RpcHttpClient } from "./client"
import { TestConstants } from "./constants"
import { Orchestrator } from "./orchestrator"
import { BaseProcess } from "./process"
import { CommonUtil } from "./utils"

export class NodeCluster extends BaseProcess {
  orchestrator?: Orchestrator

  constructor(timeout: number) {
    super(timeout)
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
    if (!this.orchestrator) {
      console.log(`[${this.name()}] error: orchestrator is not initialized`)
      process.exit(1)
    }
    const response = await this.orchestrator.getNodes(TestConstants.SupportedChains.ViteBuidl)
    return response.result.total >= 3
  }

  async stop() {
    await super.kill()
    await super.stop()
  }
}