"use client";

import type React from "react";
import { createContext, useContext, useState } from "react";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  TriangleAlert,
  Info,
  CheckCircle,
  AlertTriangle,
  MessageSquare,
} from "lucide-react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

type DialogVariant = "destructive" | "info" | "warning" | "success" | "input";

interface ConfirmationDialogOptions {
  title: string;
  description: string;
  confirmText?: string;
  cancelText?: string;
  variant?: DialogVariant;
  inputConfig?: {
    label: string;
    placeholder?: string;
    required?: boolean;
    type?: "input" | "textarea";
    minLength?: number;
  };
}

interface ConfirmationDialogResult {
  confirmed: boolean;
  inputValue?: string;
}

interface ConfirmationDialogContextType {
  showConfirmation: (
    options: ConfirmationDialogOptions
  ) => Promise<ConfirmationDialogResult>;
}

const ConfirmationDialogContext = createContext<
  ConfirmationDialogContextType | undefined
>(undefined);

const variantConfig = {
  destructive: {
    icon: TriangleAlert,
    iconBg: "bg-red-100",
    iconColor: "text-red-600",
    buttonBg: "bg-red-600 hover:bg-red-700",
    buttonVariant: "destructive" as const,
  },
  info: {
    icon: Info,
    iconBg: "bg-blue-100",
    iconColor: "text-blue-600",
    buttonBg: "bg-blue-500 hover:bg-blue-600",
    buttonVariant: "default" as const,
  },
  warning: {
    icon: AlertTriangle,
    iconBg: "bg-amber-100",
    iconColor: "text-amber-700",
    buttonBg: "bg-yellow-500 hover:bg-yellow-600 text-yellow-50",
    buttonVariant: "default" as const,
  },
  success: {
    icon: CheckCircle,
    iconBg: "bg-green-100",
    iconColor: "text-green-600",
    buttonBg: "bg-green-600 hover:bg-green-700",
    buttonVariant: "default" as const,
  },
  input: {
    icon: MessageSquare,
    iconBg: "bg-gray-100",
    iconColor: "text-gray-600",
    buttonBg: "bg-blue-600 hover:bg-blue-700",
    buttonVariant: "default" as const,
  },
};

export function ConfirmationDialogProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [open, setOpen] = useState(false);
  const [options, setOptions] = useState<ConfirmationDialogOptions | null>(
    null
  );
  const [resolveRef, setResolveRef] = useState<
    ((value: ConfirmationDialogResult) => void) | null
  >(null);
  const [inputValue, setInputValue] = useState("");
  const [inputError, setInputError] = useState("");

  const showConfirmation = (
    options: ConfirmationDialogOptions
  ): Promise<ConfirmationDialogResult> => {
    setOptions(options);
    setInputValue("");
    setInputError("");
    setOpen(true);
    return new Promise<ConfirmationDialogResult>((resolve) => {
      setResolveRef(() => resolve);
    });
  };

  const validateInput = (): boolean => {
    if (!options?.inputConfig) return true;

    const { required, minLength } = options.inputConfig;

    if (required && !inputValue.trim()) {
      setInputError("Kolom ini wajib diisi");
      return false;
    }

    if (minLength && inputValue.trim().length < minLength) {
      setInputError(`Minimum ${minLength} characters required`);
      return false;
    }

    setInputError("");
    return true;
  };

  const handleConfirm = () => {
    if (options?.variant === "input") {
      const isValid = validateInput();
      if (!isValid) return;
    }

    if (resolveRef) {
      resolveRef({
        confirmed: true,
        inputValue: options?.variant === "input" ? inputValue : undefined,
      });
    }

    setOpen(false);
  };

  const handleCancel = () => {
    if (resolveRef) {
      resolveRef({ confirmed: false });
    }
    setOpen(false);
  };

  return (
    <ConfirmationDialogContext.Provider value={{ showConfirmation }}>
      {children}
      {options && (
        <AlertDialog open={open} onOpenChange={setOpen}>
          <AlertDialogContent className="max-w-md rounded-lg p-6 shadow-lg">
            <div className="mb-6 flex flex-col items-center justify-center text-center">
              <div
                className={`mb-4 flex h-16 w-16 items-center justify-center rounded-full ${
                  options.variant
                    ? variantConfig[options.variant].iconBg
                    : variantConfig["destructive"].iconBg
                } ${
                  options.variant
                    ? variantConfig[options.variant].iconColor
                    : variantConfig["destructive"].iconColor
                }`}
              >
                {(() => {
                  const IconComponent = options.variant
                    ? variantConfig[options.variant].icon
                    : variantConfig["destructive"].icon;
                  return <IconComponent className="h-8 w-8" />;
                })()}
              </div>
              <AlertDialogHeader className="space-y-2 text-center">
                <AlertDialogTitle className="text-center text-xl font-bold">
                  {options.title}
                </AlertDialogTitle>
                <AlertDialogDescription className="text-center text-sm text-gray-500">
                  {options.description}
                </AlertDialogDescription>
              </AlertDialogHeader>
            </div>

            {options.variant === "input" && options.inputConfig && (
              <div className="mb-4 space-y-2">
                <Label htmlFor="dialog-input" className="text-sm font-medium">
                  {options.inputConfig.label}
                  {options.inputConfig.required && (
                    <span className="ml-0.5 text-red-500">*</span>
                  )}
                </Label>
                <Input
                  id="dialog-input"
                  placeholder={options.inputConfig.placeholder}
                  value={inputValue}
                  onChange={(e) => {
                    setInputValue(e.target.value);
                    if (inputError) setInputError("");
                  }}
                  className={inputError ? "border-red-500" : ""}
                />
                {inputError && (
                  <p className="text-sm text-red-500">{inputError}</p>
                )}
              </div>
            )}

            <AlertDialogFooter className="flex w-full flex-row gap-2">
              <AlertDialogCancel
                className="flex-1 border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
                onClick={handleCancel}
              >
                {options.cancelText || "Cancel"}
              </AlertDialogCancel>

              <Button
                variant={
                  options.variant
                    ? variantConfig[options.variant].buttonVariant
                    : variantConfig["destructive"].buttonVariant
                }
                onClick={handleConfirm}
                className={`flex-1 ${
                  options.variant
                    ? variantConfig[options.variant].buttonBg
                    : variantConfig["destructive"].buttonBg
                }`}
              >
                {options.confirmText || "Confirm"}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </ConfirmationDialogContext.Provider>
  );
}

export const useConfirmationDialog = () => {
  const context = useContext(ConfirmationDialogContext);
  if (!context) {
    throw new Error(
      "useConfirmationDialog must be used within a ConfirmationDialogProvider"
    );
  }
  return context;
};
