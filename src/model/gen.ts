import { iSettings } from "../types/db";

export const generateSettings = (master: string): iSettings => {
  return {
    masterId: master,
    prefix: "!",
    selfMode: false,
    admins: [master],
  };
};
