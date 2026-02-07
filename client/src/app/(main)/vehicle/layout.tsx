"use client";

import React from "react";
import { usePathname, useRouter } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DynamicBreadcrumb } from "@/components/ui/breadcrumb/dynamic-breadcrumb";

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();

  const currentTab = pathname.includes("/tronton") ? "tronton" : "non-tronton";

  const title = currentTab === "tronton" ? "Tronton" : "Kendaraan";

  const handleTabChange = (value: string) => {
    router.push(`/vehicle/${value}`);
  };

  return (
    <div className="flex w-full flex-col gap-2 p-4">
      <div>
        <h1 className="text-2xl font-bold mb-1">{title}</h1>
        <DynamicBreadcrumb />
      </div>
      <Tabs
        value={currentTab}
        onValueChange={handleTabChange}
        className="w-full"
      >
        <TabsList className="w-[400px]">
          <TabsTrigger value="tronton">Tronton</TabsTrigger>
          <TabsTrigger value="non-tronton">Kendaraan</TabsTrigger>
        </TabsList>
        <div>{children}</div>
      </Tabs>
    </div>
  );
}
