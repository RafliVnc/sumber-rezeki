"use client";

import { DynamicBreadcrumb } from "@/components/ui/breadcrumb/dynamic-breadcrumb";
import { Card, CardContent } from "@/components/ui/card";
import React from "react";
import TableUser from "./table/table";

export default function UserPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Pengguna</h1>
      <DynamicBreadcrumb />
      <Card className="mt-4">
        <CardContent>
          <TableUser />
        </CardContent>
      </Card>
    </div>
  );
}
