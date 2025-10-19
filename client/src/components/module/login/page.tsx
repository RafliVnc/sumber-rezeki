import React from "react";
import Form from "./form";

export default function LoginPage() {
  return (
    <div className="w-full h-screen grid grid-cols-10 gap-4 p-2">
      <div className="size-full border rounded-lg col-span-6 flex justify-center items-center ">
        <div className="w-1/2 flex flex-col gap-4">
          <article>
            <h1 className="text-2xl font-bold">Masuk</h1>
            <p>-</p>
          </article>
          <Form />
        </div>
      </div>
      <div className="size-full bg-black rounded-lg col-span-4"></div>
    </div>
  );
}
