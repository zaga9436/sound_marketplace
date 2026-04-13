import { AdminCardDetailPage } from "@/widgets/admin/admin-card-detail-page";

export default async function AdminCardDetailRoute({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <AdminCardDetailPage id={id} />;
}
