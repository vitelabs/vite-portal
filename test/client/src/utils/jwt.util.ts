import { sign, Secret } from "jsonwebtoken"
import { CommonUtil } from "./common.util"

export abstract class JwtUtil {
  public static CreateDefaultToken(jwtSecret: string, subject?: string): string {
    const now = Math.floor(Date.now() / 1000)
    const payload: {
      [x: string]: any;
    } = {
      "iat": now,
      "exp": now + 10
    }
    if (!CommonUtil.isNullOrWhitespace(subject)) {
      payload["sub"] = subject!
    }
    return sign(payload, jwtSecret, { algorithm: "HS256" })
  }
}