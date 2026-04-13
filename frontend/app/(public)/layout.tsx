import { PublicHeader } from "@/widgets/public/public-header";

export default function PublicLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="page-frame">
      <PublicHeader />
      {children}
    </div>
  );
}
