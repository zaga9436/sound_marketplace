import { CardForm } from "@/features/card/card-form";

export default async function EditCardPage({
  params
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return <CardForm mode="edit" cardId={id} />;
}
