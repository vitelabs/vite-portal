import * as fs from "fs";
import fetch from "node-fetch";

export abstract class FileUtil {
  abstract readFileAsync(path: string, encoding?: BufferEncoding): Promise<any>;
  abstract writeFileAsync(path: string, data: any, encoding?: BufferEncoding): Promise<any>;
}

export class LocalFileUtil extends FileUtil {
  async readFileAsync(path: string, encoding: BufferEncoding = 'utf8'): Promise<any> {
    return Promise.resolve(fs.readFileSync(path, encoding));
  }
  async writeFileAsync(path: string, data: any, encoding?: BufferEncoding): Promise<any> {
    return Promise.resolve(fs.writeFileSync(path, data, encoding));
  }
}

const localFileUtil = new LocalFileUtil();

export const getLocalFileUtil = () => {
  return localFileUtil;
}

export class NodeFileUtil extends FileUtil {
  async readFileAsync(path: string, encoding: BufferEncoding = 'utf8'): Promise<any> {
    const response = await fetch(path);
    switch (encoding) {
      case 'binary':
        return response.arrayBuffer();
      default:
        return response.text();
    }
  }
  async writeFileAsync(path: string, data: any, encoding?: BufferEncoding): Promise<any> {

  }
}

const nodeFileUtil = new NodeFileUtil();

export const getNodeFileUtil = () => {
  return nodeFileUtil;
}