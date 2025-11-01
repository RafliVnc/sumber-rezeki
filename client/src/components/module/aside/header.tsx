"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { useApp } from "@/context/app-context";
import { ConvertUserRole } from "@/lib/utils";
import { Key, LogOut } from "lucide-react";
import React from "react";

export default function Header() {
  const app = useApp();
  const { user } = app;

  return (
    <header className="bg-white sticky top-0 flex h-16 shrink-0 items-center gap-2 border-b px-2 z-50">
      <div className="flex flex-1">
        <SidebarTrigger className="-ml-1 rounded-full" />
        {/* <Separator
          orientation="vertical"
          className="mr-2 data-[orientation=vertical]:h-6"
        /> */}
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            size="lg"
            variant="ghost"
            className="data-[state=open]:!bg-primary/10 data-[state=open]:text-text-black "
          >
            <div className="grid flex-1 text-left text-sm leading-tight">
              <span className="truncate font-medium">{user.name}</span>
              <span className="truncate text-xs text-primary">
                {ConvertUserRole(user.role)}
              </span>
            </div>
            <Avatar className="h-8 w-8 rounded-lg border-2 border-primary/30">
              <AvatarImage src={""} alt={user.name} />
              <AvatarFallback className="rounded-lg">
                {user.name.charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
          side={"bottom"}
          align="end"
          sideOffset={4}
        >
          <DropdownMenuLabel className="p-0 font-normal">
            <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage src={"/logo.svg"} alt={user.name} />
                <AvatarFallback className="rounded-lg">CN</AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.name}</span>
                <span className="truncate text-xs text-primary">
                  {ConvertUserRole(user.role)}
                </span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem>
              <Key />
              Perbarui Kata Sandi
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={app.logout}>
            <LogOut />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
  );
}
