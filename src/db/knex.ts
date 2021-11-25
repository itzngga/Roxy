import env from "../config";
import { Knex, knex } from "knex";
import { iSettings, IDB } from "../types/db";
import { isModuleExist, RoxyError } from "../client/util";
import { generateSettings } from "../model/gen";

class DB implements IDB {
  public _ctx!: Knex;
  constructor(public master: string) {
    if (!isModuleExist("knex")) throw new RoxyError("knex module not found");
    this.__init();
  }

  _hasConfig(): boolean {
    if (
      (["pg", "mysql"].includes(env.DB_MODE) && env.DB_HOST) ||
      env.DB_PORT ||
      env.DB_USER ||
      env.DB_PASS ||
      env.DB_NAME
    )
      return true;
    if (env.DB_MODE === "sqlite" && env.SQLITE_FILENAME) return true;
    return false;
  }

  __init = () => {
    switch (env.DB_MODE) {
      case "pg":
        if (!isModuleExist("pg")) throw new RoxyError("pg module not found");
        if (!this._hasConfig())
          throw new RoxyError("some config field is missing");
        this._ctx = knex({
          client: "pg",
          connection: {
            host: env.DB_HOST,
            port: env.DB_PORT,
            user: env.DB_USER,
            password: env.DB_PASS,
            database: env.DB_NAME,
          },
        });
        break;
      case "mysql":
        if (!isModuleExist("mysql"))
          throw new RoxyError("mysql module not found");
        if (!this._hasConfig())
          throw new RoxyError("some config field is missing");
        this._ctx = knex({
          client: "mysql",
          connection: {
            host: env.DB_HOST,
            port: env.DB_PORT,
            user: env.DB_USER,
            password: env.DB_PASS,
            database: env.DB_NAME,
          },
        });
        break;
      case "sqlite":
        if (!isModuleExist("sqlite3"))
          throw new RoxyError("sqlite3 module not found");
        if (!this._hasConfig())
          throw new RoxyError("some config field is missing");
        this._ctx = knex({
          client: "sqlite3",
          connection: {
            filename: env.SQLITE_FILENAME,
          },
        });
        break;
      default:
        this._ctx = knex({});
        break;
    }
  };

  get _settings(): any {
    return this._ctx<iSettings>("settings")
      .select("*")
      .where("masterId", this.master)
      .first()
      .then((res) => {
        if (res) return generateSettings(this.master);
        return res;
      });
  }

  set _settings(settings: iSettings) {
    this._ctx<iSettings>("settings")
      .insert(settings)
      .where("masterId", this.master)
      .onConflict("masterId")
      .merge();
  }

  _getCustomSettings = (master: string): Promise<iSettings | undefined> => {
    return this._ctx<iSettings>("settings")
      .select("*")
      .where("masterId", master)
      .first()
      .then((res) => {
        if (res) return generateSettings(master);
        return res;
      });
  };

  _setCustomSettings = (settings: iSettings, master: string) => {
    this._ctx<iSettings>("settings")
      .insert(settings)
      .where("masterId", master)
      .onConflict("masterId")
      .merge();
  };
}

export const init = (master: string) => new DB(master);
