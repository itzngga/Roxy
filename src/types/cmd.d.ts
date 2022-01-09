import { ISimplifiedMessage } from "./handler";
import { MsgHandler, ISock } from "../handler";
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
type Categories = "admins" | "dev" | "any" | "misc";
