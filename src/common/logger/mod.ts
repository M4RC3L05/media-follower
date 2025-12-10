import {
  type BaseHandler,
  ConsoleHandler,
  type LevelName,
  Logger,
  type LogRecord,
} from "@std/log";
import { memoize } from "@std/cache";
import pineSerializer, { type SerializedError } from "pino-std-serializers";

const formatLogArg = (arg: unknown) => {
  if (
    arg &&
    typeof arg === "object" &&
    (arg as Record<string, unknown>).error instanceof Error
  ) {
    (arg as { error: SerializedError }).error = pineSerializer.errWithCause(
      (arg as { error: Error }).error,
    );
  }

  if (
    arg &&
    typeof arg === "object" &&
    (arg as Record<string, unknown>).reason instanceof Error
  ) {
    (arg as { reason: SerializedError }).reason = pineSerializer.errWithCause(
      (arg as { reason: Error }).reason,
    );
  }

  return arg;
};

const levelToEmoji: Record<LevelName, string | undefined> = {
  ERROR: "❌",
  WARN: "⚠️",
  INFO: "👀",
  DEBUG: "🔍",
  CRITICAL: "🚨",
  NOTSET: undefined,
};

const logFormatter = (
  { args, datetime, levelName, loggerName, msg }: LogRecord,
) => {
  if (Deno.env.get("ENV") !== "production") {
    return `${datetime.toISOString()} [${loggerName}](${
      [
        levelToEmoji[levelName as LevelName]
          ? levelToEmoji[levelName as LevelName]
          : undefined,
        levelName,
      ].filter(Boolean).join(" ")
    }): ${msg}${
      args[0]
        ? `\n${
          Deno.inspect(formatLogArg(args[0]), { colors: true, depth: 1000 })
        }`
        : ""
    }`;
  }

  return JSON.stringify({
    datetime: datetime.toISOString(),
    level: levelName,
    name: loggerName,
    message: msg,
    data: formatLogArg(args[0]),
  });
};

export const makeLogger = memoize((namespace: string) => {
  const handlers: BaseHandler[] = [];

  if (Deno.env.get("ENV") !== "test") {
    handlers.push(
      new ConsoleHandler("INFO", {
        formatter: logFormatter,
        useColors: false,
      }),
    );
  }

  return new Logger(namespace, "INFO", { handlers: handlers });
}, { getKey: (namespace) => namespace });
