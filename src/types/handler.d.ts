import { proto, WASocket } from "@adiwajshing/baileys-md";

export enum MsgType {
  text = "conversation",
  extendedText = "extendedTextMessage",
  contact = "contactMessage",
  contactsArray = "contactsArrayMessage",
  groupInviteMessage = "groupInviteMessage",
  listMessage = "listMessage",
  buttonsMessage = "buttonsMessage",
  location = "locationMessage",
  liveLocation = "liveLocationMessage",
  image = "imageMessage",
  video = "videoMessage",
  sticker = "stickerMessage",
  document = "documentMessage",
  audio = "audioMessage",
  product = "productMessage",
  listResponse = "listResponseMessage",
  buttonsResponse = "buttonsResponseMessage",
  templateButtonReply = "templateButtonReplyMessage",
  ephemeral = "ephemeralMessage",
  viewOnce = "viewOnceMessage",
}
export interface ISock extends WASocket {
  reply: (msg: string, M: ISimplifiedMessage) => void;
  groupQueryInvite: (code: string) => void;
}
export interface IExtendedGroupMetadata extends WAGroupMetadata {
  admins?: string[];
}

export interface ISimplifiedMessage extends proto.IWebMessageInfo {
  id?: string | null;
  fromMe?: boolean | null;
  from: string;
  sender?: string | null;
  type: MsgType;
  isGroup?: boolean | null;
  body: string;
  mentions?: string[] | null;
  quoted?: proto.IContextInfo & {
    type: string;
    message: proto.IMessage;
    fromMe: boolean;
    key: {
      remoteJid: string;
      fromMe: string;
      id: string;
    };
    delete: () => void;
    download: (pathFile: string) => Promise<string | Buffer>;
  };
  reply: (text: string) => void;
  download: (pathFile: string) => Promise<string | Buffer>;
}
