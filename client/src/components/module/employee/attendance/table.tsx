"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight, Plus, X } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface AttendanceRecord {
  employeeId: string;
  employeeName: string;
  attendance: {
    [date: string]: "hadir" | "tidak_hadir" | null;
  };
}

interface AttendanceTableProps {
  currentDate?: Date;
  employees: Array<{
    id: string;
    name: string;
  }>;
  onAttendanceChange?: (
    employeeId: string,
    date: string,
    status: "hadir" | "tidak_hadir" | null
  ) => void;
  onConfirm?: (records: AttendanceRecord[]) => void;
}

export function AttendanceTable({
  currentDate = new Date(),
  employees,
  onAttendanceChange,
  onConfirm,
}: AttendanceTableProps) {
  const [records, setRecords] = useState<AttendanceRecord[]>(
    employees.map((emp) => ({
      employeeId: emp.id,
      employeeName: emp.name,
      attendance: {},
    }))
  );
  const [isEditMode, setIsEditMode] = useState(false);
  const [displayDate, setDisplayDate] = useState(new Date(currentDate));

  const getSundayOfWeek = (date: Date): Date => {
    const d = new Date(date);
    const day = d.getDay();
    const diff = d.getDate() - day;
    return new Date(d.setDate(diff));
  };

  const weekStart = getSundayOfWeek(displayDate);
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekEnd.getDate() + 6); // Saturday

  const getIsNextWeekDisabled = () => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const currentWeekEnd = new Date(weekEnd);
    currentWeekEnd.setHours(0, 0, 0, 0);

    // If today is still within or before the current week, disable next
    return today <= currentWeekEnd;
  };

  const generateWeekDates = () => {
    const dates = [];
    const current = new Date(weekStart);

    while (current <= weekEnd) {
      const day = current.getDay();
      // Include only: Sunday (0), Monday (1), Tuesday (2), Wednesday (3), Thursday (4), Saturday (6)
      if (day !== 5) {
        dates.push(new Date(current));
      }
      current.setDate(current.getDate() + 1);
    }
    return dates;
  };

  const dates = generateWeekDates();

  const getDayName = (date: Date): string => {
    const days = [
      "Minggu",
      "Senin",
      "Selasa",
      "Rabu",
      "Kamis",
      "Jumat",
      "Sabtu",
    ];
    return days[date.getDay()];
  };

  const getDateString = (date: Date): string => {
    return date.toISOString().split("T")[0];
  };

  const toggleAttendance = (employeeId: string, date: string) => {
    setRecords((prev) =>
      prev.map((record) => {
        if (record.employeeId === employeeId) {
          const currentStatus = record.attendance[date];
          const newStatus: "hadir" | "tidak_hadir" | null =
            currentStatus === null
              ? "hadir"
              : currentStatus === "hadir"
              ? "tidak_hadir"
              : null;

          onAttendanceChange?.(employeeId, date, newStatus);

          return {
            ...record,
            attendance: {
              ...record.attendance,
              [date]: newStatus,
            },
          };
        }
        return record;
      })
    );
  };

  const fillDateWithHadir = (date: string) => {
    setRecords((prev) =>
      prev.map((record) => ({
        ...record,
        attendance: {
          ...record.attendance,
          [date]: "hadir",
        },
      }))
    );
  };

  const clearDateData = (date: string) => {
    setRecords((prev) =>
      prev.map((record) => ({
        ...record,
        attendance: {
          ...record.attendance,
          [date]: null,
        },
      }))
    );
  };

  const handleConfirm = () => {
    onConfirm?.(records);
    setIsEditMode(false);
  };

  const getStatusBadge = (status: "hadir" | "tidak_hadir" | null) => {
    if (status === "hadir") {
      return (
        <span className="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-sm font-medium text-green-800">
          Hadir
        </span>
      );
    }
    if (status === "tidak_hadir") {
      return (
        <span className="inline-flex items-center rounded-full bg-red-100 px-2.5 py-0.5 text-sm font-medium text-red-800">
          Tidak Hadir
        </span>
      );
    }
    return <span className="text-gray-400">-</span>;
  };

  const goToPreviousWeek = () => {
    const newDate = new Date(displayDate);
    newDate.setDate(newDate.getDate() - 7);
    setDisplayDate(newDate);
  };

  const goToNextWeek = () => {
    const newDate = new Date(displayDate);
    newDate.setDate(newDate.getDate() + 7);
    setDisplayDate(newDate);
  };

  const isNextWeekDisabled = getIsNextWeekDisabled();

  const hasDateData = (date: string): boolean => {
    return records.some(
      (record) =>
        record.attendance[date] !== null &&
        record.attendance[date] !== undefined
    );
  };

  return (
    <div className="flex flex-col max-h-[680px] overflow-hidden">
      {/* Header with navigation */}
      <div className="flex items-center justify-between p-4 bg-white">
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={goToPreviousWeek}>
            <ChevronLeft className="w-4 h-4" />
          </Button>
          <span className="text-sm font-medium min-w-48">
            {weekStart.toLocaleDateString("id-ID", {
              day: "numeric",
              month: "short",
            })}{" "}
            -{" "}
            {weekEnd.toLocaleDateString("id-ID", {
              day: "numeric",
              month: "short",
              year: "numeric",
            })}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={goToNextWeek}
            disabled={isNextWeekDisabled}
            className={
              isNextWeekDisabled ? "opacity-50 cursor-not-allowed" : ""
            }
          >
            <ChevronRight className="w-4 h-4" />
          </Button>
        </div>
        <div className="flex gap-2">
          {!isEditMode ? (
            <Button
              onClick={() => setIsEditMode(true)}
              className="bg-blue-500 hover:bg-blue-600"
            >
              Pengaturan Absensi
            </Button>
          ) : (
            <Button
              onClick={handleConfirm}
              className="bg-green-500 hover:bg-green-600"
            >
              Konfirmasi
            </Button>
          )}
          <Button variant="outline" size="sm">
            Export PDF
          </Button>
        </div>
      </div>

      {/* Table container with scroll */}
      <div className="relative flex-1 overflow-auto rounded-md border">
        <Table noWrapper>
          <TableHeader className="sticky top-0 bg-white z-10">
            <TableRow>
              <TableHead className="sticky left-0 bg-white min-w-20">
                Nama
              </TableHead>
              {dates.map((date) => (
                <TableHead
                  key={getDateString(date)}
                  className="text-center min-w-36 whitespace-nowrap py-3"
                >
                  <div className="flex w-full justify-center items-center gap-2">
                    <div>
                      <div>{getDayName(date)}</div>
                      <div className="text-xs font-normal text-gray-600">
                        {date.toLocaleDateString("id-ID", {
                          month: "long",
                          day: "numeric",
                          year: "numeric",
                        })}
                      </div>
                    </div>
                    {isEditMode && (
                      <div>
                        {!hasDateData(getDateString(date)) ? (
                          <Button
                            variant={"ghost"}
                            size={"sm"}
                            onClick={() =>
                              fillDateWithHadir(getDateString(date))
                            }
                            title="Isi semua hari ini dengan Hadir"
                          >
                            <Plus className="w-4 h-4 text-primary" />
                          </Button>
                        ) : (
                          <Button
                            variant={"ghost"}
                            size={"sm"}
                            onClick={() => clearDateData(getDateString(date))}
                            title="Hapus semua data pada hari ini"
                          >
                            <X className="w-4 h-4 text-red-600" />
                          </Button>
                        )}
                      </div>
                    )}
                  </div>
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>

          <TableBody>
            {records.map((record) => (
              <TableRow key={record.employeeId}>
                <TableCell className="font-medium sticky left-0 bg-white">
                  {record.employeeName}
                </TableCell>
                {dates.map((date) => {
                  const dateStr = getDateString(date);
                  const status = record.attendance[dateStr];
                  return (
                    <TableCell key={dateStr} className="text-center py-2">
                      {isEditMode ? (
                        <button
                          onClick={() =>
                            toggleAttendance(record.employeeId, dateStr)
                          }
                          className="w-full h-full py-2 rounded-full hover:bg-gray-50 transition-colors cursor-pointer inline-flex items-center justify-center"
                        >
                          {getStatusBadge(status)}
                        </button>
                      ) : (
                        <div className="py-2 inline-flex items-center justify-center">
                          {getStatusBadge(status)}
                        </div>
                      )}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
