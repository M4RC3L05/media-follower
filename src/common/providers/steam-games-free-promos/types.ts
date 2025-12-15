import z from "@zod/zod";

export enum SteamGamesFreePromoTypes {
  FREE_TO_PLAY = "free-to-play",
  FREE_TO_KEEP = "free-to-keep",
}

export const steamGamesFreePromosInputSchema = z.object({ url: z.url() });

export type SteamGamesFreePromosInput = z.infer<
  typeof steamGamesFreePromosInputSchema
>;

export const steamGamesFreePromosOutputSchema = z.object({
  id: z.union([
    z.number(),
    z.string().pipe(z.coerce.number()).pipe(z.number()),
  ]),
  image: z.url(),
  link: z.url(),
  name: z.string().min(1),
  promoType: z.enum(SteamGamesFreePromoTypes),
  startDate: z.union([
    z.string().pipe(z.coerce.date()).pipe(z.date()),
    z.date(),
  ]),
  endDate: z.union([z.string().pipe(z.coerce.date()).pipe(z.date()), z.date()]),
});

export type SteamGamesFreePromosOutput = z.infer<
  typeof steamGamesFreePromosOutputSchema
>;

export type Input = SteamGamesFreePromosInput;
export type Output = SteamGamesFreePromosOutput;
