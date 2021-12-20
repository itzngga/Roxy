import { AnyAuthenticationCredentialsBase64 } from "@adiwajshing/baileys";

type promiseOrUndefined<Type> = Type | Promise<Type | undefined> | void;
type getType<Type> = Type;

interface ISettings {
  id: string;
  prefix: string;
  selfMode: boolean;
  admins: string[];
  autoRead: boolean;
  autoReply: boolean;
  fakeReply: boolean;
  mentionReply: boolean;
  fakeReplyText: string;
  fakeReplyJID: string;
  mentionedMsg: string;
}

interface ISession extends AnyAuthenticationCredentialsBase64 {
  sessionId: string;
}

interface IDB {
  master: string;
  _ctx: getType;
  __init(): getType;
  _getCustomSettings(master: string): promiseOrUndefined<ISettings>;
  _setCustomSettings(settings: iSettings, master: string): void;
  get _settings(): getType;
  set _settings(settings: iSettings);
}
