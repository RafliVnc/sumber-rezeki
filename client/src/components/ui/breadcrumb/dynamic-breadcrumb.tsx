"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useMemo, Fragment } from "react";

import {
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const breadcrumbTranslations: Record<string, string> = {
  user: "Pengguna",
  route: "Rute",
  settings: "Pengaturan",
  attendance: "Absensi",
  payroll: "Gaji",
  list: "Daftar",
  employee: "Karyawan",
  factory: "Kilang",
  vehicle: "Kendaraan",
  "non-tronton": "Kendaraan",
};

export const translateBreadcrumb = (segment: string): string => {
  const lowercased = segment.toLowerCase();
  return breadcrumbTranslations[lowercased] || formatLabel(segment);
};

const formatLabel = (segment: string): string => {
  return segment
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
};

interface BreadcrumbProps {
  maxItems?: number;
  homeLabel?: string;
  homePath?: string;
}

export function DynamicBreadcrumb({
  maxItems = 3,
  homeLabel = "Dashboard",
  homePath = "/dashboard",
}: BreadcrumbProps) {
  const pathname = usePathname();

  const breadcrumbs = useMemo(() => {
    const segments = pathname.split("/").filter(Boolean);

    const items = segments.map((segment, index) => {
      const path = `/${segments.slice(0, index + 1).join("/")}`;
      const label = translateBreadcrumb(segment);
      const isLast = index === segments.length - 1;

      return {
        label,
        path,
        isLast,
      };
    });

    return items;
  }, [pathname]);

  const shouldCollapse = breadcrumbs.length > maxItems;
  const visibleItems = shouldCollapse
    ? [...breadcrumbs.slice(0, 1), ...breadcrumbs.slice(-(maxItems - 1))]
    : breadcrumbs;

  const collapsedItems = shouldCollapse
    ? breadcrumbs.slice(1, -(maxItems - 1))
    : [];

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink asChild>
            <Link href={homePath}>{homeLabel}</Link>
          </BreadcrumbLink>
        </BreadcrumbItem>

        {breadcrumbs.length > 0 && <BreadcrumbSeparator />}

        {shouldCollapse && collapsedItems.length > 0 && (
          <>
            <BreadcrumbItem>
              <DropdownMenu>
                <DropdownMenuTrigger className="flex items-center gap-1">
                  <BreadcrumbEllipsis className="size-4" />
                  <span className="sr-only">Toggle menu</span>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start">
                  {collapsedItems.map((item, index) => (
                    <DropdownMenuItem key={index} asChild>
                      <Link href={item.path}>{item.label}</Link>
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
          </>
        )}

        {visibleItems.map((item) => (
          <Fragment key={item.path}>
            <BreadcrumbItem>
              {item.isLast ? (
                <BreadcrumbPage>{item.label}</BreadcrumbPage>
              ) : (
                <BreadcrumbLink asChild>
                  <Link
                    href={
                      item.path === "/settings" ? "/settings/user" : item.path
                    }
                  >
                    {item.label}
                  </Link>
                </BreadcrumbLink>
              )}
            </BreadcrumbItem>
            {!item.isLast && <BreadcrumbSeparator />}
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
