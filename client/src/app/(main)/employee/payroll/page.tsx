import { Metadata } from "next";
import React from "react";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Gaji",
  };
}

export default function PayrollPage() {
  return <div>PayrollPage</div>;
}
