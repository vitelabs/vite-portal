import { exec } from "child_process";
import path from "path";
import axios, { AxiosInstance } from "axios";
import { CommonUtil } from "./utils";

const defaultBinPath = path.join(path.dirname(__dirname), "bin");

export function binPath() {
  return defaultBinPath;
}

export abstract class BaseApp {
  binPath: string;
  stopped: boolean;
  axiosClient: AxiosInstance;

  constructor(url: string, binPath: string = defaultBinPath) {
    this.binPath = binPath;
    this.stopped = false;
    this.axiosClient = axios.create({
      baseURL: url,
      timeout: 1000,
      validateStatus: function () {
        return true;
      }
    });
  }

  abstract name(): string
  abstract isUp(): Promise<boolean>

  private execCallback = (error: any | null, stdout: string, stderr: string): void => {
    if (error) {
      console.error(`[${this.name()}] exec error: ${error}`);
      return;
    }
    console.log(`[${this.name()}] stdout: ${stdout}`);
    console.error(`[${this.name()}] stderr: ${stderr}`);
  }

  async start() {
    console.log(`[${this.name()}] Starting...`);

    console.log("Binary:", this.binPath);
    exec(
      `./startup_${this.name()}.sh`,
      {
        cwd: this.binPath
      },
      //this.execCallback
    );

    await CommonUtil.waitFor(this.isUp, `Wait for [${this.name()}]`, 1000);
    console.log(`[${this.name()}] Started.`);
  }

  async stop() {
    if (this.stopped) return;
    console.log(`[${this.name()}] Stopping.`);
    exec(
      `./shutdown_${this.name()}.sh`,
      {
        cwd: this.binPath
      },
      //this.execCallback
    );
    this.stopped = true;
    console.log(`[${this.name()}] Stopped.`);
  }
}