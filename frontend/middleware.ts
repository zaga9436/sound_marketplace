import { NextRequest, NextResponse } from "next/server";

const protectedPrefixes = ["/dashboard", "/orders", "/profile", "/chats", "/notifications", "/settings"];
const adminPrefixes = ["/admin"];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get("sm_token")?.value;
  const role = request.cookies.get("sm_role")?.value;

  const isCardManagementRoute = pathname === "/cards/new" || /^\/cards\/[^/]+\/edit$/.test(pathname);
  const isProtectedRoute = protectedPrefixes.some((prefix) => pathname.startsWith(prefix)) || isCardManagementRoute;
  const isAdminRoute = adminPrefixes.some((prefix) => pathname.startsWith(prefix));

  if ((isProtectedRoute || isAdminRoute) && !token) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("next", pathname);
    return NextResponse.redirect(url);
  }

  if (isAdminRoute && role !== "admin") {
    const url = request.nextUrl.clone();
    url.pathname = "/dashboard";
    return NextResponse.redirect(url);
  }

  if ((pathname === "/login" || pathname === "/register") && token) {
    const url = request.nextUrl.clone();
    url.pathname = role === "admin" ? "/admin" : "/dashboard";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/dashboard/:path*",
    "/orders/:path*",
    "/profile/:path*",
    "/chats/:path*",
    "/notifications/:path*",
    "/settings/:path*",
    "/admin/:path*",
    "/cards/:path*",
    "/login",
    "/register"
  ]
};
