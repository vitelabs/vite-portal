import { v4 as uuidv4 } from "uuid"
import { expect } from "chai"

export abstract class CommonUtil {
  public static isString(value: any): boolean {
    return typeof value === 'string' || value instanceof String
  }

  public static isNullOrWhitespace(value: any): boolean {
    if (!CommonUtil.isString(value)) {
      return true
    } else {
      return value === null || value === undefined || value.trim() === ''
    }
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(resolve, ms)
    })
  }

  public static retry(conditionFn: () => Promise<boolean>, description: string = '', maxRetries: number = 5) {
    process.stdout.write(description + "\n")
    async function retryWithBackoff(retries: number): Promise<any> {
      try {
        // Make sure we don't wait on the first attempt
        if (retries > 0) {
          const timeToWait = 2 ** retries * 100;
          process.stdout.write(`waiting for ${timeToWait}ms...\n`);
          await CommonUtil.sleep(timeToWait);
        }
        const result = await conditionFn()
        if (result) {
          process.stdout.write("OK\n")
          return
        } else {
          throw new Error("retry failed")
        }
      } catch (e) {
        if (retries < maxRetries) {
          return retryWithBackoff(retries + 1);
        } else {
          process.stdout.write("Max retries reached. Bubbling the error up\n");
          throw e;
        }
      }
    }
    return retryWithBackoff(0);
  }

  public static uuid(): string {
    return uuidv4()
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

  public static expectAsync = async (method: () => Promise<boolean>, timeout: number) => {
    const start = Date.now()
    while (true) {
      if (Date.now() > start + timeout) {
        throw new Error("timed out")
      }
      if (await method()) {
        break
      }
      await CommonUtil.sleep(100)
    }
    expect(await method()).to.be.true
  }
}