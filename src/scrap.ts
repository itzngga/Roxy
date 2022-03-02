/* eslint-disable @typescript-eslint/ban-ts-comment */
import axios from "axios";
import cheerio from "cheerio";
import BodyForm from "form-data";
import Util from "./util";
import path from "path";
import fs from "fs";

export default class Scrap {
  static webpToMp4(buff: Buffer) {
    return new Promise((resolve, reject) => {
      const file = path.join(
        __dirname,
        "\\..\\temp",
        Util.randomString() + ".webp"
      );
      fs.promises.writeFile(file, buff).then(() => {
        const form = new BodyForm();
        form.append("new-image-url", "");
        form.append("new-image", fs.createReadStream(file));
        axios({
          method: "post",
          url: "https://s6.ezgif.com/webp-to-mp4",
          data: form,
          headers: {
            // @ts-expect-error
            "Content-Type": `multipart/form-data; boundary=${form._boundary}`,
          },
        })
          .then(({ data }) => {
            const bodyFormThen = new BodyForm();
            const $ = cheerio.load(data);
            const file = $('input[name="file"]').attr("value");
            bodyFormThen.append("file", file);
            bodyFormThen.append("convert", "Convert WebP to MP4!");
            axios({
              method: "post",
              url: "https://ezgif.com/webp-to-mp4/" + file,
              data: bodyFormThen,
              headers: {
                // @ts-expect-error
                "Content-Type": `multipart/form-data; boundary=${bodyFormThen._boundary}`,
              },
            })
              .then(({ data }) => {
                const $ = cheerio.load(data);
                const result =
                  "https:" +
                  $("div#output > p.outfile > video > source").attr("src");
                return resolve(result);
              })
              .catch(reject);
          })
          .catch(reject);
      });
    });
  }
}
