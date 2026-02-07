"use client";

import { VehicleDummy } from "@/dummy/vehicle-dummy";
import { api } from "@/lib/api";
import { VehicleValidation } from "@/validation/vehicle-validation";
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
import { Input } from "@/components/ui/input";
import Required from "@/components/ui/required";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ConvertToPlate, ConvertVehicleType } from "@/lib/utils";
import { VehicleType } from "@/type/enum/vehicle-type";

const ApiFormVehicle = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof VehicleValidation.MAIN>;
  isEdit?: boolean;
  id?: number;
}) => {
  const res = await api<{ data: Vehicle }>({
    url: isEdit ? `vehicles/${id}` : "vehicles",
    method: isEdit ? "PUT" : "POST",
    body: values,
  });

  return res;
};

export default function FormVehicle({
  vehicle,
  handleClose,
}: {
  vehicle?: Vehicle;
  handleClose: () => void;
}) {
  const queryClient = useQueryClient();
  const isEdit = !!vehicle;

  const dummyFormVehicle: z.infer<typeof VehicleValidation.MAIN> = {
    plate: ConvertToPlate(vehicle?.plate || VehicleDummy.plate),
    type: vehicle?.type || VehicleDummy.type,
  };

  const form = useForm<z.infer<typeof VehicleValidation.MAIN>>({
    resolver: zodResolver(VehicleValidation.MAIN),
    defaultValues: dummyFormVehicle,
  });

  const submitMutation = useMutation({
    mutationFn: ApiFormVehicle,
    onSuccess: () => {
      toast.success(`Kendaraan berhasil ${isEdit ? "diubah" : " ditambahkan"}`);
      queryClient.invalidateQueries({ queryKey: ["vehicles-tronton"] });
      queryClient.invalidateQueries({ queryKey: ["vehicles-non-tronton"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan kendaraan");
    },
  });

  const onSubmit = (values: z.infer<typeof VehicleValidation.MAIN>) => {
    submitMutation.mutate({ values, isEdit, id: vehicle?.id });
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
            name="plate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required children={"Plat Nomor"} />
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Masukkan Plat Nomor"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="type"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required children={"Jenis Kendaraan"} />
                </FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl className="w-full">
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(VehicleType).map((value) => (
                      <SelectItem key={value} value={value}>
                        {ConvertVehicleType(value)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
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
