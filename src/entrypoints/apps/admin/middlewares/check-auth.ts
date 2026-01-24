import type { Session, SessionData } from "@hono/session";

type CheckAuthPorps = {
  session: Session<SessionData>;
};

export const checkAuth = async (
  props: CheckAuthPorps,
  next: () => Promise<void> | void,
) => {
  const session = await props.session.get();

  if (!session?.uid) {
    return new Response(null, {
      status: 302,
      headers: { location: "/auth/login" },
    });
  }

  return next();
};
