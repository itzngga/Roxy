import DB from "../db";
import Wac from "../client/Wac";
import { WAGroupMetadata } from "@adiwajshing/baileys";

export interface iHandler {
  wac: Wac;

  __initHandler(): void;
  __run(): void;
}

export interface IExtendedGroupMetadata extends WAGroupMetadata {
  admins?: string[];
}

export interface ISimplifiedMessage {
  type: MessageType;
  content: string | null;
  args: string[];
  reply(
    content: string | Buffer,
    type?: MessageType,
    mime?: Mimetype,
    mention?: string[],
    caption?: string,
    thumbnail?: Buffer
  ): Promise<unknown>;
  mentioned: string[];
  groupMetadata: IExtendedGroupMetadata | null;
  chat: "group" | "private";
  from: string;
  sender: {
    jid: string;
    username: string;
    isAdmin: boolean;
  };
  quoted?: {
    message?: WAMessage | null;
    sender?: string | null;
  } | null;
  WAMessage: WAMessage;
  urls: string[];
}
