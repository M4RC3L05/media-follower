import type { IServerApp } from "#src/common/server/mod.ts";
import { Hono } from "hono";
import { secureHeaders } from "hono/secure-headers";
import { useSession } from "@hono/session";
import { config } from "#src/common/config/mod.ts";
import { setCookie } from "hono/cookie";
import type { IDatabase } from "#src/common/database/database.ts";
import type { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import { publicIndex } from "./route-handlers/public-handler.ts";
import {
  loginHandlerGet,
  loginHandlerPost,
  logoutHandler,
  registerHandlerGet,
  registerHandlerPost,
} from "./route-handlers/auth-handler.ts";
import favicon from "#src/common/assets/favicon.ico" with { type: "bytes" };
import { HttpError } from "../../../common/errors/mod.ts";
import { indexHandler } from "./route-handlers/index-handler.ts";
import { checkAuth, notAuth } from "./middlewares/mod.tsx";
import { outputsIndex } from "./route-handlers/outputs-handler.ts";
import {
  inputsCreateGet,
  inputsCreatePost,
  inputsIndex,
} from "./route-handlers/inputs-handler.ts";

type AppProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
  onError?: (error: unknown) => Response | Promise<Response>;
};

export class App implements IServerApp {
  #app: Hono;
  #sessionMaxTime = 60 * 60 * 1;
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
    this.#app = new Hono();

    if (props.onError) {
      this.#app.onError((err) => this.#props.onError!(err));
    }

    this.#app.use(
      secureHeaders({
        crossOriginEmbedderPolicy: "credentialless",
        contentSecurityPolicy: {
          defaultSrc: ["'self'"],
          baseUri: ["'self'"],
          childSrc: ["'self'"],
          connectSrc: ["'self'"],
          fontSrc: ["'self'"],
          formAction: ["'self'"],
          frameAncestors: ["'self'"],
          frameSrc: ["'self'"],
          imgSrc: [
            "'self'",
            "https://*.mzstatic.com",
            "https://images.blu-ray.com",
            "https://*.steamstatic.com",
          ],
          manifestSrc: ["'self'"],
          mediaSrc: ["'self'"],
          objectSrc: ["'none'"],
          sandbox: ["allow-same-origin", "allow-forms"],
          scriptSrc: ["'self'"],
          scriptSrcAttr: ["'none'"],
          scriptSrcElem: ["'self'"],
          styleSrc: ["'self'"],
          styleSrcAttr: ["none"],
          styleSrcElem: ["'self'"],
          upgradeInsecureRequests: [],
          workerSrc: ["'self'"],
        },
      }),
      useSession({
        duration: { absolute: this.#sessionMaxTime },
        secret: config().session.secret,
        setCookie: (c, name, value, opt) =>
          setCookie(c, name, value, {
            ...opt,
            secure: Deno.env.get("ENV") === "production",
            httpOnly: true,
            sameSite: "Strict",
          }),
      }),
    )
      .get(
        "/favicon.ico",
        () =>
          new Response(favicon, {
            headers: { "content-type": "image/x-icon" },
          }),
      )
      .get("/public", (c) => publicIndex({ url: new URL(c.req.url) }))
      .get(
        "/auth/login",
        (c, next) => notAuth({ session: c.var.session }, next),
        () => loginHandlerGet(),
      )
      .post(
        "/auth/login",
        (c, next) => notAuth({ session: c.var.session }, next),
        async (c) =>
          loginHandlerPost({
            database: this.#props.database,
            formData: await c.req.formData(),
            session: c.var.session,
          }),
      )
      .get(
        "/auth/register",
        (c, next) => notAuth({ session: c.var.session }, next),
        () => registerHandlerGet(),
      )
      .post(
        "/auth/register",
        (c, next) => notAuth({ session: c.var.session }, next),
        async (c) =>
          registerHandlerPost({
            database: this.#props.database,
            formData: await c.req.formData(),
            session: c.var.session,
          }),
      )
      .post(
        "/auth/logout",
        (c) => logoutHandler({ session: c.var.session }),
      )
      .get(
        "/inputs",
        (c, next) => checkAuth({ session: c.var.session }, next),
        (c) =>
          inputsIndex({
            database: this.#props.database,
            providers: this.#props.providers,
            url: new URL(c.req.url),
          }),
      )
      .get(
        "/inputs/create",
        (c, next) => checkAuth({ session: c.var.session }, next),
        () => inputsCreateGet(),
      )
      .post(
        "/inputs/create",
        (c, next) => checkAuth({ session: c.var.session }, next),
        async (c) =>
          inputsCreatePost({
            database: this.#props.database,
            formData: await c.req.formData(),
            providers: this.#props.providers,
          }),
      )
      .get(
        "/outputs",
        (c, next) => checkAuth({ session: c.var.session }, next),
        (c) =>
          outputsIndex({
            database: this.#props.database,
            providers: this.#props.providers,
            url: new URL(c.req.url),
          }),
      )
      .get(
        "/",
        (c, next) => checkAuth({ session: c.var.session }, next),
        () => indexHandler(),
      )
      .all("*", () => {
        throw new HttpError("Page not found", 404);
      });
  }

  fetch = (request: Request) => this.#app.fetch(request);
}
