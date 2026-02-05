"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import {
  ChevronLeft,
  ChevronRight,
  Plus,
  X,
  Save,
  Wrench,
  FileDown,
} from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { formatDate, getDayName, getSundayOfWeek } from "@/lib/utils";
import { Skeleton } from "@/components/ui/skeleton";
import { AttendanceStatus } from "@/type/enum/attendance-status";
import { AttendanceUIData } from "./page";
import {
  AttendanceValidation,
  BatchAttendanceInput,
  UIFormInput,
} from "@/validation/attendance-validation";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

interface AttendanceTableProps {
  currentDate?: Date;
  employees: Array<{
    id: number;
    name: string;
  }>;
  initialRecords?: AttendanceUIData[];
  totalEmployees: number;
  onWeekChange?: (newDate: Date) => void;
  isLoading?: boolean;
  weekStartStr: string;
  weekEndStr: string;
}

const getNextStatus = (
  current: AttendanceStatus | null,
): AttendanceStatus | null => {
  if (current === null) {
    return AttendanceStatus.PRESENT;
  } else if (current === AttendanceStatus.PRESENT) {
    return AttendanceStatus.LEAVE;
  } else if (current === AttendanceStatus.LEAVE) {
    return AttendanceStatus.SICK;
  } else if (current === AttendanceStatus.SICK) {
    return AttendanceStatus.ABSENT;
  } else {
    return null;
  }
};

const ApiSaveAttendance = async (data: BatchAttendanceInput) => {
  return await api<{ message: string }, BatchAttendanceInput>({
    url: "attendance/batch",
    method: "POST",
    body: data,
  });
};

