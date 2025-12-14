import type { IDatabase } from "../database/database.ts";
import { EInputProvider } from "../database/enums/input-provider.ts";
import type { IHttpFetch } from "../http/mod.ts";
import { BluRayComPhysicalReleasesProvider } from "./blu-ray-com-physical-releases/provider.ts";
import { ItunesMusicReleasesProvider } from "./itunes-music-releases/provider.ts";

type ProviderFactoryProps = { database: IDatabase; httpClient: IHttpFetch };

export const providerFactory = <P extends EInputProvider>(
  p: P,
  props: ProviderFactoryProps,
) => {
  switch (p) {
    case EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE: {
      return new BluRayComPhysicalReleasesProvider({
        database: props.database,
        httpClient: props.httpClient,
      });
    }
    case EInputProvider.ITUNES_MUSIC_RELEASE: {
      return new ItunesMusicReleasesProvider({
        database: props.database,
        httpClient: props.httpClient,
      });
    }
  }
};
