import type { Session, SessionData } from "@hono/session";

type NoAuthProps = {
  session: Session<SessionData>;
};

export const notAuth = async (
  props: NoAuthProps,
  next: () => Promise<void> | void,
) => {
  const session = await props.session.get();

  if (!session?.uid) {
    return next();
  }

  return new Response(null, { status: 302, headers: { location: "/" } });
};
