import { readdirSync, statSync } from "fs";
import { join } from "path";

export class RoxyError extends Error {
  status?: number;
  context: any;

  constructor(message: string, stack?: string) {
    super(message);
    this.name = this.constructor.name;
    this.context = 404;
    if (stack) {
      this.stack = stack;
    }
  }
}

export const isModuleExist = (name: string): boolean => {
  try {
    require.resolve(name);
  } catch (err: any) {
    if (err.code === "MODULE_NOT_FOUND") {
      return false;
    }
  }
  return true;
};

export const readDirRecursive = (dir: string): string[] => {
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

export function applyMixins(derivedCtor: any, constructors: any[]) {
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
}
