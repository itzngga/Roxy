/* eslint-disable @typescript-eslint/ban-ts-comment */
import ffmpeg from "fluent-ffmpeg";
import fs = require("fs");
import path from "path";
import util from "../../util";
import concat from "concat-stream";
import Scrap from "../../scrap";
import axios from "axios";
import { on } from "process";
const ffmpegPath = require("@ffmpeg-installer/ffmpeg").path;
ffmpeg.setFfmpegPath(ffmpegPath);
export default class Sticker {
  protected buffer: Buffer;
  protected operation: "convert" | "toImage" | "toVideo" = "convert";
  protected done = false;
  protected opt: {
    packInfo?: {
      packname?: string;
      author?: string;
    };
    type: "image" | "video";
  } = { type: "image" };

  constructor(file: Buffer | string = "") {
    if (file instanceof Buffer) {
      this.buffer = file;
    } else {
      this.buffer = fs.readFileSync(file);
    }
  }

  public setPackInfo(opt: {
    packInfo?: {
      packname?: string;
      author?: string;
    };
    type: "image" | "video";
  }) {
    this.operation = "convert";
    this.opt = opt;
    return this;
  }

  public toImage() {
    this.opt.type = "image";
    this.operation = "toImage";
    return this;
  }

  public toVideo() {
    this.opt.type = "video";
    this.operation = "toVideo";
    return this;
  }

  public async build(): Promise<Sticker> {
    if (this.operation === "convert") {
      if (this.opt.type === "image") {
        return new Promise((resolve, reject) => {
          ffmpeg(util.buffer2Stream(this.buffer))
            .videoCodec("libwebp")
            .addOutputOptions("-fs", "1M")
            .format("webp")
            // .save(
            //   path.join(
            //     "c:\\Users\\itzngga\\Desktop\\Project\\Roxy\\temp\\",
            //     util.randomString() + ".webp"
            //   )
            // )
            .stream(concat((buf: Buffer) => (this.buffer = buf)))
            .on("error", (err: string) => console.error(err))
            .once("finish", () => {
              this.done = true;
              return resolve(this);
            });
        });
      } else {
        return new Promise((resolve, reject) => {
          ffmpeg(util.buffer2Stream(this.buffer))
            .videoCodec("libwebp")
            .addOutputOptions("-fs", "1M")
            .videoFilter(
              "scale=512:512:flags=lanczos:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=#00000000,setsar=1"
            )
            .format("webp")
            .stream(concat((buf: Buffer) => (this.buffer = buf)))
            .on("error", (err: string) => reject(err))
            .once("finish", () => {
              this.done = true;
              return resolve(this);
            });
        });
      }
    } else {
      if (this.opt.type === "image") {
        return new Promise((resolve, reject) => {
          ffmpeg(util.buffer2Stream(this.buffer))
            .toFormat("png")
            .stream(concat((buf: Buffer) => (this.buffer = buf)))
            .once("error", (err: string) => reject(err))
            .once("finish", () => {
              this.done = true;
              return resolve(this);
            });
        });
      } else {
        try {
          const res = await Scrap.webpToMp4(this.buffer);
          try {
            const { data } = await axios.get(res as string, {
              responseType: "arraybuffer",
            });
            this.buffer = data as Buffer;
            this.done = true;
            return await Promise.resolve(this);
          } catch (err_3) {
            return await Promise.reject(err_3);
          }
        } catch (err_4) {
          return await Promise.reject(err_4);
        }
      }
    }
  }

  public getBuffer(): Promise<Buffer> {
    if (!this.done) return Promise.reject("not builded yet");
    return Promise.resolve(this.buffer);
  }

  public getPath() {
    if (!this.done) return Promise.reject("not builded yet");
    const ext =
      this.operation === "convert"
        ? ".webp"
        : this.opt.type === "image"
        ? ".png"
        : ".mp4";
    fs.promises
      .writeFile(
        path.join(
          "c:\\Users\\itzngga\\Desktop\\Project\\Roxy\\temp\\",
          util.randomString() + ext
        ),
        this.buffer
      )
      .then(() => {
        return path.join(
          "c:\\Users\\itzngga\\Desktop\\Project\\Roxy\\temp\\",
          util.randomString() + ext
        );
      });
  }
}
