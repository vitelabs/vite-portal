import { BaseApp } from "./app";
import { Orchestrator } from "./orchestrator";
import { Relayer } from "./relayer"

export abstract class VitePortal {
  static handleShutdown(app: BaseApp) {
    process.on("SIGINT", async function () {
      await app.stop();
    });
    process.on("SIGTERM", async function () {
      await app.stop();
    });
    process.on("SIGQUIT", async function () {
      await app.stop();
    });
  }

  public static async startOrchestrator(url: string) {
    const app = new Orchestrator(url);
    this.handleShutdown(app)
    await app.start();
    return app;
  }

  public static async startRelayer(url: string) {
    const app = new Relayer(url);
    this.handleShutdown(app)
    await app.start();
    return app;
  }
}