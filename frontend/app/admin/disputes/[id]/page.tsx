import { AdminDisputeDetailPage } from "@/widgets/admin/admin-dispute-detail-page";

export default async function AdminDisputeDetailRoute({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <AdminDisputeDetailPage id={id} />;
}
