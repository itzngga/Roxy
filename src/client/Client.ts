import Wac from "./Wac";
import { MsgHandler } from "../handler";

export default class Client {
  public wac!: Wac;
  constructor(id: string) {
    this.wac = new Wac(id);
  }
}
