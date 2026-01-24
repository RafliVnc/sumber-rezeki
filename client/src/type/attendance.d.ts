enum AttendanceStatus {
  PRESENT = "PRESENT",
  LEAVE = "LEAVE",
  SICK = "SICK",
  ABSENT = "ABSENT",
}

interface Attendance {
  id: number;
  status: AttendanceStatus;
  date: string;
  periodId: number;
}
