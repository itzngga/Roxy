import SettingService from "../services/settings";
import { WAConnection } from "@adiwajshing/baileys";
import { ISettings } from "../model/types";

export default class Wac extends WAConnection {
  public settings: ISettings;
  constructor(public id: string) {
    super();
    this.connectOptions.connectCooldownMs = 15 * 1000;
    this.connectOptions.alwaysUseTakeover = true;
    this.connectOptions.queryChatsTillReceived = true;
    this.browserDescription = ["Roxy", "Safari", "10.0"];

    this.__generateSettings(id);
  }

  private __generateSettings = async (id: string) => {
    const res = await SettingService.getOneSetting(id);
    if (!res) return console.error("Could not find setting with id: " + id);
    this.settings = res as ISettings;
  };

  public setMessageHandler = (handler: any) => {
    this.removeAllListeners("chat-update");
    return this.on("chat-update", (chat) => handler.run(chat));
  };
}
