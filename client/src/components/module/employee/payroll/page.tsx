"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { MONTHS, YEARS } from "@/constants/cons";
import { ChevronLeft, ChevronRight, HandCoins, Pencil } from "lucide-react";
import React, { useState } from "react";

export default function PayrolePage() {
  const now = new Date();
  const [selectedMonth, setSelectedMonth] = useState(
    MONTHS[now.getMonth()].toLowerCase(),
  );
  const [selectedYear, setSelectedYear] = useState(
    now.getFullYear().toString(),
  );
  const [searchQuery, setSearchQuery] = useState("");

  const monthIndex = MONTHS.findIndex((m) => m.toLowerCase() === selectedMonth);

  const previousMonth = () => {
    if (monthIndex === 0) {
      setSelectedYear((parseInt(selectedYear) - 1).toString());
      setSelectedMonth(MONTHS[11].toLowerCase());
    } else {
      setSelectedMonth(MONTHS[monthIndex - 1].toLowerCase());
    }
  };

  const nextMonth = () => {
    if (monthIndex === 11) {
      setSelectedYear((parseInt(selectedYear) + 1).toString());
      setSelectedMonth(MONTHS[0].toLowerCase());
    } else {
      setSelectedMonth(MONTHS[monthIndex + 1].toLowerCase());
    }
  };

  return (
    <div>
      <Card>
        <CardContent>
          {/* Toolbar */}
          <div className="flex items-center justify-between gap-2 mb-4">
            <div className="flex items-center gap-2 flex-1">
              {/* Year filter */}
              <div className="flex items-center gap-1">
                {/* Previous Month Button */}
                <Button
                  variant="outline"
                  size="icon-sm"
                  onClick={previousMonth}
                  className="p-2 h-auto"
                >
                  <ChevronLeft className="w-4 h-4 " />
                </Button>

                {/* Month & Year Display with Dropdowns */}
                <div className="flex items-center gap-2">
                  <Select
                    value={selectedMonth}
                    onValueChange={setSelectedMonth}
                  >
                    <SelectTrigger className="text-md font-bold">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {MONTHS.map((month) => (
                        <SelectItem
                          key={month}
                          value={month.toLowerCase()}
                          className="capitalize"
                        >
                          {month}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>

                  <Select value={selectedYear} onValueChange={setSelectedYear}>
                    <SelectTrigger className="text-md font-bold">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {YEARS.map((year) => (
                        <SelectItem key={year} value={year.toString()}>
                          {year}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {/* Next Month Button */}
                <Button
                  variant="outline"
                  size="icon-sm"
                  onClick={nextMonth}
                  className="p-2 h-auto"
                >
                  <ChevronRight className="w-4 h-4" />
                </Button>
              </div>
            </div>

            <Button size="sm" className="h-8">
              <HandCoins className="size-4" />
              Lihat Bonus
            </Button>

            <Button size="sm" className="h-8">
              <Pencil className="size-4" />
              Ubah Bonus
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
