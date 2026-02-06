"use client";

import { UserDummy } from "@/dummy/user-dummy";
import { api } from "@/lib/api";
import { UserValidation } from "@/validation/user-validation";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { InputPassword } from "@/components/ui/form/password-form";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { UserRole } from "@/type/enum/user-role";
import Required from "@/components/ui/required";

const apiFormUser = async ({
  values,
  isEdit,
  id,
}: {
  values: z.infer<typeof UserValidation.CREATE | typeof UserValidation.UPDATE>;
  isEdit?: boolean;
  id?: string;
}): Promise<{ data: User }> => {
  const res = await api<{ data: User }>({
    url: isEdit ? `users/${id}` : "users",
    method: isEdit ? "PUT" : "POST",
    body: values,
  });

  return res;
};

export default function FormUser({
  handleClose,
  user,
}: {
  handleClose: () => void;
  user?: User;
}) {
  const queryClient = useQueryClient();
  const isEdit = !!user;

  const dummyFormUser: z.infer<typeof UserValidation.BASE> = {
    name: user?.name || UserDummy.name,
    username: user?.username || UserDummy.username,
    phone: user?.phone || UserDummy.phone,
    role: user?.role || UserDummy.role,
  };

  const form = useForm<
    z.infer<typeof UserValidation.CREATE | typeof UserValidation.UPDATE>
  >({
    resolver: zodResolver(
      isEdit ? UserValidation.UPDATE : UserValidation.CREATE,
    ),
    defaultValues: isEdit ? dummyFormUser : { ...dummyFormUser, password: "" },
  });

  const submitMutation = useMutation({
    mutationFn: apiFormUser,
    onSuccess: () => {
      toast.success("Pengguna berhasil ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["users"] });
      form.reset();
      handleClose();
    },
    onError: (error: Error) => {
      toast.error(error?.message || "Gagal menambahkan pengguna");
    },
  });

  const onSubmit = (
    values: z.infer<
      typeof UserValidation.CREATE | typeof UserValidation.UPDATE
    >,
  ) => {
    submitMutation.mutate({ values, isEdit, id: user?.id });
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
                    <SelectItem value={UserRole.OWNER}>Owner</SelectItem>
                    <SelectItem value={UserRole.TREASURER}>
                      Bendahara
                    </SelectItem>
                    <SelectItem value={UserRole.WAREHOUSE_HEAD}>
                      Kepala Gudang
                    </SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required>Nama</Required>
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
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Required>Username</Required>
                </FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Username" {...field} />
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
                  <Required>No. Hp</Required>
                </FormLabel>
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
          {!user && (
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    <Required>Kata Sandi</Required>
                  </FormLabel>
                  <FormControl>
                    <InputPassword
                      placeholder="Kata Sandi"
                      {...field}
                      value={field.value ?? ""}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
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
