require("dotenv").config({ path: ".env" });

import db from "./db";

console.log(new db("master"));
