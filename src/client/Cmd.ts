import { ICmd } from "../types/cmd";

export default class Cmd implements ICmd {
  constructor(public wac: any, public config: ICmd["config"]) {}
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  run(Message: any, args: any[]) {
    throw new Error("Method not implemented.");
  }
}
