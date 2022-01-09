import { ICmd } from "../types/cmd";
import { MsgHandler } from "../handler";
import { ISimplifiedMessage, ISock } from "../types/handler";

export default class implements ICmd {
  sock: ISock;
  msgHandler: MsgHandler;
  config: ICmd["config"];
  //config
  constructor() {
    this.config = {
      cmd: "tes",
      aliases: ["anj"],
    };
  }

  //run
  run(
    msg: ISimplifiedMessage,
    args: any[]
  ): Promise<void | never> | void | never {
    return void msg.reply("yes");
  }
}
