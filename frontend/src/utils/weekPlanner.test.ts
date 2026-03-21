import { TimeEntry } from "../types";
import {
  buildMonthlyTotals,
  buildWeekRows,
  isoWeekInfo,
  roundWorkedMinutes,
  rowRoundedTotal,
} from "./weekPlanner";

function makeEntry(dateISO: string, roundedWorkedMinutes: number): TimeEntry {
  return {
    id: Number(`${Date.parse(dateISO)}`.slice(-6)),
    start_time: dateISO,
    end_time: dateISO,
    is_open: false,
    lunch_duration: 30,
    source: "manual",
    created_at: dateISO,
    updated_at: dateISO,
    worked_minutes: roundedWorkedMinutes,
    rounded_worked_minutes: roundedWorkedMinutes,
    rounded_worked_hours: roundedWorkedMinutes / 60,
  };
}

function makeOpenEntry(dateISO: string): TimeEntry {
  return {
    id: Number(`${Date.parse(dateISO)}`.slice(-6)) + 1,
    start_time: dateISO,
    end_time: dateISO,
    is_open: true,
    lunch_duration: 30,
    source: "manual",
    created_at: dateISO,
    updated_at: dateISO,
    worked_minutes: 0,
    rounded_worked_minutes: 0,
    rounded_worked_hours: 0,
  };
}

describe("weekPlanner utilities", () => {
  test("buildWeekRows returns Monday-Sunday for the selected week", () => {
    const baseDate = new Date("2026-03-18T12:00:00Z"); // Wednesday
    const rows = buildWeekRows(0, baseDate);

    expect(rows).toHaveLength(7);
    expect(rows[0].weekday).toBe("Monday");
    expect(rows[0].date).toBe("2026-03-16");
    expect(rows[6].weekday).toBe("Sunday");
    expect(rows[6].date).toBe("2026-03-22");

    const nextWeek = buildWeekRows(1, baseDate);
    expect(nextWeek[0].date).toBe("2026-03-23");
    expect(nextWeek[6].date).toBe("2026-03-29");
  });

  test("roundWorkedMinutes follows 5-minute threshold", () => {
    expect(roundWorkedMinutes(454)).toBe(450);
    expect(roundWorkedMinutes(455)).toBe(480);
    expect(roundWorkedMinutes(4)).toBe(0);
    expect(roundWorkedMinutes(5)).toBe(30);
  });

  test("rowRoundedTotal computes and rounds final worked minutes", () => {
    expect(
      rowRoundedTotal({
        weekday: "Monday",
        date: "2026-03-16",
        entryId: null,
        startTime: "09:00",
        endTime: "17:05",
        lunchDuration: 30,
      })
    ).toBe(480);
  });

  test("buildMonthlyTotals lists all month weeks with weekdays and excludes leisure-only week", () => {
    const entries = [
      makeEntry("2026-03-03T09:00:00Z", 480), // week 10
      makeEntry("2026-03-31T09:00:00Z", 240), // week 14
      makeOpenEntry("2026-03-18T09:00:00Z"), // open entry should not contribute
      makeEntry("2026-04-01T09:00:00Z", 300), // outside March
    ];

    const totals = buildMonthlyTotals(entries, new Date("2026-03-20T00:00:00Z"));

    expect(totals.weeks.map((week) => week.weekNumber)).toEqual([10, 11, 12, 13, 14]);
    expect(totals.monthTotal).toBe(720);

    const byWeek = new Map(totals.weeks.map((week) => [week.weekNumber, week.minutes]));
    expect(byWeek.get(10)).toBe(480);
    expect(byWeek.get(11)).toBe(0);
    expect(byWeek.get(12)).toBe(0);
    expect(byWeek.get(13)).toBe(0);
    expect(byWeek.get(14)).toBe(240);
  });

  test("isoWeekInfo returns expected ISO week/year", () => {
    const info = isoWeekInfo(new Date("2026-03-16T00:00:00Z"));
    expect(info).toEqual({ week: 12, year: 2026 });
  });
});
