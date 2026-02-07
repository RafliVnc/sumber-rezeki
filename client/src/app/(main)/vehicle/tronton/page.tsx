import TrontonPage from "@/components/module/vehicle/tronton/page";
import { Metadata } from "next";
import React from "react";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Tronton",
  };
}

export default function trontonPage() {
  return <TrontonPage />;
}
