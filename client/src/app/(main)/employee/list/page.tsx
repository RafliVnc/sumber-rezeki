import EmployeePage from "@/components/module/employee/list/page";
import { Metadata } from "next";
import React from "react";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Karyawan",
  };
}

export default function ListPage() {
  return <EmployeePage />;
}
