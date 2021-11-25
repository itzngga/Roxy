import { AnyAuthenticationCredentialsBase64 } from "@adiwajshing/baileys";

type promiseOrUndefined<Type> = Type | Promise<Type | undefined> | void;
type getType<Type> = Type;

interface iSettings {
  masterId: string;
  prefix: string;
  selfMode: boolean;
  admins: string[];
}

interface iSession extends AnyAuthenticationCredentialsBase64 {
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
