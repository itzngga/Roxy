import env from "../config";
import { Sequelize } from "sequelize-typescript";

export default new Sequelize({
  dialect: "mysql",
  host: env.DB_HOST,
  username: env.DB_USER,
  password: env.DB_PASS,
  database: env.DB_NAME,
  repositoryMode: true,
});
