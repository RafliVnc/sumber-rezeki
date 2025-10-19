import { NextRequest, NextResponse } from "next/server";
import {
  // superAdminRoutes,
  authRoutes,
  DEFAULT_LOGIN_REDIRECT,
  noSiteRoutes,
} from "./routes";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  const token = req.cookies.get("token");

  const isAuthRoutes = authRoutes.includes(pathname);
  // const isSuperAdminRoutes = superAdminRoutes.includes(pathname);
  const isNoSiteRoutes = noSiteRoutes.includes(pathname);

  if (isNoSiteRoutes) {
    return Response.redirect(new URL("/login", req.url));
  }

  if (!isAuthRoutes && !token) {
    return Response.redirect(new URL("/login", req.url));
  }

  if (isAuthRoutes) {
    if (token) {
      return Response.redirect(new URL(DEFAULT_LOGIN_REDIRECT, req.url));
    }
    return;
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!.*\\..*|_next).*)", "/", "/(api|trpc)(.*)"],
};
