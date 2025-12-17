import { createRouter, form, get, route } from "@remix-run/fetch-router";
import { createCookieSessionStorage } from "@remix-run/session/cookie-storage";
import { createCookie } from "@remix-run/cookie";
import { session } from "@remix-run/session-middleware";
import { formData } from "@remix-run/form-data-middleware";
import { publicIndex } from "./route-handlers/public-handler.ts";
import { indexHandler } from "./route-handlers/index-handler.ts";
import { outputsIndex } from "./route-handlers/outputs-handler.ts";
import {
  inputsCreateGet,
  inputsCreatePost,
  inputsIndex,
} from "./route-handlers/inputs-handler.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import type { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";

export const routes = route({
  home: get("/"),
  public: get("/public"),
  inputs: {
    index: get("/inputs"),
    create: form("/inputs/create"),
  },
  outputs: get("/outputs"),
});

const sessionCookie = createCookie("__sid", {
  secrets: ["s3cr3t"],
  httpOnly: true,
  secure: Deno.env.get("ENV") === "production",
  sameSite: "Strict",
});

const sessionStorage = createCookieSessionStorage();

type MakeRouterProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
};

export const makeRouter = (props: MakeRouterProps) => {
  const router = createRouter({
    middleware: [session(sessionCookie, sessionStorage), formData()],
  });

  router.map(routes, {
    public: publicIndex,
    home: indexHandler,
    outputs: {
      middleware: [],
      action: outputsIndex(props),
    },
    inputs: {
      middleware: [],
      actions: {
        create: {
          action: inputsCreatePost(props),
          index: inputsCreateGet,
        },
        index: inputsIndex(props),
      },
    },
  });

  return router;
};
