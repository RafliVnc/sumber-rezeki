import RoutePage from "@/components/module/route/page";
import { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Rute",
  };
}

export default function RouteSettingsPage() {
  return <RoutePage />;
}
