/* eslint-disable @typescript-eslint/no-unused-vars */
import { IDB, iSession, iSettings, promiseOrUndefined } from "../types/db";
import { generateSettings } from "../model/gen";

class DB implements IDB {
  public settings: Map<string, iSettings> = new Map();
  public session: Map<string, iSession> = new Map();
  public _ctx: any;
  constructor(public master: string) {
    this.__init();
  }

  __init = () => {};
  _getCustomSettings = (master: string): promiseOrUndefined<iSettings> => {};
  _setCustomSettings = (settings: iSettings, master: string): void => {};

  get _settings(): iSettings {
    // when get the settings
    return generateSettings("sample");
  }
  set _settings(settings: iSettings) {
    // when set the settings
  }
}

export const init = (master: string) => new DB(master);
