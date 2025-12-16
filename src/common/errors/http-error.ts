export class HttpError extends Error {
  readonly status: number;
  readonly headers?: Headers;

  constructor(message: string, status: number, headers?: Headers) {
    super(message);

    this.status = status;

    if (headers) {
      this.headers = headers;
    }
  }
}
