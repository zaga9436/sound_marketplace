import { AuthRedirectGuard } from "@/widgets/auth/auth-redirect-guard";

export default function AuthLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <main className="page-frame">
      <div className="container flex min-h-screen items-center justify-center py-16">
        <AuthRedirectGuard>{children}</AuthRedirectGuard>
      </div>
    </main>
  );
}
