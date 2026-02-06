"use client";

import { api } from "@/lib/api";
import { SalesValidation } from "@/validation/sales-validation";
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
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { SalesDummy } from "@/dummy/sales-dummy";
import { MultiSelect } from "@/components/ui/multiple-select";

// TODO: Refactor delete Post request
const apiFormSales = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof SalesValidation.MAIN>;
  isEdit?: boolean;
  id?: number;
}): Promise<{ data: Sales }> => {
  const res = await api<{ data: Sales }>({
    url: isEdit ? `sales/${id}` : "sales",
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

export default function FormSales({
  handleClose,
  sales,
}: {
  handleClose: () => void;
  sales?: Sales;
}) {
  const queryClient = useQueryClient();
  const isEdit = !!sales;

  const dummyFormSales: z.infer<typeof SalesValidation.MAIN> = {
    name: sales?.Employee.name || SalesDummy.Employee.name,
    phone: sales?.phone || SalesDummy.phone,
    routeIds: sales?.Routes?.map((r) => r.id) || [],
  };

  const { data: routes, isLoading } = useQuery({
    queryKey: ["form-sales-routes"],
    queryFn: apiGetRoute,
  });

  const routeList =
    routes?.data?.map((r) => ({
      value: r.id.toString(),
      label: r.name,
    })) || [];

  const form = useForm<z.infer<typeof SalesValidation.MAIN>>({
    resolver: zodResolver(SalesValidation.MAIN),
    defaultValues: dummyFormSales,
  });

  const submitMutation = useMutation({
    mutationFn: apiFormSales,
    onSuccess: () => {
      toast.success(`Sales berhasil ${isEdit ? "diubah" : "ditambahkan"}`);
      queryClient.invalidateQueries({ queryKey: ["sales"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan sales");
    },
  });

  const onSubmit = (values: z.infer<typeof SalesValidation.MAIN>) => {
    submitMutation.mutate({ values, isEdit, id: sales?.id });
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col justify-between h-full px-4 pb-4"
      >
        <div className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nama*</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Nama" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="phone"
            render={({ field }) => (
              <FormItem>
                <FormLabel>No. HP*</FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="08123456789"
                    {...field}
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
                <FormLabel>Rute</FormLabel>
                <FormControl>
                  <MultiSelect
                    options={routeList}
                    maxCount={6}
                    isLoading={isLoading}
                    defaultValue={field.value ? field.value.map(String) : []}
                    onValueChange={(val) => field.onChange(val.map(Number))}
                    placeholder="Pilih rute"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        {/* Button Group */}
        <div className="flex space-x-2">
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
