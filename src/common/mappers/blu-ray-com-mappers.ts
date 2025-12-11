import { ReleaseSourceProvider } from "../database/enums/release-source-provider.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../database/types.ts";
import {
  type BluRayComCountry,
  bluRayComCountrySchema,
  type BluRayComReleaseWithExtra,
  bluRayComReleaseWithExtraSchema,
} from "../services/service.ts";

export const fromReleaseToPersistance = (
  item: BluRayComReleaseWithExtra,
): DbReleasesTable => {
  return {
    releasedAt: item.releasedate.toISOString(),
    id: String(item.id),
    provider: item.extra.provider,
    type: item.extra.type,
    raw: JSON.stringify(item),
  };
};

export const fromPersistanceToRelease = (
  row: DbReleasesTable,
): BluRayComReleaseWithExtra => {
  const parsed = JSON.parse(row.raw);

  return bluRayComReleaseWithExtraSchema.parse(parsed);
};

export const fromReleaseSourceToPersistance = (
  item: BluRayComCountry,
): DbReleaseSourcesTable => {
  return {
    id: item.code,
    provider: ReleaseSourceProvider.BLU_RAY_COM,
    raw: JSON.stringify(item),
  };
};

export const fromPersistanceToReleaseSurce = (
  row: DbReleaseSourcesTable,
): BluRayComCountry => {
  const parsed = JSON.parse(row.raw);

  return bluRayComCountrySchema.parse(parsed);
};
