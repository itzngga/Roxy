import { JsonDB } from "node-json-db";
import { iSettings, IDB } from "../types/db";
import { generateSettings } from "../model/gen";
import { isModuleExist, RoxyError } from "../client/util";

class DB implements IDB {
  public settings: Map<string, iSettings> = new Map();
  public _ctx!: JsonDB;
  constructor(public master: string) {
    if (!isModuleExist("node-json-db"))
      throw new RoxyError("knex module not found");
    this.__init();
  }

  __init = () => {
    return (this._ctx = new JsonDB("json/db.json")), true, false;
  };

  get _settings() {
    if (this.settings.has(this.master)) {
      return this.settings.get(this.master) || generateSettings(this.master);
    } else {
      const setting = generateSettings(this.master);
      this._settings = setting;
      return setting;
    }
  }

  set _settings(settings: iSettings) {
    this.settings.set(this.master, settings);
    this._ctx.push("/settings", Array.from(Object.entries(this._settings)));
  }

  _getCustomSettings = (master: string) => {
    return this.settings.get(master);
  };

  _setCustomSettings = (settings: iSettings, master: string) => {
    return this.settings.set(master, settings);
  };
}

export const init = (master: string) => new DB(master);
