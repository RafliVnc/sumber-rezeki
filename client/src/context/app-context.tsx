"use client";

import { UserDummy } from "@/dummy/user-dummy";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";

type AuthContextType = {
  user: User;
  setUser: (value: User) => void;
};

const AuthContext = createContext<AuthContextType>({
  user: UserDummy,
  setUser: (_value: User) => {},
});

const AppProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User>(UserDummy);
  const router = useRouter();

  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem("token");
      if (token) {
        try {
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
          localStorage.removeItem("token");
          setUser(UserDummy);
        }
      }
      Cookies.remove("token");
      router.push("/login");
    };

    fetchUser();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        setUser: (value: User) => {
          setUser(value);
        },
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
