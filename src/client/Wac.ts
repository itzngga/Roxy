/* eslint-disable @typescript-eslint/ban-ts-comment */
import pino from "pino";
import SettingService from "../services/settings";
import { ISettings } from "../model/types";
import WAC, {
  DisconnectReason,
  getBinaryNodeChild,
  useSingleFileAuthState,
} from "@adiwajshing/baileys-md";
import { MsgHandler } from "../handler";
import { ISimplifiedMessage, ISock } from "../types/handler";

const { state, saveState } = useSingleFileAuthState(
  "session-md.json",
  pino({ level: "info" })
);
export default class Wac {
  protected msgHandler: MsgHandler;
  protected autoReconnect = true;
  public sock: ISock;

  constructor(public id: string) {
    // this.connectOptions.connectCooldownMs = 15 * 1000;
    // this.connectOptions.alwaysUseTakeover = true;
    // this.connectOptions.queryChatsTillReceived = true;
    // this.browserDescription = ["Roxy", "Safari", "10.0"];
    // this.__generateSettings(id);

    this.start();
  }

  protected prepareCustomFunction = (): void => {
    this.sock.groupQueryInvite = async (code) => {
      const results = await this.sock.query({
        tag: "iq",
        attrs: {
          type: "get",
          xmlns: "w:g2",
          to: "@g.us",
        },
        content: [{ tag: "invite", attrs: { code } }],
      });
      const group = getBinaryNodeChild(results, "group");
      return group.attrs;
    };
  };
  protected start = () => {
    this.sock = WAC({
      printQRInTerminal: true,
      auth: state,
      logger: pino({ level: "silent" }),
    }) as ISock;

    this.prepareCustomFunction();
    this.prepareMsgHandler();

    // creds.update
    this.sock.ev.on("creds.update", saveState);
    // connection.update
    this.sock.ev.on("connection.update", (up) => {
      const { lastDisconnect, connection } = up;
      if (connection === "close") {
        if (
          //@ts-expect-error
          lastDisconnect?.error?.output?.statusCode !==
          DisconnectReason?.loggedOut
        ) {
          this.start();
        } else {
          console.log("Closed");
        }
      }
    });
    // group-participants.update
    // this.sock.ev.on("group-participants.update", (json) => {
    //   joinhandler(json, this.sock);
    // });
  };

  // private __generateSettings = async (id: string) => {
  //   const res = await SettingService.getOneSetting(id);
  //   if (!res) return console.error("Could not find setting with id: " + id);
  //   this.settings = res as ISettings;
  // };

  public prepareMsgHandler = () => {
    // this.sock.ev.removeAllListeners("messages.upsert");
    this.msgHandler = new MsgHandler(this.sock);
    return this.msgHandler.__initHandler();
  };
}
