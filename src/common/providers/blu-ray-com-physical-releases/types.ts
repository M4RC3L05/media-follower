import z from "@zod/zod";

export const bluRayComPhysicalReleaseOutputSchema = z.object({
  id: z.number(),
  casing: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  artworkurl: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  title_sort: z.string().trim(),
  title: z.string().trim(),
  edition: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  extended: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  title_keywords: z.string().trim(),
  studio: z.string().trim().transform((val) => val === "" ? undefined : val)
    .optional(),
  year: z.union([
    z.number(),
    z.string().trim().pipe(z.coerce.number()).pipe(z.number()),
  ]),
  yearend: z.union([
    z.number(),
    z.string().trim().pipe(z.coerce.number()).pipe(z.number()),
  ]),
  releasedate: z.union([
    z.date(),
    z.string().trim().pipe(z.coerce.date()).pipe(z.date()),
  ]),
  popularity: z.number(),
  width: z.number(),
  height: z.number(),
});

export type BluRayComPhysicalReleaseOutput = z.infer<
  typeof bluRayComPhysicalReleaseOutputSchema
>;

export enum EBLuRayComPhysicalReleaseType {
  BLURAY = "bluray",
  DVD = "dvd",
}

export const bluRayComPhysicalReleaseOutputWithExtraSchema =
  bluRayComPhysicalReleaseOutputSchema.extend({
    extra: z.object({
      artworkUrl: z.url(),
      type: z.enum(EBLuRayComPhysicalReleaseType),
      link: z.url(),
    }),
  });

export type BluRayComPhysicalReleaseOutputWithExtra = z.infer<
  typeof bluRayComPhysicalReleaseOutputWithExtraSchema
>;

export const bluRayComPhysicalReleaseInputSchema = z.object({
  code: z.string(),
  name: z.string(),
});

export type BluRayComPhysicalReleaseInput = z.infer<
  typeof bluRayComPhysicalReleaseInputSchema
>;

export type Input = BluRayComPhysicalReleaseInput;
export type Output = BluRayComPhysicalReleaseOutputWithExtra;
