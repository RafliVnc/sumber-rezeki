import { UserRole } from "@/type/enum/user-role";

export const UserDummy: User = {
  id: "",
  name: "",
  username: "",
  phone: "",
  role: UserRole.OWNER,
  createdAt: new Date(),
};
