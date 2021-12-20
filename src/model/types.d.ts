export interface ISettings {
  id: string;
  prefix: string;
  selfMode: boolean;
  admins: string[];
  autoRead: boolean;
  autoReply: boolean;
  fakeReply: boolean;
  mentionReply: boolean;
  fakeReplyText: string;
  fakeReplyJID: string;
  mentionedMsg: string;
}
