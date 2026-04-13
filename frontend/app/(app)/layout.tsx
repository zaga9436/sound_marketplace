import { AppShell } from "@/widgets/app-shell/app-shell";

export default function ProtectedAppLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <AppShell>{children}</AppShell>;
}
