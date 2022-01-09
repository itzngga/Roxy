/* eslint-disable @typescript-eslint/ban-ts-comment */
import util from "../util";
import { join } from "path";
import { ICmd } from "../types/cmd";
import { EventEmitter } from "events";
import { proto } from "@adiwajshing/baileys-md";
import { MsgType, ISimplifiedMessage, ISock } from "../types/handler";
const multi_pref = new RegExp(
  "^[" + "!#$%&?/;:,.<>~-+=".replace(/[|\\{}()[\]^$+*?.\-^]/g, "\\$&") + "]"
);
const cmdEvent = new EventEmitter();
export class Cmd {
  static cmdMap: Map<string, ICmd> = new Map();

  static init = async () => {
    this.cmdMap.clear();
    const files = util.readDirRecursive(join(__dirname, "\\..\\cmd"));
    for await (const path of files) {
      try {
        import(path)
          .then((imp) => {
            if (imp.default.prototype) {
              const cmd = new imp.default();
              this.cmdMap.set(cmd.config.cmd, cmd);
            } else {
              this.cmdMap.set(imp.default.config.cmd, imp.default);
            }
          })
          .catch((error) => util.roxyLog("fatal", "some command error", error));
      } catch (error) {
        return util.roxyLog("fatal", "some command error", error);
      }
    }
    return util.roxyLog("info", `loaded cmd: ${this.cmdMap.size}`);
  };
  static loadCmd = async () => {
    await this.init();
    cmdEvent.emit("cmdLoaded", null);
  };
  static reloadCmd = async () => {
    await this.init();
    cmdEvent.emit("cmdLoaded", null);
    return this.cmdMap.size;
  };
}
export class MsgHandler {
  public cmdMap: Map<string, ICmd> = new Map();
  public aliasMap: Map<string, ICmd> = new Map();
  public cooldown: Map<string, any> = new Map();

  constructor(protected sock: ISock) {}

  public __initHandler(): void {
    this.sock.ev.on("messages.upsert", async (m) => this.handle(m));
    this.loadCmd();
  }

  public reloadCmd = () => Cmd.reloadCmd();

  protected loadCmd(): void {
    cmdEvent.on("cmdLoaded", () => {
      this.cmdMap.clear();
      this.aliasMap.clear();
      Cmd.cmdMap.forEach((cmd, key) => {
        try {
          cmd.sock = this.sock;
          cmd.msgHandler = this;
          this.cmdMap.set(key, cmd);
          if (cmd.config.aliases)
            cmd.config.aliases.forEach((alias) =>
              this.aliasMap.set(alias, cmd)
            );
        } catch (error) {
          return util.roxyLog("fatal", "some command error", error);
        }
      });
    });
  }

  protected async handle(m: any): Promise<void> {
    if (m.type !== "notify") return;
    const msg = this.simplified(JSON.parse(JSON.stringify(m.messages[0])));
    if (!msg.message || msg.fromMe) return;
    if (msg.key && msg.key.remoteJid === "status@broadcast") return;
    if (
      (msg.type as string) === "protocolMessage" ||
      (msg.type as string) === "senderKeyDistributionMessage" ||
      !msg.type
    )
      return;

    let { body } = msg;
    const temp_pref = multi_pref.test(body)
      ? body?.split("").shift() || ""
      : "!";
    const { type, isGroup, sender, from } = msg;
    body =
      type === "conversation" && body?.startsWith(temp_pref)
        ? body
        : (type === "imageMessage" || type === "videoMessage") &&
          body &&
          body?.startsWith(temp_pref)
        ? body
        : type === "extendedTextMessage" && body?.startsWith(temp_pref)
        ? body
        : type === "buttonsResponseMessage" && body?.startsWith(temp_pref)
        ? body
        : type === "listResponseMessage" && body?.startsWith(temp_pref)
        ? body
        : type === "templateButtonReplyMessage" && body?.startsWith(temp_pref)
        ? body
        : "";
    const arg = body.substring(body.indexOf(" ") + 1);
    const args = body.trim().split(/ +/).slice(1);
    const isCmd = body.startsWith(temp_pref);
    const gcMeta = isGroup ? await this.sock.groupMetadata(from) : null;
    const gcName = isGroup ? gcMeta?.subject : null;

    // Log
    // printLog(isCmd, sender, gcName, isGroup);

    // @ts-expect-error
    const cmd = body
      .slice(temp_pref?.length)
      .trim()
      .split(/ +/)
      .shift()
      .toLowerCase();

    // if (!this.cooldown.has(from)) {
    //   this.cooldown.set(from, Date.now() + 5 * 1000);
    // }

    const command = this.cmdMap.get(cmd) || this.aliasMap.get(cmd);

    if (!command)
      return void this.sock.sendMessage(
        from,
        {
          text: "No such command, Baka! Have you never seen someone use the command *!help*.",
        },
        { quoted: msg }
      );

    // const now = Date.now();
    // const timestamps = this.cooldown.get(from);
    // const cdAmount = (timestamps || 5) * 1000;
    // if (timestamps.has(from)) {
    //   const expiration = timestamps.get(from) + cdAmount;

    //   if (now < expiration) {
    //     if (isGroup) {
    //       const timeLeft = (expiration - now) / 1000;
    //       // printSpam(isGroup, sender, gcName);
    //       return void this.sock.sendMessage(
    //         from,
    //         {
    //           text: `This group is on cooldown, please wait another _${timeLeft.toFixed(
    //             1
    //           )} second(s)_`,
    //         },
    //         { quoted: msg }
    //       );
    //     } else if (!isGroup) {
    //       const timeLeft = (expiration - now) / 1000;
    //       // printSpam(isGroup, sender);
    //       return void this.sock.sendMessage(
    //         from,
    //         {
    //           text: `You are on cooldown, please wait another _${timeLeft.toFixed(
    //             1
    //           )} second(s)_`,
    //         },
    //         { quoted: msg }
    //       );
    //     }
    //   }
    // }
    // timestamps.set(from, now);
    // setTimeout(() => timestamps.delete(from), cdAmount);

    try {
      await command.run(msg, args);
    } catch (e) {
      console.error(e);
    }
  }

