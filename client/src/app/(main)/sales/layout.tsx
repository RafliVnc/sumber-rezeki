import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Sales",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <section className="p-4">{children}</section>;
}
