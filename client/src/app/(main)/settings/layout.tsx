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

  const currentTab = pathname.includes("/route") ? "route" : "user";

  const title = currentTab === "user" ? "Pengguna" : "Rute";

  const handleTabChange = (value: string) => {
    router.push(`/settings/${value}`);
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
          <TabsTrigger value="user">Pengguna</TabsTrigger>
          <TabsTrigger value="route">Rute</TabsTrigger>
        </TabsList>
        <div>{children}</div>
      </Tabs>
    </div>
  );
}
