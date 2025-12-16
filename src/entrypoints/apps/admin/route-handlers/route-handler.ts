import type { METHOD } from "@std/http/unstable-method";
import { STATUS_CODE, STATUS_TEXT } from "@std/http";
import type { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import { HttpError } from "#src/common/errors/http-error.ts";

export type SupportedHttpMethods =
  | typeof METHOD.Get
  | typeof METHOD.Put
  | typeof METHOD.Patch
  | typeof METHOD.Delete
  | typeof METHOD.Post;

export interface IRouteHandler extends
  Record<
    SupportedHttpMethods,
    (request: Request) => Response | Promise<Response>
  > {
}

export type AbstractRouteHandlerProps = {
  providers: Record<EInputProvider, IProvider>;
  database: IDatabase;
};

export class AbstractRouteHandler implements IRouteHandler {
  static readonly PATH: string;

  constructor(protected props: AbstractRouteHandlerProps) {}

  DELETE(_request: Request): Response | Promise<Response> {
    throw new HttpError(
      STATUS_TEXT[STATUS_CODE.MethodNotAllowed],
      STATUS_CODE.MethodNotAllowed,
    );
  }

  GET(_request: Request): Response | Promise<Response> {
    throw new HttpError(
      STATUS_TEXT[STATUS_CODE.MethodNotAllowed],
      STATUS_CODE.MethodNotAllowed,
    );
  }

  PATCH(_request: Request): Response | Promise<Response> {
    throw new HttpError(
      STATUS_TEXT[STATUS_CODE.MethodNotAllowed],
      STATUS_CODE.MethodNotAllowed,
    );
  }

  PUT(_request: Request): Response | Promise<Response> {
    throw new HttpError(
      STATUS_TEXT[STATUS_CODE.MethodNotAllowed],
      STATUS_CODE.MethodNotAllowed,
    );
  }

  POST(_request: Request): Response | Promise<Response> {
    throw new HttpError(
      STATUS_TEXT[STATUS_CODE.MethodNotAllowed],
      STATUS_CODE.MethodNotAllowed,
    );
  }
}
