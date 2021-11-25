import { env } from "../types/constant";

const env: env = {
  DB_MODE: <string>(<string>process.env.DB_MODE),
  DB_HOST: <string>(<string>process.env.DB_HOST),
  DB_PORT: <number>parseInt(<string>process.env.DB_PORT),
  DB_USER: <string>(<string>process.env.DB_USER),
  DB_PASS: <string>(<string>process.env.DB_PASS),
  DB_NAME: <string>(<string>process.env.DB_NAME),
  SQLITE_FILENAME: <string>(<string>process.env.SQLITE_FILENAME),
};

export default env;
