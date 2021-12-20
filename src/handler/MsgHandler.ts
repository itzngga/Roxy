import Wac from "../client/Wac";
import { join } from "path";
import { ICmd } from "../types/cmd";
import { readDirRecursive } from "../client/util";
import { WAChatUpdate, proto, MessageType } from "@adiwajshing/baileys";
import { IExtendedGroupMetadata, ISimplifiedMessage } from "../types/handler";

export class MsgHandler {
  public cmdMap: Map<string, ICmd> = new Map();
  public aliasMap: Map<string, ICmd> = new Map();

  constructor(public wac: Wac) {
    this._loadCmd();
  }

  __initHandler(): void {
    return this._loadCmd();
  }

  public run(chat: WAChatUpdate): void {
    if (!chat.hasNewMessage || typeof chat.messages === "undefined") return;
    const message = chat.messages.first;
    if (!message.message) return;
    if (message.key && message.key.remoteJid === "status@broadcast") return;
    const parsedMessage = this.parseRawMessage(message);
  }

  private parseRawMessage(message: proto.WebMessageInfo) {}

  private _loadCmd = () => {
    const files = readDirRecursive(join(__dirname, "\\..\\cmd"));
    for (const path of files) {
      const cmd: ICmd = new (require(path).default)(this.wac, this);
      this.cmdMap.set(cmd.config.cmd, cmd);
      if (cmd.config.aliases)
        cmd.config.aliases.forEach((alias) => this.aliasMap.set(alias, cmd));
    }
    console.log("loaded cmd: " + this.cmdMap.size);
  };
}

class ParseRawMessage {
  constructor(protected wac: Wac, protected M: proto.WebMessageInfo) {
    
  }

  simple(M: proto.WebMessageInfo): Promise<ISimplifiedMessage> {
    if (M.message?.ephemeralMessage)
			M.message = M.message.ephemeralMessage.message;
		const jid = M.key.remoteJid || "";
		const chat = jid.endsWith("g.us") ? "group" : "dm";
		const type = (Object.keys(M.message || {})[0] || "") as MessageType;
		const user = chat === "group" ? M.participant : jid;
		const info = this.wac.getContact(user);
		const groupMetadata: IExtendedGroupMetadata | null =
			chat === "group" ? await this.wac.groupMetadata(jid) : null;
		if (groupMetadata)
			groupMetadata.admins = groupMetadata.participants
				.filter((user) => user.isAdmin)
				.map((user) => user.jid);
		const sender = {
			jid: user,
			username: info.notify || info.vname || info.name || "User",
			isAdmin:
				groupMetadata && groupMetadata.admins
					? groupMetadata.admins.includes(user)
					: false,
		};
		const content: string | null =
			type === MessageType.text && M.message?.conversation
				? M.message.conversation
				: this.supportedMediaMessages.includes(type)
				? this.supportedMediaMessages
						.map(
							(type) =>
								M.message?.[type as MessageType.image | MessageType.video]
									?.caption
						)
						.filter((caption) => caption)[0] || ""
				: type === MessageType.extendedText &&
				  M.message?.extendedTextMessage?.text
				? M.message?.extendedTextMessage.text
				: null;
		const quoted: ISimplifiedMessage["quoted"] = {};
		quoted.message = M?.message?.[type as MessageType.extendedText]?.contextInfo
			?.quotedMessage
			? JSON.parse(JSON.stringify(M).replace("quotedM", "m")).message?.[
					type as MessageType.extendedText
			  ].contextInfo
			: null;
		quoted.sender =
			M.message?.[type as MessageType.extendedText]?.contextInfo?.participant ||
			null;
		return {
			type,
			content,
			chat,
			sender,
			quoted,
			args: content?.split(" ") || [],
			reply: async (
				content: string | Buffer,
				type?: MessageType,
				mime?: Mimetype,
				mention?: string[],
				caption?: string,
				thumbnail?: Buffer
			) => {
				const options = {
					quoted: M,
					caption,
					mimetype: mime,
					contextInfo: { mentionedJid: mention },
				};
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				if (thumbnail) (options as any).thumbnail = thumbnail;
				await this.sendMessage(jid, content, type || MessageType.text, options);
			},
			mentioned: this.getMentionedUsers(M, type),
			from: jid,
			groupMetadata,
			WAMessage: M,
			urls: this.util.getUrls(content || ""),
		};
  }
  parse() {
    return {
      serial: this.__serial(),
      fromMe: this.__fromMe(),
      chatId: this.__chatId(),
      msgId: this.__msgId(),
      isBaileys: this.__isBaileys(),
      from: this.__from(),
    };
  }
  __serial = () => {
    const jid = this.__fromMe() ? this.wac.user.jid
    : this.message.key.remoteJid?.endsWith("@g.us")
    ? (this.message.participant as string)
    : (this.message.key.remoteJid as string)
    const contact = this.wac.contacts[jid];
    return {
      jid,
      name: contact,
      isAdmin: 
    }
  }
    
