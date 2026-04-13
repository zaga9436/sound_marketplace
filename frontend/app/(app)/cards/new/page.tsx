import { CardForm } from "@/features/card/card-form";
import { CardType } from "@/shared/types/api";

export default async function CreateCardPage({
  searchParams
}: {
  searchParams?: Promise<{ type?: string }>;
}) {
  const resolvedSearchParams = searchParams ? await searchParams : undefined;
  const type = resolvedSearchParams?.type;
  const initialCardType: CardType | undefined = type === "offer" || type === "request" ? type : undefined;

  return <CardForm mode="create" initialCardType={initialCardType} />;
}
