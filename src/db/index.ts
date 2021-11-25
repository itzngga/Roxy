import env from "../config";
import { RoxyError } from "../client/util";
import { IDB, iSettings } from "../types/db";

export default class {
  protected DB!: IDB;
  public _settings?: iSettings;
  constructor(public master: string) {
    this._initDB();
  }

  _initDB = () => {
    if (["pg", "mysql", "sqlite"].includes(env.DB_MODE)) {
      return this._importDB("./knex");
    } else if (env.DB_MODE === "json") {
      return this._importDB("./json");
    } else {
      throw new RoxyError(
        `
Invalid DB_MODE
Supported Type:
- json
- mysql
- pg
- sqlite
`
      );
    }
  };

  _importDB = (path: string) =>
    import(path).then(async (imp) => {
      this.DB = imp.init(this.master);
      this._settings = await this.DB._settings;
    });
}
