export const MONTHS = [
  "Januari",
  "Februari",
  "Maret",
  "April",
  "Mei",
  "Juni",
  "Juli",
  "Agustus",
  "September",
  "Oktober",
  "November",
  "Desember",
];

const startYear = new Date().getFullYear() - 2;
const totalYears = 5;

export const YEARS = Array.from(
  { length: totalYears },
  (_, i) => startYear + i,
);
