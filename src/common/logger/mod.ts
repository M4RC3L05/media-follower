import logger, { stdSerializers } from "pino";
import { memoize } from "@std/cache";

export const makeLogger = memoize((namespace: string) => {
  return logger({
    name: namespace,
    serializers: {
      error: (value) => stdSerializers.errWithCause(value),
      reason: (value) => stdSerializers.errWithCause(value),
    },
    formatters: {
      level: (label) => {
        return {
          level: label,
        };
      },
    },
  });
}, { getKey: (namespace) => namespace });
