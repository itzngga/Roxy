export interface ICmd {
  wac: any;
  handler: any;
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
