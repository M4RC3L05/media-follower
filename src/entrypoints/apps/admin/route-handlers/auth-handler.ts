import type { IDatabase } from "#src/common/database/database.ts";
import { authPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import z from "@zod/zod";
import type { DbUsersTable } from "#src/common/database/types.ts";
import * as passwordCrypto from "#src/common/crypto/mod.ts";
import type { Session, SessionData } from "@hono/session";

export const loginHandlerGet = () => pageToHtmlResponse(authPages.loginPage());

type LoginHandlerPostProps = {
  database: IDatabase;
  formData: FormData;
  session: Session<SessionData>;
};

export const loginHandlerPost = async (props: LoginHandlerPostProps) => {
  const { username, password } = z.object({
    username: z.string(),
    password: z.string(),
  }).parse(
    Object.fromEntries(props.formData.entries()),
  );

  const [user] = props.database.sql.all`
    select * from users
    where username = ${username}
  ` as DbUsersTable[];

  if (!user) {
    return new Response(null, {
      status: 302,
      headers: { location: "/auth/login" },
    });
  }

  const isValid = await passwordCrypto.verify(password, user.password);

  if (!isValid) {
    return new Response(null, {
      status: 302,
      headers: { location: "/auth/login" },
    });
  }

  await props.session.update({ uid: user.id });

  return new Response(null, {
    status: 302,
    headers: { location: "/" },
  });
};

export const registerHandlerGet = () =>
  pageToHtmlResponse(authPages.registerPage());

type RegisterHandlerPostProps = {
  database: IDatabase;
  session: Session<SessionData>;
  formData: FormData;
};

export const registerHandlerPost = async (props: RegisterHandlerPostProps) => {
  const { username, password } = z.object({
    username: z.string(),
    password: z.string(),
  }).parse(
    Object.fromEntries(props.formData!.entries()),
  );

  const [user] = props.database.sql.all`
    select * from users
    where username = ${username}
  ` as DbUsersTable[];

  if (user) {
    return new Response(null, {
      status: 302,
      headers: { location: "/auth/login" },
    });
  }

  const [createdUser] = props.database.sql.all`
    insert into users
      (id, username, password)
    values
      (${crypto.randomUUID()}, ${username}, ${await passwordCrypto.hash(
    password,
  )})
    returning *;
  ` as DbUsersTable[];

  if (!createdUser) {
    return new Response(null, {
      status: 302,
      headers: { location: "/auth/register" },
    });
  }

  return new Response(null, {
    status: 302,
    headers: { location: "/auth/login" },
  });
};

type LogoutHandlerProps = {
  session: Session<SessionData>;
};

export const logoutHandler = (props: LogoutHandlerProps) => {
  props.session.delete();

  return new Response(null, {
    status: 302,
    headers: { location: "/auth/login" },
  });
};
