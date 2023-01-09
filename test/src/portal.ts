import { exec } from "child_process"
import { BaseProcess } from "./process"
import { TestConstants } from "./constants"
import { Orchestrator } from "./orchestrator"
import { Relayer } from "./relayer"
import { OrchestratorConfig, RelayerConfig } from "./types"

export abstract class VitePortal {
  static handleShutdown(p: BaseProcess) {
    process.on("SIGINT", async function () {
      await p.stop()
    })
    process.on("SIGTERM", async function () {
      await p.stop()
    })
    process.on("SIGQUIT", async function () {
      await p.stop()
    })
  }

  public static startCleanup() {
    exec(
      `./start_cleanup.sh`,
      {
        cwd: TestConstants.DefaultBinPath
      },
    )
  }

  public static async startOrchestrator(config: OrchestratorConfig, timeout: number) {
    const app = new Orchestrator(config, timeout, "1.1.1.1")
    this.handleShutdown(app)
    await app.start()
    return app
  }

  public static async startRelayer(config: RelayerConfig, timeout: number) {
    const app = new Relayer(config, timeout, "1.1.1.2")
    this.handleShutdown(app)
    await app.start()
    return app
  }
}