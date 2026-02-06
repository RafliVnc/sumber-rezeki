"use client";

import { FactoryDummy } from "@/dummy/factory-dummy";
import { api } from "@/lib/api";
import { FactoryValidation } from "@/validation/factory-validation";
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
import Required from "@/components/ui/required";

const ApiFormFactory = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof FactoryValidation.MAIN>;
  isEdit?: boolean;
  id?: number;
}) => {
  const res = await api<{ data: Factory }>({
    url: isEdit ? `factories/${id}` : "factories",
    method: isEdit ? "PUT" : "POST",
    body: values,
  });

  return res;
};

export default function FormFactory({
  factory,
  handleClose,
}: {
  factory?: Factory;
  handleClose: () => void;
}) {
  const queryClient = useQueryClient();
  const isEdit = !!factory;

  const dummyFormFactory: z.infer<typeof FactoryValidation.MAIN> = {
    name: factory?.name || FactoryDummy.name,
    phone: factory?.phone || FactoryDummy.phone,
    dueDate: factory?.dueDate || FactoryDummy.dueDate,
    description: factory?.description || FactoryDummy.description,
  };

  const form = useForm<z.infer<typeof FactoryValidation.MAIN>>({
    resolver: zodResolver(FactoryValidation.MAIN),
    defaultValues: dummyFormFactory,
  });

  const submitMutation = useMutation({
    mutationFn: ApiFormFactory,
    onSuccess: () => {
      toast.success(`Kilang berhasil ${isEdit ? "diubah" : " ditambahkan"}`);
      queryClient.invalidateQueries({ queryKey: ["factories"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan kilang");
    },
  });

  const onSubmit = (values: z.infer<typeof FactoryValidation.MAIN>) => {
    submitMutation.mutate({ values, isEdit, id: factory?.id });
  };
  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col gap-4 h-full px-4 pb-4"
      >
        <div className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required children={"Nama"} />
                </FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Nama" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="dueDate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required children={"Jatuh Tempo (Hari)"} />
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Ketik hari"
                    value={field.value || ""}
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
          <FormField
            control={form.control}
            name="phone"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required children={"No. HP"} />
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Masukkan Nomor Hp"
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
