import {
  type FunctionComponent,
  isValidElement,
  toChildArray,
  type VNode,
} from "preact";
import { renderToString } from "preact-render-to-string";
import { defaults, nosecone } from "nosecone";

export const Page: FunctionComponent & {
  Head: FunctionComponent;
  Body: FunctionComponent;
} = ({ children }) => {
  const childArray = toChildArray(children);
  const head = childArray.find((item) =>
    isValidElement(item) && item.type === Head
  );
  const body = childArray.find((item) =>
    isValidElement(item) && item.type === Body
  );

  return (
    <html lang="en">
      {head}
      {body}
    </html>
  );
};

const Head: FunctionComponent = ({ children }) => (
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="stylesheet" href="/public/styles.css" />
    {children}
  </head>
);
Page.Head = Head;

const Body: FunctionComponent = ({ children }) => (
  <body>
    {children}
  </body>
);

Page.Body = Body;

export const pageToHtmlResponse = (
  page: VNode,
  status?: number,
  headers?: Headers,
) => {
  const secureHeaders = nosecone({
    ...defaults,
    crossOriginEmbedderPolicy: { policy: "credentialless" },
    contentSecurityPolicy: {
      directives: {
        ...defaults.contentSecurityPolicy.directives,
        imgSrc: [
          "'self'",
          "https://*.mzstatic.com",
          "https://images.blu-ray.com",
          "https://*.steamstatic.com",
        ],
      },
    },
  });
  const innerHeaders = new Headers(headers);
  innerHeaders.set("content-type", "text/html");

  return new Response(`<!DOCTYPE html>${renderToString(page)}`, {
    status: status ?? 200,
    headers: {
      ...Object.fromEntries(secureHeaders.entries()),
      ...Object.fromEntries(innerHeaders.entries()),
    },
  });
};
