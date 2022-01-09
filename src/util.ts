/* eslint-disable @typescript-eslint/ban-ts-comment */
import chalk from "chalk";
import { readdirSync, statSync } from "fs";
import { join } from "path";
import { proto, downloadContentFromMessage } from "@adiwajshing/baileys-md";
import { promises as fs } from "fs";
import pino from "pino";
import util from "util";

export default class {
  static isModuleExist = (name: string): boolean => {
    try {
      require.resolve(name);
    } catch (err: any) {
      if (err.code === "MODULE_NOT_FOUND") {
        return false;
      }
    }
    return true;
  };

  static readDirRecursive = (dir: string): string[] => {
    const results: string[] = [];
    const read = (path: string): void => {
      const files = readdirSync(path);
      for (const file of files) {
        const dir = join(path, file);
        if (statSync(dir).isDirectory()) read(dir);
        else results.push(dir);
      }
    };
    read(dir);
    return results;
  };

  static applyMixins = (derivedCtor: any, constructors: any[]) => {
    constructors.forEach((baseCtor) => {
      Object.getOwnPropertyNames(baseCtor.prototype).forEach((name) => {
        Object.defineProperty(
          derivedCtor.prototype,
          name,
          Object.getOwnPropertyDescriptor(baseCtor.prototype, name) ||
            Object.create(null)
        );
      });
    });
  };

  static downloadMedia = async (message: proto.IMessage, pathFile: string) => {
    if (pathFile) {
      const type = Object.keys(message)[0];
      const mimeMap = {
        imageMessage: "image",
        videoMessage: "video",
        stickerMessage: "sticker",
        documentMessage: "document",
        audioMessage: "audio",
      };
      const stream = await downloadContentFromMessage(
        // @ts-expect-error
        message[type],
        // @ts-expect-error
        mimeMap[type]
      );
      let buffer = Buffer.from([]);
      for await (const chunk of stream) {
        buffer = Buffer.concat([buffer, chunk]);
      }
      await fs.writeFile(pathFile, buffer);
      return pathFile;
    } else {
      const type = Object.keys(message)[0];
      const mimeMap = {
        imageMessage: "image",
        videoMessage: "video",
        stickerMessage: "sticker",
        documentMessage: "document",
        audioMessage: "audio",
      };
      const stream = await downloadContentFromMessage(
        // @ts-expect-error
        message[type],
        // @ts-expect-error
        mimeMap[type]
      );
      let buffer = Buffer.from([]);
      for await (const chunk of stream) {
        buffer = Buffer.concat([buffer, chunk]);
      }
      return buffer;
    }
  };

  static roxyLog = <T>(code: pino.Level, args: T, error?: any) => {
    const icon =
      code === "info"
        ? chalk.blue("ℹ")
        : code === "error"
        ? chalk.red("✖")
        : code === "warn"
        ? chalk.yellow("⚠")
        : code === "fatal"
        ? chalk.red("✖")
        : chalk.green("?");
    if (typeof args === "object" && args !== null) {
      const formated = util.inspect(args, {
        showHidden: true,
        depth: 1,
        colors: true,
      });
      if (code === "fatal")
        return console.log(
          `[${icon}] [${chalk.red("FATAL")}] [${chalk.green(
            "OBJECT"
          )}]\n${formated}\n${error.stack ? error.stack : error}`
        );
      return void console.log(
        `[${icon}] [${chalk.blue(
          new Date().toLocaleTimeString()
        )}] [${chalk.green("OBJECT")}]\n${formated}`
      );
    } else {
      if (code === "fatal")
        return console.log(
          `[${icon}] [${chalk.red("FATAL")}] ${chalk.yellowBright(args)}\n${
            error.stack ? error.stack : error
          }`
        );
      return void console.log(
        `[${icon}] [${chalk.blueBright(
          new Date().toLocaleTimeString()
        )}] ${chalk.green(args)}`
      );
    }
  };
}
