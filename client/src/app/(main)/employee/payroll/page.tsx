import PayrolePage from "@/components/module/employee/payroll/page";
import { Metadata } from "next";
import React from "react";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Gaji",
  };
}

export default function payrollPage() {
  return <PayrolePage />;
}
