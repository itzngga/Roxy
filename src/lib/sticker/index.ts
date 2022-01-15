import ffmpeg from "fluent-ffmpeg";
import fs = require("fs");
import path from "path";
import { exec } from "child_process";
import Exif from "./exif";

export default class Sticker {
  public static async convert(
    file: Buffer,
    ext1: string,
    ext2: string,
    options: string[]
  ) {
    return new Promise((resolve, reject) => {
      const temp = path.join(__dirname, "../temp", Date.now() + "." + ext1),
        out = temp + "." + ext2;
      return fs.promises.writeFile(temp, file).then(() => {
        return ffmpeg(temp)
          .on("start", (cmd: unknown) => {
            console.log(cmd);
          })
          .on("error", (e: unknown) => {
            fs.unlinkSync(temp);
            reject(e);
          })
          .on("end", () => {
            console.log("Finish");
            setTimeout(() => {
              fs.unlinkSync(temp);
              fs.unlinkSync(out);
            }, 2000);
            resolve(fs.readFileSync(out));
          })
          .addOutputOptions(options)
          .toFormat(ext2)
          .save(out);
      });
    });
  }

  public static convert2(
    file: Buffer,
    ext1: string,
    ext2: string,
    options: string[]
  ) {
    return new Promise((resolve, reject) => {
      const temp = path.join(__dirname, "../temp", Date.now() + "." + ext1),
        out = temp + "." + ext2;
      return fs.promises.writeFile(temp, file).then(() => {
        return ffmpeg(temp)
          .on("start", (cmd: unknown) => {
            console.log(cmd);
          })
          .on("error", (e: unknown) => {
            fs.unlinkSync(temp);
            reject(e);
          })
          .on("end", () => {
            console.log("Finish");
            setTimeout(() => {
              fs.unlinkSync(temp);
              fs.unlinkSync(out);
            }, 2000);
            resolve(fs.readFileSync(out));
          })
          .addOutputOptions(options)
          .seekInput("00:00")
          .setDuration("00:05")
          .toFormat(ext2)
          .save(out);
      });
    });
  }

  public static async WAVideo(file: Buffer, ext1: string) {
    return Sticker.convert(file, ext1, "mp4", [
      "-c:a aac",
      "-c:v libx264",
      "-b:a 128K",
      "-ar 44100",
      "-crf 28",
      "-preset slow",
    ] as never);
  }

  public static WAAudio(file: Buffer, ext1: string) {
    return Sticker.convert(file, ext1, "mp3", [
      "-b:a 128K",
      "-ar 44100",
      "-ac 2",
      "-vn",
    ] as never);
  }

  public static WAOpus(file: Buffer, ext1: string) {
    return Sticker.convert(file, ext1, "opus", [
      "-vn",
      "-c:a libopus",
      "-b:a 128K",
      "-vbr on",
      "-compression_level 10",
    ] as never);
  }

  public static sticker(
    file: Buffer,
    opts: {
      cmdType: string;
      withPackInfo: boolean;
      packInfo: {
        packname: string;
        author: string;
      };
      isImage: boolean;
      isVideo: boolean;
    }
  ) {
    if (typeof opts.cmdType === "undefined") opts.cmdType = "1";
    const cmd: { [index: number]: string[] } = {
      1: [
        "-fs 1M",
        "-vcodec",
        "libwebp",
        "-vf",
        "scale=512:512:flags=lanczos:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=#00000000,setsar=1",
      ],
      2: ["-fs 1M", "-vcodec", "libwebp"],
    };
    if (opts.withPackInfo) {
      if (!opts.packInfo)
        throw Error("'packInfo' must be filled when using 'withPackInfo'");
      const ext =
        opts.isImage !== undefined || false
          ? "jpg"
          : opts.isVideo !== undefined || false
          ? "mp4"
          : null;
      return Sticker.stickerWithExif(
        file,
        ext as string,
        opts.packInfo,
        cmd[parseInt(opts.cmdType)]
      );
    }

    if (opts.isImage) {
      return Sticker.convert(file, "jpg", "webp", cmd[parseInt(opts.cmdType)]);
    }
    if (opts.isVideo) {
      return Sticker.convert2(file, "mp4", "webp", cmd[parseInt(opts.cmdType)]);
    }
  }

  public static stickerWithExif(
    file: Buffer,
    ext: string,
    packInfo: {
      packname: string;
      author: string;
    },
    cmd: string[]
  ): Promise<Buffer> {
    return new Promise(() => {
      const { packname, author } = packInfo;
      const filename = Date.now();
      if (ext === "jpg") {
        return Sticker.convert(file, ext, "webp", cmd).then((res) => run(res));
      } else {
        return Sticker.convert2(file, ext, "webp", cmd).then((res) => run(res));
      }
      function run(buff: unknown) {
        Exif.create(
          packname !== undefined || "" ? packname : "Original",
          author !== undefined || "" ? author : "SMH-BOT",
          filename.toString()
        );
        return fs.promises
          .writeFile(`./temp/${filename}.webp`, buff as Buffer)
          .then(() => {
            exec(
              `webpmux -set exif ./temp/${filename}.exif ./temp/${filename}.webp -o ./temp/${filename}.webp`,
              async (err: unknown) => {
                if (err)
                  return (
                    (await Promise.reject(err)) &&
                    (await Promise.all([
                      fs.promises.unlink(`./temp/${filename}.webp`),
                      fs.promises.unlink(`./temp/${filename}.exif`),
                    ]))
                  );
                setTimeout(() => {
                  fs.unlinkSync(`./temp/${filename}.exif`);
                  fs.unlinkSync(`./temp/${filename}.webp`);
                }, 2000);
                Promise.resolve(fs.readFileSync(`./temp/${filename}.webp`));
              }
            );
          });
      }
    });
  }
}
