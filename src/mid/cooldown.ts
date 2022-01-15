import { MsgHandler } from "../handler";
import { ISimplifiedMessage, ISock, IMid } from "../types";

export default class implements IMid {
  public cooldown: Map<string, number> = new Map();

  sock: ISock;
  msgHandler: MsgHandler;
  config: IMid["config"];
  //config
  constructor() {
    this.config = {
      name: "cooldown",
      mode: "before",
    };
  }

  //run
  run(msg: ISimplifiedMessage, args: any[]): any {
    const now = new Date().getTime();
    if (!this.cooldown.has(msg.from)) {
      this.cooldown.set(msg.from, now + 5 * 1000);
      return true;
    } else {
      const expiration = this.cooldown.get(msg.from) || new Date().getTime();
      if (now <= expiration) {
        if (msg.isGroup) {
          const timeLeft = expiration - now;
          // printSpam(isGroup, sender, gcName);
          return msg.reply(
            `This group is on cooldown, please wait another _${(
              timeLeft / 1000
            ).toFixed(1)} second(s)_`
          );
        } else if (!msg.isGroup) {
          const timeLeft = expiration - now;
          // printSpam(isGroup, sender);
          return msg.reply(
            `You are on cooldown, please wait another _${(
              timeLeft / 1000
            ).toFixed(1)} second(s)_`
          );
        }
      }
      setTimeout(() => this.cooldown.delete(msg.from), 5 * 1000);
      return true;
    }
  }
}
