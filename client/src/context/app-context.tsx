"use client";

import { UserDummy } from "@/dummy/user-dummy";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";
import { toast } from "sonner";

type AuthContextType = {
  user: User;
  setUser: (value: User) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType>({
  user: UserDummy,
  setUser: (_value: User) => {},
  logout: () => {},
});

const AppProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User>(UserDummy);
  const router = useRouter();

  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem("token");
      if (token) {
        try {
          Cookies.set("token", token, { expires: 1 });
          const res = await api<{ data: User }>({
            url: "current",
            method: "GET",
          });
          const data = await res.data;
          setUser(data);
          return;
        } catch (err) {
          console.error(err);
          // token invalid
          if (err instanceof Error && err.message === "Unauthorized") {
            localStorage.removeItem("token");
          }
          setUser(UserDummy);
        }
      }
      Cookies.remove("token");
      router.push("/login");
    };

    fetchUser();
  }, []);

  const logout = () => {
    localStorage.removeItem("token");
    Cookies.remove("token");
    setUser(UserDummy);
    toast.success("Anda telah keluar dari akun ini");
    router.push("/login");
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        setUser: (value: User) => {
          setUser(value);
        },
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

const useApp = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useApp must be used within a AppProvider");
  }
  return context;
};

export { AppProvider, useApp };
