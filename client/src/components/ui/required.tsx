import React from "react";

interface RequiredProps {
  children: React.ReactNode;
}

export default function Required({ children }: RequiredProps) {
  return (
    <div className="inline-flex items-center gap-1">
      {children}
      <span className="text-red-500">*</span>
    </div>
  );
}
