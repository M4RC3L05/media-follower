export interface IHttpFetch {
  fetch<T>(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<T>;
  fetchText(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<string>;
  fetchReadable(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<ReadableStream>;
}

type HttpFetchProps = {
  signal?: AbortSignal;
  timeout?: number;
};

export class HttpFetch implements IHttpFetch {
  #props?: HttpFetchProps;

  constructor(props?: HttpFetchProps) {
    if (props) {
      this.#props = props;
    }
  }

  async #request(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<Response> {
    const signals: AbortSignal[] = [];

    if (init?.signal) signals.push(init.signal);
    if (this.#props?.signal) signals.push(this.#props.signal);
    signals.push(AbortSignal.timeout(this.#props?.timeout ?? 5000));

    return await fetch(path, {
      ...init,
      signal: AbortSignal.any(signals),
    });
  }

  async fetch<T>(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<T> {
    const response = await this.#request(path, init);

    if (!response.ok) throw new Error("Request failed");

    const data: T = await response.json();

    return data;
  }

  async fetchText(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<string> {
    const response = await this.#request(path, init);

    if (!response.ok) throw new Error("Request failed");

    const data = await response.text();

    return data;
  }

  async fetchReadable(
    path: Parameters<typeof fetch>[0],
    init?: Parameters<typeof fetch>[1],
  ): Promise<ReadableStream> {
    const response = await this.#request(path, init);

    if (!response.ok) throw new Error("Request failed");
    if (!response.body) throw new Error("Response has no body");

    return response.body;
  }
}
