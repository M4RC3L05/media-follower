import type { NextFunction, RequestContext } from "@remix-run/fetch-router";
import { createRedirectResponse } from "@remix-run/response/redirect";
import { routes } from "../router.ts";

export const notAuth = (ctx: RequestContext, next: NextFunction) => {
  const authUserId = ctx.session.get("uid") as string;

  if (!authUserId) {
    return next();
  }

  return createRedirectResponse(routes.home.href());
};
