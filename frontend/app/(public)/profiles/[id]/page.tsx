import { PublicProfile } from "@/widgets/profile/public-profile";

export default async function PublicProfilePage({
  params
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <main className="page-frame">
      <div className="container space-y-8 py-10">
        <PublicProfile id={id} />
      </div>
    </main>
  );
}
