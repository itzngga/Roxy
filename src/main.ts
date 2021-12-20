require("dotenv").config({ path: ".env" });

import client from "./client/Client";
import SettingService from "./services/settings";

SettingService.__init();
