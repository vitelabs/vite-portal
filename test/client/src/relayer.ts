import { exec } from "child_process";
import path from "path";
import axios, { AxiosInstance } from "axios";
import { CommonUtil } from "./utils";

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
      validateStatus: function () {
        return true;
      }
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

    await CommonUtil.waitFor(this.isUp, "Wait for VitePortal", 1000);
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

  getNodes = async (chain: string, offset?: number, limit?: number): Promise<GenericPage<NodeEntity>> => {
    const params = new URLSearchParams({
      chain
    })
    !!offset && params.append("offset", offset.toString())
    !!limit && params.append("limit", limit.toString())
    const response = await this.provider.get(`/api/v1/db/nodes?${params.toString()}`)
    return response.data
  }

  getNode = (id: string) => {
    return this.provider.get(`/api/v1/db/nodes/${id}`)
  }

  putNode = (node: NodeEntity) => {
    return this.provider.put(`/api/v1/db/nodes/${node.id}`, node)
  }

  deleteNode = (id: string) => {
    return this.provider.delete(`/api/v1/db/nodes/${id}`)
  }
}

export type NodeEntity = {
  id: string
  chain: string
  ipAddress: string
  rewardAddress?: string
}

export type GenericPage<T> = {
  entries: T[]
  limit: number
  offset: number
  total: number
}