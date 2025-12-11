import AttendancePage from "@/components/module/employee/attendance/page";
import { Metadata } from "next";
import React from "react";
export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Absensi",
  };
}

export default function attendancePage() {
  return <AttendancePage />;
}
