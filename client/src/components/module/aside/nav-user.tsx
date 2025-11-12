"use client";

import { LogOut, Settings } from "lucide-react";

import {
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useApp } from "@/context/app-context";
import { usePathname } from "next/navigation";

export function NavUser() {
  const app = useApp();
  const pathname = usePathname();
  const isActive = pathname.startsWith("/settings");

  return (
    <SidebarMenu>
      <SidebarMenuItem className="space-y-1">
        <SidebarGroupLabel>Pengaturan</SidebarGroupLabel>
        <SidebarMenuButton asChild tooltip="Pengaturan" isActive={isActive}>
          <a href={"/settings/user"}>
            <Settings />
            <span>Pengaturan</span>
          </a>
        </SidebarMenuButton>

        <SidebarMenuButton onClick={app.logout}>
          <LogOut className="text-red-500" />
          <span className="text-red-500">Keluar</span>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
