import { exec } from "child_process";
import path from "path";
import axios, { AxiosInstance } from "axios";
import * as utils from "./utils";

const defaultBinPath = path.join(path.dirname(__dirname), "bin");

export function binPath() {
  return defaultBinPath;
}

export class Relayer {
  binPath: string;
  stopped: boolean;
  provider: AxiosInstance;

  constructor(url: string, binPath: string = defaultBinPath) {
    this.binPath = binPath;
    this.stopped = false;
    this.provider = axios.create({
      baseURL: url,
      timeout: 1000,
    });
  }

  async start() {
    console.log("[VitePortal] Starting...");

    console.log("Node binary:", this.binPath);
    exec(
      `./startup.sh`,
      {
        cwd: this.binPath
      }
    );

    await utils.waitFor(this.isUp, "Wait for VitePortal", 1000);
    console.log("[VitePortal] Started.");
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const response = await this.provider.get("/api")
    return response.data === "vite-portal-relayer"
  }

  async stop() {
    if (this.stopped) return;
    console.log("[VitePortal] Stopping.");
    exec(
      `./shutdown.sh`,
      {
        cwd: this.binPath
      }
    );
    this.stopped = true;
    console.log("[VitePortal] Stopped.");
  }
}