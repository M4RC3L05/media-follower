type ServerPorps = {
  app: IServerApp;
  onListen?: (host: string, post: number) => void;
  onError?: (error: unknown) => Response | Promise<Response>;
  hostname: string;
  port: number;
};

export class Server {
  #server: Deno.HttpServer;

  constructor(props: ServerPorps) {
    const serverConfig: Deno.ServeTcpOptions = {
      hostname: props.hostname,
      port: props.port,
    };

    if (props.onError) {
      serverConfig.onError = props.onError;
    }

    if (props.onListen) {
      serverConfig.onListen = ({ hostname, port }) =>
        props.onListen!(hostname, port);
    }

    this.#server = Deno.serve(serverConfig, props.app.fetch);
  }

  async [Symbol.asyncDispose]() {
    await this.#server.shutdown();
  }

  get finished() {
    return this.#server.finished;
  }
}

export interface IServerApp {
  fetch: (request: Request) => Response | Promise<Response>;
}
