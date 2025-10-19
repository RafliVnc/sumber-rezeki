"use client";

import { useApp } from "@/context/app-context";
import React from "react";

export default function UserPage() {
  const app = useApp();
  console.log(app.user);
  return <div>{app.user.name}</div>;
}
