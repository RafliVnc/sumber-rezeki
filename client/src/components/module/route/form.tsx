"use client";

import { RouteDummy } from "@/dummy/route-dummy";
import { api } from "@/lib/api";
import { RouteValidation } from "@/validation/route-validation";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
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
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";

const ApiFormRoute = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof RouteValidation.MAIN>;
  isEdit?: boolean;
  id?: number;
}) => {
  const res = await api<{ data: Route }>({
    url: isEdit ? `routes/${id}` : "routes",
    method: isEdit ? "PUT" : "POST",
    body: values,
  });

  return res;
};

export default function FormRoute({
  route,
  handleClose,
}: {
  route?: Route;
  handleClose: () => void;
}) {
  const queryClient = useQueryClient();
  const isEdit = !!route;

  const dummyFormRoute: z.infer<typeof RouteValidation.MAIN> = {
    name: route?.name || RouteDummy.name,
    description: route?.description || RouteDummy.description,
  };

  const form = useForm<z.infer<typeof RouteValidation.MAIN>>({
    resolver: zodResolver(RouteValidation.MAIN),
    defaultValues: dummyFormRoute,
  });

  const submitMutation = useMutation({
    mutationFn: ApiFormRoute,
    onSuccess: () => {
      toast.success("Sales berhasil ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["routes"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan sales");
    },
  });

  const onSubmit = (values: z.infer<typeof RouteValidation.MAIN>) => {
    submitMutation.mutate({ values, isEdit, id: route?.id });
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
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Deskripsi</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="Deskripsi"
                    className="min-h-40"
                    {...field}
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
