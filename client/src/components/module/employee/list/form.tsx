"use client";

import { api } from "@/lib/api";
import { EmployeeDummy } from "@/dummy/employee-dummy";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import React from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { Input } from "@/components/ui/input";
import { EmployeeValidation } from "@/validation/employee-validation";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { EmployeeRole } from "@/type/enum/employee-role";
import { ConvertEmployeeRole, formatIDR } from "@/lib/utils";
import Required from "@/components/ui/required";
import { MultiSelect } from "@/components/ui/multiple-select";

const ApiFormEmployee = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof EmployeeValidation.MAIN>;
  isEdit?: boolean;
  id?: number;
}) => {
  const res = await api<{ data: Employee }>({
    url: isEdit ? `employees/${id}` : "employees",
    method: isEdit ? "PUT" : "POST",
    body: values,
  });

  return res;
};

const apiGetRoute = async (): Promise<{ data: Route[] }> => {
  const res = await api<{ data: Route[] }>({
    url: "routes",
    method: "GET",
  });

  return res;
};

const apiGetEmployeeSales = async (): Promise<{ data: Employee[] }> => {
  const res = await api<{ data: Employee[] }>({
    url: "employees",
    method: "GET",
    params: { roles: [EmployeeRole.SALES] },
  });

  return res;
};

const apiGetEmployeeById = async (id: number): Promise<{ data: Employee }> => {
  const res = await api<{ data: Employee }>({
    url: `employees/${id}`,
    method: "GET",
  });

  return res;
};

export default function FormEmployee({
  id,
  handleClose,
}: {
  id?: number;
  handleClose: () => void;
}) {
  const queryClient = useQueryClient();

  const { data: employee, isLoading: isLoadingEmployee } = useQuery({
    queryKey: ["form-employee", id],
    queryFn: () => apiGetEmployeeById(id!),
    staleTime: 0,
    enabled: !!id,
  });

  const isEdit = !!employee;

  const defaultValues: z.infer<typeof EmployeeValidation.MAIN> = {
    name: EmployeeDummy.name,
    salary: EmployeeDummy.salary,
    role: EmployeeDummy.role,
  };

  const form = useForm<z.infer<typeof EmployeeValidation.MAIN>>({
    resolver: zodResolver(EmployeeValidation.MAIN),
    values:
      isEdit && employee?.data
        ? {
            name: employee.data.name,
            salary: employee.data.salary,
            role: employee.data.role,
            supervisorId: employee.data?.supervisorId,
            phone: employee.data.Sales?.phone,
            routeIds: employee.data.Sales?.Routes?.map((r) => r.id) || [],
          }
        : defaultValues,
  });

  const { data: routes, isLoading: isLoadingRoutes } = useQuery({
    queryKey: ["form-employee-sales-routes"],
    queryFn: apiGetRoute,
    enabled:
      form.watch("role") === EmployeeRole.SALES ||
      employee?.data.role == EmployeeRole.SALES,
    staleTime: 0,
  });

  const routeList =
    routes?.data?.map((r) => ({
      value: r.id.toString(),
      label: r.name,
    })) || [];

  const { data: sales, isLoading: isLoadingSales } = useQuery({
    queryKey: ["form-employee-sales", id],
    queryFn: apiGetEmployeeSales,
    enabled:
      form.watch("role") === EmployeeRole.HELPER ||
      employee?.data.role == EmployeeRole.HELPER ||
      form.watch("role") === EmployeeRole.DRIVER ||
      employee?.data.role == EmployeeRole.DRIVER,
    staleTime: 0,
  });

  const salesList =
    sales?.data?.map((s) => ({
      value: s.id.toString(),
      label: s.name,
    })) || [];

  const isLoadingData = isLoadingEmployee || isLoadingRoutes || isLoadingSales;

  const submitMutation = useMutation({
    mutationFn: ApiFormEmployee,
    onSuccess: () => {
      toast.success("Karyawan berhasil ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["employees"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan karyawan");
    },
  });

  const onSubmit = (values: z.infer<typeof EmployeeValidation.MAIN>) => {
    submitMutation.mutate({ values, isEdit, id: employee?.data.id });
  };

  return (
    <Form {...form}>
      {isLoadingData && (
        <div className="absolute inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="flex flex-col items-center gap-2">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-sm text-muted-foreground">Memuat data...</p>
          </div>
        </div>
      )}
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="space-y-4 h-full px-4 pb-4"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Required>Nama</Required>
              </FormLabel>
              <FormControl>
                <Input type="text" placeholder="Masukan Nama" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Required>Role</Required>
              </FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl className="w-full">
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {Object.values(EmployeeRole).map((value) => (
                    <SelectItem key={value} value={value}>
                      {ConvertEmployeeRole(value)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="salary"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Required>Gaji</Required>
              </FormLabel>
              <FormControl>
                <Input
                  type="text"
                  placeholder="gaji"
                  value={formatIDR(field.value || 0)}
                  onChange={(e) => {
                    const raw = e.target.value.replace(/\D/g, "");
                    const numericValue = Number(raw || 0);

                    field.onChange(numericValue);
                  }}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {(form.watch("role") === EmployeeRole.DRIVER ||
          form.watch("role") === EmployeeRole.HELPER) && (
          <FormField
            control={form.control}
            name="supervisorId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required>Sales</Required>
                </FormLabel>
                <FormControl>
                  <Select
                    onValueChange={(value) => field.onChange(parseInt(value))}
                    value={field.value?.toString()}
                    defaultValue={field.value?.toString()}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Pilih Sales" />
                    </SelectTrigger>
                    <SelectContent>
                      {salesList && salesList.length === 0 && (
                        <SelectItem value="no-sales" disabled>
                          Silahkan tambahkan sales terlebih dahulu
                        </SelectItem>
                      )}

                      {salesList &&
                        salesList.length > 0 &&
                        salesList.map((s) => (
                          <SelectItem key={s.value} value={s.value}>
                            {s.label}
                          </SelectItem>
                        ))}
                    </SelectContent>
                  </Select>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />
        )}

        {form.watch("role") === EmployeeRole.SALES && (
          <>
            <FormField
              control={form.control}
              name="phone"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    <Required>No. HP</Required>
                  </FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="08123456789"
                      {...field}
                      value={field.value || ""}
                      onChange={(e) => {
                        const value = e.target.value;
                        const numericValue = value.replace(/[^0-9]/g, "");
                        if (numericValue.length <= 13) {
                          field.onChange(numericValue);
                        }
                      }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="routeIds"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    <Required>Rute</Required>
                  </FormLabel>
                  <FormControl>
                    <MultiSelect
                      options={routeList}
                      maxCount={6}
                      isLoading={isLoadingRoutes}
                      defaultValue={field.value ? field.value.map(String) : []}
                      onValueChange={(val) => field.onChange(val.map(Number))}
                      placeholder="Pilih rute"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </>
        )}

        {/* Button Group */}
        <div className="flex space-x-2 mt-8">
          <Button
            className="flex-1"
            type="button"
            variant={"outline"}
            disabled={submitMutation.isPending || submitMutation.isSuccess}
            onClick={handleClose}
          >
            {submitMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              "Batal"
            )}
          </Button>
          <Button
            className="flex-1"
            type="submit"
            disabled={submitMutation.isPending || submitMutation.isSuccess}
          >
            {submitMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              "Simpan"
            )}
          </Button>
        </div>
      </form>
    </Form>
  );
}
