"use client";

import * as React from "react";
import { CircleUser, Truck, User, Waypoints } from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { TeamSwitcher } from "./team-switcher";
import { NavMain } from "./nav-main";
import { NavUser } from "./nav-user";

// This is sample data.
const data = {
  navMain: [
    {
      title: "Karyawan",
      url: "/employee/attendance",
      activePath: "/employee",
      icon: CircleUser,
    },
    {
      title: "Sales",
      url: "/sales",
      activePath: "/sales",
      icon: User,
    },
    {
      title: "Kilang",
      url: "/factories",
      activePath: "/factories",
      icon: Waypoints,
    },
    {
      title: "Kendaraan",
      url: "/Vehicles",
      activePath: "/Vehicles",
      icon: Truck,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="border-sidebar-border h-16 border-b flex justify-center items-center">
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
