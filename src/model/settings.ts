import { Table, Column, Model, PrimaryKey } from "sequelize-typescript";
import { ISettings } from "./types";

@Table
export default class Settings extends Model<ISettings> {
  @PrimaryKey
  @Column
  id: string;
  @Column
  prefix: string;
  @Column
  selfMode: boolean;
  @Column
  admins: string[];
  @Column
  autoRead: boolean;
  @Column
  autoReply: boolean;
  @Column
  fakeReply: boolean;
  @Column
  mentionReply: boolean;
  @Column
  fakeReplyText: string;
  @Column
  fakeReplyJID: string;
  @Column
  mentionedMsg: string;
}
