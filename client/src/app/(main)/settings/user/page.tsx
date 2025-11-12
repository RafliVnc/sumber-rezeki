import UserPage from "@/components/module/user/page";
import { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  return {
    title: "Pengguna",
  };
}

export default function UserSettingsPage() {
  return <UserPage />;
}
