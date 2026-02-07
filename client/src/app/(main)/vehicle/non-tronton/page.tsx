import NonTrontonPage from "@/components/module/vehicle/non-tronton/page";
import { Metadata } from "next";
import React from "react";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Kendaraan",
  };
}

export default function nonTrontonPage() {
  return <NonTrontonPage />;
}