export function AttendanceTable({
  currentDate = new Date(),
  employees,
  initialRecords = [],
  totalEmployees,
  onWeekChange,
  isLoading,
  weekStartStr,
  weekEndStr,
}: AttendanceTableProps) {
  const [isEditMode, setIsEditMode] = useState(false);
  const [displayDate, setDisplayDate] = useState(new Date(currentDate));
  const queryClient = useQueryClient();

  const form = useForm<UIFormInput>({
    resolver: zodResolver(AttendanceValidation.UI_FORM),
    defaultValues: {
      dates: {},
    },
    mode: "onChange",
  });

  const saveMutation = useMutation({
    mutationFn: ApiSaveAttendance,
    onSuccess: () => {
      toast.success("Data absensi berhasil disimpan!");
      queryClient.invalidateQueries({
        queryKey: ["attendance", weekStartStr, weekEndStr],
      });
      setIsEditMode(false);
    },
    onError: (error: Error) => {
      console.error("Error saving attendance:", error);
      toast.error(
        error?.message || "Gagal menyimpan data absensi. Silakan coba lagi.",
      );
    },
  });

  useEffect(() => {
    setDisplayDate(new Date(currentDate));
  }, [currentDate]);

  const weekStart = getSundayOfWeek(displayDate);
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekEnd.getDate() + 6);

  const generateWeekDates = () => {
    const dates = [];
    const current = new Date(weekStart);
    while (current <= weekEnd) {
      const day = current.getDay();
      if (day !== 5) {
        dates.push(new Date(current));
      }
      current.setDate(current.getDate() + 1);
    }
    return dates;
  };

  const dates = generateWeekDates();

  useEffect(() => {
    const initialDates: UIFormInput["dates"] = {};

    dates.forEach((date) => {
      const dateStr = formatDate(date);
      const hasData = initialRecords.some(
        (r) =>
          r.attendance[dateStr] !== null && r.attendance[dateStr] !== undefined,
      );

      if (hasData) {
        initialDates[dateStr] = {
          isActive: true,
          employees: employees.map((emp) => {
            const record = initialRecords.find((r) => r.employeeId === emp.id);
            return {
              id: emp.id,
              name: emp.name,
              status: record?.attendance[dateStr] || null,
            };
          }),
        };
      }
    });

    form.reset({ dates: initialDates });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialRecords, employees, weekStart.getTime()]);

  const activateDateForEditing = (dateStr: string) => {
    const currentDates = form.getValues("dates");

    if (!currentDates[dateStr] || !currentDates[dateStr].isActive) {
      form.setValue(`dates.${dateStr}`, {
        isActive: true,
        employees: employees.map((emp) => ({
          id: emp.id,
          name: emp.name,
          status: AttendanceStatus.PRESENT,
        })),
      });
    }
  };

  const deactivateDate = (dateStr: string) => {
    const currentDates = form.getValues("dates");

    const hasExistingData = initialRecords.some(
      (r) =>
        r.attendance[dateStr] !== null && r.attendance[dateStr] !== undefined,
    );

    if (hasExistingData) {
      form.setValue(`dates.${dateStr}`, {
        isActive: true,
        employees: employees.map((emp) => ({
          id: emp.id,
          name: emp.name,
          status: null,
        })),
      });
    } else {
      const newDates = { ...currentDates };
      delete newDates[dateStr];
      form.setValue("dates", newDates);
    }
  };

  const toggleAttendance = (dateStr: string, employeeId: number) => {
    const currentDates = form.getValues("dates");
    if (!currentDates[dateStr]) return;

    const employeeIndex = currentDates[dateStr].employees.findIndex(
      (e) => e.id === employeeId,
    );

    if (employeeIndex === -1) return;

    const currentStatus = currentDates[dateStr].employees[employeeIndex].status;
    const newStatus = getNextStatus(currentStatus);

    form.setValue(
      `dates.${dateStr}.employees.${employeeIndex}.status`,
      newStatus,
      {
        shouldValidate: true,
        shouldDirty: true,
      },
    );
  };

  const handleFormSubmit = (data: UIFormInput) => {
    const attendances: BatchAttendanceInput["attendances"] = [];

    Object.entries(data.dates).forEach(([date, dateData]) => {
      if (!dateData.isActive) return;

      const allNull = dateData.employees.every((e) => e.status === null);

      if (allNull) {
        attendances.push({
          date,
          employees: [],
          action: "delete",
        });
      } else {
        const validEmployees = dateData.employees.filter(
          (e) => e.status !== null,
        );

        if (validEmployees.length !== totalEmployees) {
          toast.error(`Tanggal ${date}: Semua karyawan harus memiliki status`);
          throw new Error("Incomplete data");
        }

        const employees = validEmployees.map((e) => ({
          id: e.id,
          status: e.status!,
        }));

        attendances.push({
          date,
          employees,
          action: "update",
        });
      }
    });

    if (attendances.length === 0) {
      toast.error("Tidak ada data untuk disimpan");
      return;
    }

    try {
      const validatedData = AttendanceValidation.BATCH(totalEmployees).parse({
        attendances,
      });

      saveMutation.mutate(validatedData);
    } catch (error) {
      console.error("Validation error:", error);
      toast.error("Validasi gagal. Pastikan semua data lengkap.");
    }
  };

  const handleCancel = () => {
    const initialDates: UIFormInput["dates"] = {};

    dates.forEach((date) => {
      const dateStr = formatDate(date);
      const hasData = initialRecords.some(
        (r) =>
          r.attendance[dateStr] !== null && r.attendance[dateStr] !== undefined,
      );

      if (hasData) {
        initialDates[dateStr] = {
          isActive: true,
          employees: employees.map((emp) => {
            const record = initialRecords.find((r) => r.employeeId === emp.id);
            return {
              id: emp.id,
              name: emp.name,
              status: record?.attendance[dateStr] || null,
            };
          }),
        };
      }
    });

    form.reset({ dates: initialDates });
    setIsEditMode(false);
  };

  const getStatusBadge = (status: AttendanceStatus | null) => {
    if (status === AttendanceStatus.PRESENT) {
      return (
        <span className="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-sm font-medium text-green-800">
          Hadir
        </span>
      );
    }
    if (status === AttendanceStatus.LEAVE) {
      return (
        <span className="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-sm font-medium text-blue-800">
          Izin
        </span>
      );
    }
    if (status === AttendanceStatus.SICK) {
      return (
        <span className="inline-flex items-center rounded-full bg-yellow-100 px-2.5 py-0.5 text-sm font-medium text-yellow-800">
          Sakit
        </span>
      );
    }
    if (status === AttendanceStatus.ABSENT) {
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
    onWeekChange?.(newDate);
  };

  const goToNextWeek = () => {
    const newDate = new Date(displayDate);
    newDate.setDate(newDate.getDate() + 7);
    setDisplayDate(newDate);
    onWeekChange?.(newDate);
  };

  const getIsNextWeekDisabled = () => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const nextWeekStart = new Date(weekStart);
    nextWeekStart.setDate(nextWeekStart.getDate() + 7);
    nextWeekStart.setHours(0, 0, 0, 0);
    return nextWeekStart > today;
  };

  const isNextWeekDisabled = getIsNextWeekDisabled();

  const formDates = form.watch("dates");
  const activeDatesCount = Object.values(formDates).filter(
    (d) => d?.isActive,
  ).length;

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleFormSubmit)}
        className="flex flex-col max-h-[680px] overflow-hidden"
      >
        <div className="flex items-center justify-between mb-4 bg-white">
          <div className="flex items-center  gap-1">
            <Button
              type="button"
              variant="outline"
              size="icon-sm"
              onClick={goToPreviousWeek}
              disabled={isEditMode}
            >
              <ChevronLeft className="w-4 h-4" />
            </Button>
            <span className="md:text-md font-bold text-center min-w-48">
              {weekStart.toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}{" "}
              -{" "}
              {weekEnd.toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
            </span>
            <Button
              type="button"
              variant="outline"
              size="icon-sm"
              onClick={goToNextWeek}
              disabled={isNextWeekDisabled || isEditMode}
            >
              <ChevronRight className="w-4 h-4" />
            </Button>
          </div>
          <div className="flex gap-2">
            {!isEditMode ? (
              <>
                <Button
                  type="button"
                  size="sm"
                  onClick={(e) => {
                    e.preventDefault();
                    setIsEditMode(true);
                  }}
                  disabled={isLoading || employees.length === 0}
                >
                  <Wrench className="size-4" />
                  Pengaturan Absensi
                </Button>
                <Button
                  type="button"
                  variant="default"
                  size="sm"
                  disabled={isLoading || employees.length === 0}
                >
                  <FileDown className="size-4" />
                  Export PDF
                </Button>
              </>
            ) : (
              <>
                <Button
                  type="submit"
                  size="sm"
                  disabled={saveMutation.isPending || activeDatesCount === 0}
                  className="bg-green-500 hover:bg-green-600"
                >
                  <Save className="w-4 h-4 mr-2" />
                  {saveMutation.isPending ? "Menyimpan..." : `Simpan`}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleCancel}
                  disabled={saveMutation.isPending}
                >
                  Batal
                </Button>
              </>
            )}
          </div>
        </div>

        <div className="relative flex-1 overflow-auto rounded-md border">
          <Table noWrapper>
            <TableHeader className="sticky top-0 bg-white z-30">
              <TableRow>
                <TableHead className="sticky left-0 bg-white min-w-20 z-30">
                  Nama
                </TableHead>
                {dates.map((date) => {
                  const dateStr = formatDate(date);
                  const isActive = formDates[dateStr]?.isActive || false;

                  return (
                    <TableHead
                      key={dateStr}
                      className={`text-center min-w-36 whitespace-nowrap py-3 ${
                        isActive && isEditMode
                          ? formDates[dateStr]?.employees.every(
                              (e) => e.status === null,
                            )
                            ? "bg-red-50"
                            : "bg-blue-50"
                          : ""
                      }`}
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
                            {!isActive ? (
                              <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                onClick={() => activateDateForEditing(dateStr)}
                                title="Tambahkan hari ini ke form"
                              >
                                <Plus className="w-4 h-4 text-primary" />
                              </Button>
                            ) : (
                              <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                onClick={() => deactivateDate(dateStr)}
                                title="Hapus hari ini dari form"
                              >
                                <X className="w-4 h-4 text-red-600" />
                              </Button>
                            )}
                          </div>
                        )}
                      </div>
                    </TableHead>
                  );
                })}
              </TableRow>
            </TableHeader>

            <TableBody>
              {isLoading ? (
                [...Array(10)].map((_, rowIndex) => (
                  <TableRow key={rowIndex}>
                    <TableCell className="sticky left-0 bg-white">
                      <Skeleton className="h-5 w-full" />
                    </TableCell>
                    {[...Array(dates.length)].map((_, colIndex) => (
                      <TableCell key={colIndex}>
                        <Skeleton className="h-5 w-20 mx-auto" />
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : employees.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={dates.length + 1}
                    className="h-64 text-center"
                  >
                    <div className="flex flex-col items-center justify-center text-gray-500">
                      <p className="text-lg font-medium">
                        Tidak ada data karyawan
                      </p>
                      <p className="text-sm">
                        Silakan tambahkan karyawan terlebih dahulu
                      </p>
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                employees.map((employee) => (
                  <TableRow key={employee.id}>
                    <TableCell className="font-semibold sticky left-0 bg-white z-10">
                      {employee.name}
                    </TableCell>
                    {dates.map((date) => {
                      const dateStr = formatDate(date);
                      const dateData = formDates[dateStr];
                      const isActive = dateData?.isActive || false;
                      const isMarkedForDeletion =
                        isActive &&
                        dateData?.employees.every((e) => e.status === null);

                      const initialRecord = initialRecords.find(
                        (r) => r.employeeId === employee.id,
                      );
                      const displayStatus: AttendanceStatus | null =
                        initialRecord?.attendance?.[dateStr] ?? null;

                      const employeeIndex =
                        dateData?.employees.findIndex(
                          (e) => e.id === employee.id,
                        ) ?? -1;

                      return (
                        <TableCell
                          key={dateStr}
                          className={`text-center py-2 ${
                            isActive && isEditMode
                              ? isMarkedForDeletion
                                ? "bg-red-50"
                                : "bg-blue-50"
                              : ""
                          }`}
                        >
                          {isActive && isEditMode && employeeIndex !== -1 ? (
                            <FormField
                              control={form.control}
                              name={`dates.${dateStr}.employees.${employeeIndex}.status`}
                              render={({ field }) => (
                                <FormItem>
                                  <FormControl>
                                    <button
                                      type="button"
                                      onClick={() =>
                                        toggleAttendance(dateStr, employee.id)
                                      }
                                      className="w-full h-full py-1 rounded-full hover:bg-gray-50 transition-colors cursor-pointer inline-flex items-center justify-center"
                                    >
                                      {getStatusBadge(field.value)}
                                    </button>
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                          ) : (
                            <div className="py-1 inline-flex items-center justify-center">
                              {getStatusBadge(displayStatus)}
                            </div>
                          )}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </form>
    </Form>
  );
}
