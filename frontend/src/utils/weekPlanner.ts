import { TimeEntry } from "../types";

export type PlannerRow = {
  weekday: string;
  date: string;
  entryId: number | null;
  startTime: string;
  endTime: string;
  lunchDuration: number;
};

export type MonthTotals = {
  monthLabel: string;
  weeks: Array<{
    key: string;
    weekNumber: number;
    minutes: number;
  }>;
  monthTotal: number;
};

const WEEKDAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

export function toMondayUTC(base: Date): Date {
  const utcDate = new Date(Date.UTC(base.getUTCFullYear(), base.getUTCMonth(), base.getUTCDate()));
  const weekday = utcDate.getUTCDay() === 0 ? 7 : utcDate.getUTCDay();
  utcDate.setUTCDate(utcDate.getUTCDate() - (weekday - 1));
  return utcDate;
}

export function dateInputFromUTC(date: Date): string {
  const y = date.getUTCFullYear();
  const m = String(date.getUTCMonth() + 1).padStart(2, "0");
  const d = String(date.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

export function buildWeekRows(weekOffset: number, baseDate: Date = new Date()): PlannerRow[] {
  const monday = toMondayUTC(baseDate);
  monday.setUTCDate(monday.getUTCDate() + weekOffset * 7);
  return WEEKDAYS.map((weekday, index) => {
    const day = new Date(monday);
    day.setUTCDate(monday.getUTCDate() + index);
    return {
      weekday,
      date: dateInputFromUTC(day),
      entryId: null,
      startTime: "",
      endTime: "",
      lunchDuration: 30,
    };
  });
}

export function isoWeekInfo(date: Date): { week: number; year: number } {
  const target = new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()));
  const dayNr = (target.getUTCDay() + 6) % 7;
  target.setUTCDate(target.getUTCDate() - dayNr + 3);

  const firstThursday = new Date(Date.UTC(target.getUTCFullYear(), 0, 4));
  const firstDayNr = (firstThursday.getUTCDay() + 6) % 7;
  firstThursday.setUTCDate(firstThursday.getUTCDate() - firstDayNr + 3);

  const week = 1 + Math.round((target.getTime() - firstThursday.getTime()) / (7 * 24 * 60 * 60 * 1000));
  return { week, year: target.getUTCFullYear() };
}

export function minutesFromTimeInput(value: string): number | null {
  if (!/^\d{2}:\d{2}$/.test(value)) {
    return null;
  }
  const [hoursRaw, minsRaw] = value.split(":");
  const hours = Number(hoursRaw);
  const mins = Number(minsRaw);
  if (Number.isNaN(hours) || Number.isNaN(mins) || hours < 0 || hours > 23 || mins < 0 || mins > 59) {
    return null;
  }
  return hours * 60 + mins;
}

export function roundWorkedMinutes(minutes: number): number {
  if (minutes <= 0) {
    return 0;
  }
  const remainder = minutes % 30;
  if (remainder === 0) {
    return minutes;
  }
  if (remainder < 5) {
    return minutes - remainder;
  }
  return minutes + (30 - remainder);
}

export function formatHoursFromMinutes(minutes: number): string {
  const hours = minutes / 60;
  return `${hours.toFixed(1).replace(/\.0$/, "")}h`;
}

export function rowRoundedTotal(row: PlannerRow): number | null {
  const start = minutesFromTimeInput(row.startTime);
  const end = minutesFromTimeInput(row.endTime);
  if (start === null || end === null) {
    return null;
  }
  const rawWorked = Math.max(0, end - start - row.lunchDuration);
  return roundWorkedMinutes(rawWorked);
}

export function buildMonthlyTotals(allEntries: TimeEntry[], monthBase: Date = new Date()): MonthTotals {
  const monthStart = new Date(Date.UTC(monthBase.getUTCFullYear(), monthBase.getUTCMonth(), 1, 0, 0, 0, 0));
  const monthEnd = new Date(
    Date.UTC(monthBase.getUTCFullYear(), monthBase.getUTCMonth() + 1, 0, 23, 59, 59, 999)
  );

  const monthLabel = monthStart.toLocaleString(undefined, {
    month: "long",
    year: "numeric",
    timeZone: "UTC",
  });

  const weekOrder: string[] = [];
  const weekNumberByKey = new Map<string, number>();
  const weeklyMinutes = new Map<string, number>();
  const hasWeekdayInMonth = new Map<string, boolean>();

  const cursor = new Date(monthStart);
  while (cursor <= monthEnd) {
    const { week, year } = isoWeekInfo(cursor);
    const key = `${year}-${week}`;
    if (!weeklyMinutes.has(key)) {
      weekOrder.push(key);
      weeklyMinutes.set(key, 0);
      weekNumberByKey.set(key, week);
      hasWeekdayInMonth.set(key, false);
    }

    const day = cursor.getUTCDay();
    const isWeekday = day >= 1 && day <= 5;
    if (isWeekday) {
      hasWeekdayInMonth.set(key, true);
    }
    cursor.setUTCDate(cursor.getUTCDate() + 1);
  }

  let monthTotal = 0;

  allEntries.forEach((entry) => {
    const start = new Date(entry.start_time);
    if (Number.isNaN(start.getTime())) {
      return;
    }

    const roundedMinutes =
      typeof entry.rounded_worked_minutes === "number"
        ? entry.rounded_worked_minutes
        : roundWorkedMinutes(entry.worked_minutes || 0);

    if (start >= monthStart && start <= monthEnd) {
      monthTotal += roundedMinutes;
      const { week, year } = isoWeekInfo(start);
      const key = `${year}-${week}`;
      weeklyMinutes.set(key, (weeklyMinutes.get(key) || 0) + roundedMinutes);
    }
  });

  const weeks = weekOrder
    .filter((key) => hasWeekdayInMonth.get(key))
    .map((key) => ({
      key,
      weekNumber: weekNumberByKey.get(key) || 0,
      minutes: weeklyMinutes.get(key) || 0,
    }));

  return {
    monthLabel,
    weeks,
    monthTotal,
  };
}
