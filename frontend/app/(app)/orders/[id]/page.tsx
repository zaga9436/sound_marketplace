import { OrderDetailPage } from "@/widgets/orders/order-detail-page";

export default async function OrderDetailRoutePage({
  params
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return <OrderDetailPage id={id} />;
}
