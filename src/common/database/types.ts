/* GENERATED FILE CONTENT DO NOT EDIT */

import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { ReleaseType } from "#src/common/database/enums/release-type.ts";

export type DbReleaseSourcesTable = {
  id: string;
  provider: ReleaseSourceProvider;
  raw: string;
};

export type DbReleasesTable = {
  id: string;
  provider: ReleaseSourceProvider;
  type: ReleaseType;
  releasedAt: string;
  raw: string;
};
