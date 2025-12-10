import { makeLogger } from "#src/common/logger/mod.ts";

const log = makeLogger("server");

type ServerPorps = {
  app: IServerApp;
  hostname: string;
  port: number;
};

export class Server {
  #server: Deno.HttpServer;

  constructor(props: ServerPorps) {
    log.info("Creating server");

    this.#server = Deno.serve({
      onError: (error) => {
        console.log(error);
        return new Response(`Unknown error:\n${error}`, {
          status: 500,
          headers: { "content-type": "text/plain" },
        });
      },
      hostname: props.hostname,
      port: props.port,
      onListen: ({ hostname, port }) => {
        log.info(`Serving on http://${hostname}:${port}`);
        log.info("Server created successfully");
      },
    }, props.app.fetch);
  }

  async [Symbol.asyncDispose]() {
    log.info("Closing server");

    await this.#server.shutdown();

    log.info("Server closed successfully");
  }

  get finished() {
    return this.#server.finished;
  }
}

export interface IServerApp {
  fetch: (request: Request) => Response | Promise<Response>;
}
