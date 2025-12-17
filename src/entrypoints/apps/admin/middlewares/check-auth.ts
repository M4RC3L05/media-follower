import { createRedirectResponse } from "@remix-run/response/redirect";
import { routes } from "../router.ts";
import type { NextFunction, RequestContext } from "@remix-run/fetch-router";

export const checkAuth = (ctx: RequestContext, next: NextFunction) => {
  const authUserId = ctx.session.get("uid") as string;

  if (!authUserId) {
    return createRedirectResponse(routes.auth.login.index.href());
  }

  return next();
};
