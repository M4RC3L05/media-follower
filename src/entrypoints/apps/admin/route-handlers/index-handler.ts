import { rootPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";

export const indexHandler = () => pageToHtmlResponse(rootPages.indexPage());
