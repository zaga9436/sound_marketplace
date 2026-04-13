import { CardDetail } from "@/widgets/cards/card-detail";

export default async function CardDetailsPage({
  params
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <main className="page-frame">
      <div className="container space-y-8 py-10">
        <CardDetail id={id} />
      </div>
    </main>
  );
}
