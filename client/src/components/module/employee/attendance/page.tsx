"use client";
import React, { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { AttendanceTable } from "./table";
import { Card, CardContent } from "@/components/ui/card";
import { api } from "@/lib/api";
import { formatDate, getSundayOfWeek } from "@/lib/utils";
import { AttendanceStatus } from "@/type/enum/attendance-status";

interface EmployeeAttendance {
  id: number;
  name: string;
  Attendaces?: Array<{
    date: string;
    status: AttendanceStatus;
  }>;
}

// Format data untuk UI (per employee)
export interface AttendanceUIData {
  employeeId: number;
  employeeName: string;
  attendance: {
    [date: string]: AttendanceStatus | null;
  };
}

export default function AttendancePage() {
  const [currentWeekStart, setCurrentWeekStart] = useState<Date>(new Date());

  const weekStart = getSundayOfWeek(currentWeekStart);
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekEnd.getDate() + 6);

  const startDateStr = formatDate(weekStart);
  const endDateStr = formatDate(weekEnd);

  const { data: attendanceData, isLoading } = useQuery({
    queryKey: ["attendance", startDateStr, endDateStr],
    queryFn: async () => {
      const result = await api<{ data: EmployeeAttendance[] }>({
        url: "attendance",
        method: "GET",
        params: {
          startDate: startDateStr,
          endDate: endDateStr,
        },
      });
      return result.data;
    },
    staleTime: 1000 * 60 * 5,
  });

  // Transform data untuk UI (struktur per-employee untuk easy editing)
  const uiData: AttendanceUIData[] =
    attendanceData?.map((emp) => {
      const attendance: { [date: string]: AttendanceStatus | null } = {};

      if (emp.Attendaces && emp.Attendaces.length > 0) {
        emp.Attendaces.forEach((att) => {
          const dateStr = att.date.split("T")[0];
          attendance[dateStr] = att.status;
        });
      }

      return {
        employeeId: emp.id,
        employeeName: emp.name,
        attendance,
      };
    }) || [];

  const employees =
    attendanceData?.map((emp) => ({
      id: emp.id,
      name: emp.name,
    })) || [];

  const handleWeekChange = (newDate: Date) => {
    setCurrentWeekStart(newDate);
  };

  return (
    <Card>
      <CardContent>
        <AttendanceTable
          currentDate={currentWeekStart}
          employees={employees}
          initialRecords={uiData}
          totalEmployees={attendanceData?.length || 0}
          onWeekChange={handleWeekChange}
          isLoading={isLoading}
          weekStartStr={startDateStr}
          weekEndStr={endDateStr}
        />
      </CardContent>
    </Card>
  );
}
