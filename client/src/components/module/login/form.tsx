"use client";

import React from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { AuthValidation } from "@/validation/auth-validation";
import { InputPassword } from "@/components/ui/form/password-form";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useApp } from "@/context/app-context";
import Cookies from "js-cookie";

const postAuth = async (values: z.infer<typeof AuthValidation.MAIN>) => {
  const res = await api<{ data: { User: User; token: string } }>({
    url: "login",
    method: "POST",
    body: values,
  });

  return res;
};

export default function form() {
  const router = useRouter();
  const app = useApp();

  const form = useForm<z.infer<typeof AuthValidation.MAIN>>({
    resolver: zodResolver(AuthValidation.MAIN),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const submitMutation = useMutation({
    mutationFn: postAuth,
    onSuccess: ({ data }) => {
      // set token
      Cookies.set("token", data.token, { expires: 1 });
      localStorage.setItem("token", data.token);

      app.setUser(data.User);
      // rederect
      router.push("/user");
      toast.success("Selamat datang kembali!");
    },
    onError: (error: unknown) => {
      toast.error(error instanceof Error ? error.message : "Operasi gagal");
    },
  });

  function handleSubmit(values: z.infer<typeof AuthValidation.MAIN>) {
    submitMutation.mutate(values);
  }

  function onSubmit(values: z.infer<typeof AuthValidation.MAIN>) {
    handleSubmit(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 ">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username/ No Handphone</FormLabel>
              <FormControl>
                <Input
                  type="email"
                  autoComplete="username"
                  placeholder="Username/ No Handphone"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Kata Sandi</FormLabel>
              <FormControl>
                <InputPassword
                  type="password"
                  placeholder="Kata Sandi"
                  autoComplete="current-password"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          className="w-full"
          type="submit"
          disabled={submitMutation.isPending || submitMutation.isSuccess}
        >
          {submitMutation.isPending ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            "Masuk"
          )}
        </Button>
      </form>
    </Form>
  );
}
