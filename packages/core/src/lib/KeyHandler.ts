import { EventEmitter } from "events"
import { parseKeypress, type ParsedKey } from "./parse.keypress"
import { singleton } from "../singleton"

type KeyHandlerEventMap = {
  keypress: [ParsedKey]
}

export class KeyHandler extends EventEmitter<KeyHandlerEventMap> {
  constructor() {
    super()

    if (process.stdin.setRawMode) {
      process.stdin.setRawMode(true)
    }
    process.stdin.resume()
    process.stdin.setEncoding("utf8")

    process.stdin.on("data", (key: Buffer) => {
      const str = key.toString()
      
      // Filter out mouse escape sequences
      // SGR mouse format: ESC[<...M or ESC[<...m
      if (/\x1b\[<[^mM]*[mM]/.test(str)) {
        return // Ignore SGR mouse events
      }
      
      // Legacy mouse format: ESC[M followed by 3 bytes
      if (/\x1b\[M/.test(str)) {
        return // Ignore legacy mouse events
      }
      
      // X10 mouse format
      if (/\x1b\[\d+;\d+;\d+M/.test(str)) {
        return // Ignore X10 mouse events
      }
      
      const parsedKey = parseKeypress(key)
      this.emit("keypress", parsedKey)
    })
  }

  public destroy(): void {
    process.stdin.removeAllListeners("data")
  }
}

let keyHandler: KeyHandler | null = null

export function getKeyHandler(): KeyHandler {
  if (!keyHandler) {
    keyHandler = singleton("KeyHandler", () => new KeyHandler())
  }
  return keyHandler
}
