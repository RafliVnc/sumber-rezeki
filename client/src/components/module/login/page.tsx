import React from "react";
import Image from "next/image";
import FormLogin from "./form";

export default function LoginPage() {
  return (
    <div className="w-full h-screen grid grid-cols-11 gap-4 p-4">
      <div className="size-full border rounded-lg col-span-6 flex justify-center items-center ">
        <div className="w-1/2 md:space-y-6">
          <div className="flex justify-center items-center w-full flex-col">
            <Image
              className="dark:invert object-contain"
              src="/logo.svg"
              alt="Next.js logo"
              width={70}
              height={18}
              priority
            />
            <h1 className="text-2xl font-bold">SUMBER REZEKI</h1>
          </div>
          <div className="flex flex-col gap-4">
            <article>
              <h1 className="text-2xl font-bold">Masuk</h1>
              <p>Selamat Datang Kembali 👋🏻</p>
            </article>
            <FormLogin />
          </div>
        </div>
      </div>
      <div className="size-full bg-black rounded-lg col-span-5"></div>
    </div>
  );
}
