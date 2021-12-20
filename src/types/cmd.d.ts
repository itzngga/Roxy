export interface ICmd {
  wac: any;
  run(Message: ParsedMessage, args: ParsedArgs[]);
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