  protected simplified(M: proto.IWebMessageInfo): ISimplifiedMessage {
    const msg = M as ISimplifiedMessage;
    if (
      msg?.message?.protocolMessage ||
      msg?.message?.senderKeyDistributionMessage
    ) {
      return { ...msg, type: Object.keys(msg.message)[0] as MsgType };
    } else {
      if (msg.key) {
        msg.id = msg.key.id;
        msg.fromMe = msg.key.fromMe;
        msg.from = msg.key.remoteJid || "";
        msg.sender = msg.fromMe
          ? this.sock.user.id.split(":")[0] + "@s.whatsapp.net" ||
            this.sock.user.id
          : msg.key.participant || msg.key.remoteJid;
        msg.isGroup = msg.from?.endsWith("@g.us");
      }
      if (msg.message) {
        msg.type =
          Object.keys(msg.message)[0] === "messageContextInfo"
            ? (Object.keys(msg.message)[1] as MsgType)
            : (Object.keys(msg.message)[0] as MsgType);
        msg.body =
          msg.message.conversation ||
          msg.message[msg.type as MsgType.text] ||
          msg.message[msg.type as MsgType.video | MsgType.image]?.caption ||
          (msg.type === "listResponseMessage" &&
            msg.message[msg.type as MsgType.listResponse]?.singleSelectReply
              ?.selectedRowId) ||
          (msg.type === "buttonsResponseMessage" &&
            msg.message[
              msg.type as MsgType.buttonsResponse
            ]?.selectedButtonId?.includes("ASU") &&
            msg.message[msg.type as MsgType.buttonsResponse]
              ?.selectedButtonId) ||
          (msg.type === "templateButtonReplyMessage" &&
            msg.message[msg.type as MsgType.templateButtonReply]?.selectedId) ||
          "";
        if (msg.type === "ephemeralMessage") {
          msg.message = msg.message[msg.type as MsgType.ephemeral]?.message;
          this.simplified(msg);
        }
        if (msg.type === "viewOnceMessage") {
          msg.message = msg.message?.[msg.type as MsgType.viewOnce]?.message;
          this.simplified(msg);
        }
        // @ts-expect-error
        msg.mentions = msg.message[msg.type].contextInfo
          ? // @ts-expect-error
            msg.message[msg.type].contextInfo.mentionedJid
          : null;
        try {
          // @ts-expect-error
          msg.quoted = msg.message[msg.type].contextInfo
            ? // @ts-expect-error
              msg.message[msg.type].contextInfo
            : null;
          if (msg.quoted) {
            const type = Object.keys(
              msg.quoted.quotedMessage as proto.IMessage
            )[0];
            if (type === "ephemeralMessage") {
              const tipe = Object.keys(
                // @ts-expect-error
                msg.quoted.quotedMessage[type].message
              )[0];
              msg.quoted = {
                ...msg.quoted,
                // @ts-expect-error
                message: msg.quoted.quotedMessage[type].message,
                type: "ephemeral",
              };
              if (tipe === "viewOnceMessage") {
                msg.quoted = {
                  ...msg.quoted,
                  // @ts-expect-error
                  message: msg.quoted.quotedMessage[type].message[tipe].message,
                  type: "view_once",
                };
              }
            }
            if (type === "viewOnceMessage") {
              msg.quoted = {
                ...msg.quoted,
                // @ts-expect-error
                message: msg.quoted.quotedMessage[type].message,
                type: "view_once",
              };
            }
            // @ts-expect-error
            msg.quoted.type = msg.quoted.type || "normal";
            // @ts-expect-error
            msg.quoted.message = {
              // @ts-expect-error
              ...(msg.quoted.message || msg.quoted.quotedMessage),
            };
            // @ts-expect-error
            msg.quoted.fromMe =
              // @ts-expect-error
              msg.quoted.participant ===
              (this.sock.user &&
                this.sock.user.id.split(":")[0] + "@s.whatsapp.net");
            // @ts-expect-error
            msg.quoted.key = {
              remoteJid: msg.from,
              // @ts-expect-error
              fromMe: msg.quoted.fromMe,
              // @ts-expect-error
              id: msg.quoted.stanzaId,
            };
            // @ts-expect-error
            msg.quoted.delete = () =>
              // @ts-expect-error
              this.sock.sendMessage(msg.from, { delete: msg.quoted.key });
            // @ts-expect-error
            msg.quoted.download = (pathFile) =>
              // @ts-expect-error
              downloadMedia(msg.quoted.message, pathFile);
            // @ts-expect-error
            delete msg.quoted.quotedMessage;
          }
        } catch {
          // @ts-expect-error
          msg.quoted = null;
        }
        msg.reply = (text) =>
          this.sock.sendMessage(msg.from, { text }, { quoted: msg });
        msg.download = (pathFile) =>
          // @ts-expect-error
          downloadMedia(msg.message, pathFile);
      }
      return msg;
    }
  }
}
