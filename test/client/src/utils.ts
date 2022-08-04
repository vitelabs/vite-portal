import { v4 as uuidv4 } from "uuid";
import { expect } from "chai"

export abstract class CommonUtil {
  public static isString(value: any): boolean {
    return typeof value === 'string' || value instanceof String;
  }

  public static isNullOrWhitespace(value: any): boolean {
    if (!CommonUtil.isString(value)) {
      return true;
    } else {
      return value === null || value === undefined || value.trim() === '';
    }
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(resolve, ms);
    });
  }

  public static waitFor(conditionFn: () => Promise<boolean>, description: string = '', pollInterval: number = 1000) {
    process.stdout.write(description);
    const poll = (resolve: any) => {
      conditionFn().then((result) => {
        if (result) {
          console.log(" OK");
          resolve();
        } else {
          process.stdout.write('.');
          setTimeout(() => poll(resolve), pollInterval);
        }
      }).catch(() => {
        process.stdout.write('.');
        setTimeout(() => poll(resolve), pollInterval);
      });
    }
    return new Promise(poll);
  }

  public static uuid(): string {
    return uuidv4();
  }

  public static expectThrowsAsync = async (method: () => Promise<any>, errorMessage?: string) => {
    let error: any
    try {
      await method()
    }
    catch (err) {
      error = err
    }
    expect(error).to.not.be.undefined
    if (errorMessage) {
      expect(error.message).to.equal(errorMessage)
    }
    return error
  }
}