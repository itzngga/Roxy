import db from "../db";
import Settings from "../model/settings";
import { ISettings } from "../model/types";
import { Repository } from "sequelize-typescript";

export default class SettingService {
  public static ctx: Repository<Settings>;

  public static __init() {
    if (!SettingService.ctx) {
      SettingService.ctx = db.getRepository(Settings);
    }
  }

  public static getService() {
    if (!SettingService.ctx) {
      SettingService.ctx = db.getRepository(Settings);
      return SettingService.ctx;
    } else {
      return SettingService.ctx;
    }
  }

  public static getDefaultSettings(id: string): ISettings {
    return {
      id,
      prefix: "!",
      selfMode: false,
      admins: [id],
      autoRead: false,
      autoReply: false,
      fakeReply: false,
      mentionReply: false,
      fakeReplyText: "FakeReply!",
      fakeReplyJID: "0@s.whatsapp.net",
      mentionedMsg: "Hi! im here",
    };
  }

  public static getOneSetting(id: string) {
    return SettingService.ctx.findOne({ where: { id } });
  }

  public static getAllSetting() {
    return SettingService.ctx.findAll();
  }

  public static setSetting(setting: ISettings) {
    SettingService.ctx.create(setting);
  }

  public static updateSetting(setting: ISettings) {
    return SettingService.ctx.update(setting, { where: { id: setting.id } });
  }

  public static deleteSetting(setting: ISettings) {
    return SettingService.ctx.destroy({ where: { id: setting.id } });
  }
}
