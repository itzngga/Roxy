import { MsgHandler } from "../handler";
import { ISimplifiedMessage, ISock, ICmd } from "../types";
import Sticker from "../lib/sticker";
import { MsgType } from "../types/handler";
import fs from "fs";

export default class implements ICmd {
  sock: ISock;
  msgHandler: MsgHandler;
  config: ICmd["config"];

  //config
  constructor() {
    this.config = {
      cmd: "sticker",
    };
  }

  //run
  async run(msg: ISimplifiedMessage, args: any[]): Promise<void> {
    const target = await msg.quoted?.download();
    new Sticker(target)
      .setPackInfo({
        type: msg.type === "videoMessage" ? "video" : "image",
      })
      .build()
      .then((resp) => {
        resp.getBuffer().then((res) => {
          this.sock.sendMessage(msg.from, { sticker: res });
        });
      });
  }
}
