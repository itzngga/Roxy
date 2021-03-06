import { MsgHandler } from "../handler";
import { ISimplifiedMessage, ISock, ICmd } from "../types";

export default class implements ICmd {
  sock: ISock;
  msgHandler: MsgHandler;
  config: ICmd["config"];
  //config
  constructor() {
    this.config = {
      cmd: "reload",
      aliases: ["r", "rd"],
    };
  }

  //run
  run(
    msg: ISimplifiedMessage,
    args: any[]
  ): Promise<void | never> | void | never {
    this.msgHandler.reloadCmd();
    return void msg.reply("done");
  }
}