  __fromMe = () => this.message.key.fromMe as boolean;
  __chatId = () => this.message.key.remoteJid as string;
  __msgId = () => this.message.key.id as string;
  __isBaileys = () => this.__msgId().startsWith("3EB0");
  __from = () => this.message.key.remoteJid as string;
  __isGroupMsg = () => this.__from().endsWith("g.us");
  __isSadmin = () => this.wac.settings.admins.includes(this.__serial());
  __type = () => Object.keys(this.message.message!)[0];
  __content = () => {
    try {
      return Object.keys(
        this.message.message!.extendedTextMessage!.contextInfo!.quotedMessage!
      )[0];
    } catch (error) {
      return false;
    }
  };
  __body = () =>
    this.__type() === "conversation"
      ? this.message.message!.conversation!
      : this.__type() === "imageMessage"
      ? this.message.message!.imageMessage!.caption!
      : this.__type() === "videoMessage"
      ? this.message.message!.videoMessage!.caption!
      : this.__type() === "extendedTextMessage"
      ? this.message.message!.extendedTextMessage!.text!
      : "";
  __timestamp = () => Number(this.message.messageTimestamp) * 1000;
  __cmd = () =>
    this.__body()[0] === this.wac.settings.prefix ? this.__body() : "";
  __args = () => this.__cmd().split(" ");
  __getGroupMetadata = async () => {
    const groupMetadata = await this.wac.groupMetadata(this.__from());
    const groupAdmins = groupMetadata?.participants
      .filter((a) => a.isAdmin || a.isSuperAdmin)
      .map((a) => a.jid);
    return {
      groupMetadata,
      groupId: groupMetadata?.id,
      groupName: groupMetadata?.subject,
      groupDesc: groupMetadata?.desc,
      groupMembers: groupMetadata?.participants,
      groupAdmins: groupAdmins,
      isGroupAdmins: groupAdmins.includes(this.__serial()),
      isBotGroupAdmins: groupAdmins.includes(this.wac.user.jid),
      isGroupAdminsOnly: groupMetadata.announce,
    };
  };

  quotedMsgObj = () =>
    this.__type() === "extendedTextMessage
      ? JSON.parse(JSON.stringify(this.message).replace("quotedM", "m")).message
          .extendedTextMessage.contextInfo
      : this.message;
  getQuotedText = () =>
    isQuotedText
      ? this.message.message!.extendedTextMessage!.contextInfo!.quotedMessage!
          .conversation
      : cmd;
  fileSize = (quoted: boolean) => {
    const type = this.__type()
        if (quoted) {
          const msg: any = quotedMsgObj().message ?? [type];
          return msg.fileLength.low;
        } else {
          const msg: any = this.message.message ?? [type];
          return msg.fileLength.low;
        }
      };
    pushname = (target?: string) => {
        const targets: string = target || this.__serial();
        const v =
          targets === "0@s.whatsapp.net"
            ? { jid: targets, vname: "WhatsApp" }
            : targets === this.wac.user.jid
            ? this.wac.user
            : this.wac.addContact(targets);
        return (
          v.name ||
          v.vname ||
          v.notify ||
          new libphonenumber("+" + v.jid.replace("@s.whatsapp.net", "")).getNumber(
            "international"
          )
        );
      };
      const mentionedJidList = () => {
        if (!is_undefined(quotedMsgObj().mentionedJid)) {
          return quotedMsgObj().mentionedJid;
        }
        return [];
      };
}

// if (debug && fromMe) console.log(JSON.stringify(message, null, "\n"));
// if (autoRead) client.chatRead(from);

const isQuoted = type === "extendedTextMessage";
const t = <number>message.messageTimestamp * 1000;
const time = moment(t).format("HH:mm:ss");
const cmd: string = body.substr(0, 1) === prefix ? body : "";
const args: any = cmd.split(" ");

const isMedia =
  type === "imageMessage" ||
  type === "videoMessage" ||
  type === "audioMessage" ||
  type === "stickerMessage" ||
  type === "documentMessage";
const isImage = isQuoted
  ? content() === "imageMessage"
  : type === "imageMessage";
const isVideo = isQuoted
  ? content() === "videoMessage"
  : type === "videoMessage";
// const isAudio = isQuoted ? content() === 'audioMessage' : type === 'audioMessage';
// const isSticker = isQuoted ? content() === 'stickerMessage' : type === 'stickerMessage';
// const isDocument = isQuoted ? content() === 'documentMessage' : type === 'documentMessage';
const isQuotedText = content() === "conversation";
const isQuotedImage = content() === "imageMessage";
const isQuotedVideo = content() === "videoMessage";
const isQuotedAudio = content() === "audioMessage";
const isQuotedSticker = content() === "stickerMessage";
const isQuotedDocument = content() === "documentMessage";
const quotedMsgObj = () =>
  isQuoted
    ? JSON.parse(JSON.stringify(message).replace("quotedM", "m")).message
        .extendedTextMessage.contextInfo
    : message;
const getQuotedText = () =>
  isQuotedText
    ? message.message!.extendedTextMessage!.contextInfo!.quotedMessage!
        .conversation
    : cmd;
if (!isBotGroupAdmins && isGroupAdminOnly) return;
const replyMode = setting.fakeReply
  ? client.generateFakeReply(setting.fakeText)
  : message;
const fileSize = (quoted: boolean) => {
  if (quoted) {
    const msg: any = quotedMsgObj().message ?? [type];
    return msg.fileLength.low;
  } else {
    const msg: any = message.message ?? [type];
    return msg.fileLength.low;
  }
};
const pushname = (target?: string) => {
  const targets: string = target || serial;
  const v =
    targets === "0@s.whatsapp.net"
      ? { jid: targets, vname: "WhatsApp" }
      : targets === client.user.jid
      ? client.user
      : client.addContact(targets);
  return (
    v.name ||
    v.vname ||
    v.notify ||
    new libphonenumber("+" + v.jid.replace("@s.whatsapp.net", "")).getNumber(
      "international"
    )
  );
};
const mentionedJidList = () => {
  if (!is_undefined(quotedMsgObj().mentionedJid)) {
    return quotedMsgObj().mentionedJid;
  }
  return [];
};
