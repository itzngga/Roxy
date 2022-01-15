import { ISimplifiedMessage } from "./handler";
import { MsgHandler, ISock } from "../handler";

export * from "./cmd";
export * from "/constant";
export * from "./db";
export * from "./handler";

export interface ICmd {
  sock: ISock;
  msgHandler: MsgHandler;
  run: (
    msg: ISimplifiedMessage,
    args: ParsedArgs[]
  ) => Promise<void | never> | void | never;
  config: {
    onlyGroupAdmin?: boolean;
    onlyBotAdmin?: boolean;
    aliases?: string[];
    description?: string;
    cmd: string;
    id?: string;
    category?: Categories;
  };
}
export type Categories = "admins" | "dev" | "any" | "misc";
type;
export interface IMid {
  sock: ISock;
  msgHandler: MsgHandler;
  run: (msg: ISimplifiedMessage, args: ParsedArgs[]) => Promise<void | never>;
  config: {
    name: string;
    mode: "before" | "after";
    times?: number;
    cd?: number;
  };
}
