import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Pengguna",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <main className="p-4">{children}</main>;
}
