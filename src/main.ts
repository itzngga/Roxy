require("dotenv").config({ path: ".env" });

import { Cmd } from "./handler";
import process from "process";
import client from "./client/Client";
import util from "./util";
// import SettingService from "./services/settings";
Cmd.loadCmd();
new client("628568970792@s.whatsapp.net");
// SettingService.__init();
process.on("uncaughtException", (err) => {
  util.roxyLog("fatal", "uncaughtException", err);
  return process.exit(2);
});
process.on("unhandledRejection", (err) => {
  util.roxyLog("fatal", "unhandledRejection", err);
  return process.exit(2);
});
