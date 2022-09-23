import { sign, Secret } from "jsonwebtoken"
import { CommonUtil } from "./common.util"
import { Jwt } from "../types";

export abstract class JwtUtil {
  public static CreateDefaultToken(jwt: Jwt): string {
    const now = Math.floor(Date.now() / 1000)
    const payload: {
      [x: string]: any;
    } = {
      "iat": now,
      "exp": now + 10
    }
    if (!CommonUtil.isNullOrWhitespace(jwt.subject)) {
      payload["sub"] = jwt.subject!
    }
    if (!CommonUtil.isNullOrWhitespace(jwt.issuer)) {
      payload["iss"] = jwt.issuer!
    }
    return sign(payload, jwt.secret, { algorithm: "HS256" })
  }
}