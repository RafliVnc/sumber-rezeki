"use client";

import { useApp } from "@/context/app-context";
import React from "react";

export default function UserPage() {
  const app = useApp();
  return <div>{app.user.name}</div>;
}
