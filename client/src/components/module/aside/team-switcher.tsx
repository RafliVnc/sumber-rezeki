"use client";

import * as React from "react";

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import Image from "next/image";

export function TeamSwitcher() {
  const { state } = useSidebar();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size="lg"
          className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground hover:!bg-primary/10 hover:!text-black"
        >
          <div className="text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg overflow-hidden">
            <Image
              className="dark:invert object-contain"
              src="/logo.svg"
              alt="Next.js logo"
              width={state === "collapsed" ? 23 : 24}
              height={state === "collapsed" ? 21 : 18}
              priority
            />
          </div>
          {state !== "collapsed" && (
            <>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <h5 className="font-medium text-md">SUMBER REZEKI</h5>
              </div>
            </>
          )}
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
