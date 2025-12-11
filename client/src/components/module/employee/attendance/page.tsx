"use client";
import React from "react";
import { AttendanceTable } from "./table";
import { Card, CardContent } from "@/components/ui/card";

export default function AttendancePage() {
  const employees = [
    { id: "1", name: "Stapan" },
    { id: "2", name: "Stapan" },
    { id: "3", name: "Stapan" },
    { id: "4", name: "Stapan" },
    { id: "5", name: "Stapan" },
    { id: "6", name: "Stapan" },
    { id: "7", name: "Stapan" },
    { id: "8", name: "Stapan" },
    { id: "9", name: "Stapan" },
    { id: "10", name: "Stapan" },
    { id: "11", name: "Stapan" },
    { id: "12", name: "Stapan" },
    { id: "13", name: "Stapan" },
  ];

  const handleAttendanceChange = (
    employeeId: string,
    date: string,
    status: "hadir" | "tidak_hadir" | null
  ) => {
    console.log(`Employee ${employeeId} on ${date}: ${status}`);
  };

  const handleConfirm = (records: any) => {
    console.log("Attendance confirmed:", records);
    // Here you can save to database or API
  };
  return (
    <Card>
      <CardContent>
        <AttendanceTable
          currentDate={new Date()}
          employees={employees}
          onAttendanceChange={handleAttendanceChange}
          onConfirm={handleConfirm}
        />
      </CardContent>
    </Card>
  );
}
