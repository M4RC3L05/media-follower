import { createRouter, form, get, post, route } from "@remix-run/fetch-router";
import { createCookieSessionStorage } from "@remix-run/session/cookie-storage";
import { createCookie } from "@remix-run/cookie";
import { session } from "@remix-run/session-middleware";
import { formData } from "@remix-run/form-data-middleware";
import { staticFiles } from "@remix-run/static-middleware";
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
import {
  loginHandlerGet,
  loginHandlerPost,
  logoutHandler,
  registerHandlerGet,
  registerHandlerPost,
} from "./route-handlers/auth-handler.ts";
import { config } from "#src/common/config/mod.ts";
import { checkAuth, notAuth } from "./middlewares/mod.tsx";
import { HttpError } from "#src/common/errors/mod.ts";

export const routes = route({
  home: get("/"),
  public: get("/public"),
  inputs: {
    index: get("/inputs"),
    create: form("/inputs/create"),
  },
  outputs: get("/outputs"),
  auth: {
    login: form("/auth/login"),
    register: form("/auth/register"),
    logout: post("/auth/logout"),
  },
});

type MakeRouterProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
};

export const makeRouter = (props: MakeRouterProps) => {
  const sessionCookie = createCookie("__sid", {
    secrets: [config().session.secret],
    httpOnly: true,
    secure: Deno.env.get("ENV") === "production",
    sameSite: "Strict",
  });

  const sessionStorage = createCookieSessionStorage();

  const router = createRouter({
    middleware: [
      staticFiles("src/entrypoints/apps/admin/public"),
      session(sessionCookie, sessionStorage),
      formData(),
    ],
  });

  router.map(routes, {
    public: publicIndex,
    auth: {
      login: {
        middleware: [notAuth],
        actions: {
          action: loginHandlerPost({ database: props.database }),
          index: loginHandlerGet,
        },
      },
      register: {
        middleware: [notAuth],
        actions: {
          index: registerHandlerGet,
          action: registerHandlerPost({ database: props.database }),
        },
      },
      logout: { middleware: [checkAuth], action: logoutHandler },
    },
    home: { action: indexHandler, middleware: [checkAuth] },
    outputs: {
      middleware: [checkAuth],
      action: outputsIndex(props),
    },
    inputs: {
      middleware: [checkAuth],
      actions: {
        create: {
          action: inputsCreatePost(props),
          index: inputsCreateGet,
        },
        index: inputsIndex(props),
      },
    },
  });

  router.route("ANY", "*", () => {
    throw new HttpError("Page not found", 404);
  });

  return router;
};
