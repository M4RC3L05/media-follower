import type { RequestContext } from "@remix-run/fetch-router";
import type { IDatabase } from "#src/common/database/database.ts";
import { authPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import z from "@zod/zod";
import type { DbUsersTable } from "#src/common/database/types.ts";
import { routes } from "../router.ts";
import * as passwordCrypto from "#src/common/crypto/mod.ts";
import { createRedirectResponse } from "@remix-run/response/redirect";

export const loginHandlerGet = () => pageToHtmlResponse(authPages.loginPage());

type LoginHandlerPostProps = {
  database: IDatabase;
};

export const loginHandlerPost =
  (props: LoginHandlerPostProps) => async (ctx: RequestContext) => {
    const { username, password } = z.object({
      username: z.string(),
      password: z.string(),
    }).parse(
      Object.fromEntries(ctx.formData!.entries()),
    );

    const [user] = props.database.sql<DbUsersTable>`
      select * from users
      where username = ${username}
    `;

    if (!user) {
      return createRedirectResponse(routes.auth.login.index.href());
    }

    const isValid = await passwordCrypto.verify(password, user.password);

    if (!isValid) {
      return createRedirectResponse(routes.auth.login.index.href());
    }

    ctx.session.regenerateId();
    ctx.session.set("uid", user.id);

    return createRedirectResponse(routes.home.href());
  };

export const registerHandlerGet = () =>
  pageToHtmlResponse(authPages.registerPage());

type RegisterHandlerPostProps = {
  database: IDatabase;
};

export const registerHandlerPost =
  (props: RegisterHandlerPostProps) => async (ctx: RequestContext) => {
    const { username, password } = z.object({
      username: z.string(),
      password: z.string(),
    }).parse(
      Object.fromEntries(ctx.formData!.entries()),
    );

    const [user] = props.database.sql<DbUsersTable>`
      select * from users
      where username = ${username}
    `;

    if (user) {
      return createRedirectResponse(routes.auth.login.index.href());
    }

    const [createdUser] = props.database.sql<DbUsersTable>`
      insert into users
        (id, username, password)
      values
        (${crypto.randomUUID()}, ${username}, ${await passwordCrypto.hash(
      password,
    )})
      returning *;
    `;

    if (!createdUser) {
      return createRedirectResponse(routes.auth.register.index.href());
    }

    return createRedirectResponse(routes.auth.login.index.href());
  };

export const logoutHandler = (ctx: RequestContext) => {
  ctx.session.destroy();

  return createRedirectResponse(routes.auth.login.index.href());
};
