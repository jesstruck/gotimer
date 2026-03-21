import React, { useEffect, useMemo, useState } from "react";
import { TimeEntry, TimeEntryPayload } from "../types";
import {
  formatDateUTC,
  formatDuration,
  isoToUTCDateInput,
  isoToUTCTimeInput,
  nowUTCTimeInput,
} from "../utils/datetime";
import {
  buildMonthlyTotals,
  buildWeekRows,
  formatHoursFromMinutes,
  isoWeekInfo,
  PlannerRow,
  rowRoundedTotal,
} from "../utils/weekPlanner";

type SaveOptions = {
  auto: boolean;
};

async function parseJSON(response: Response): Promise<any> {
  if (!response.ok) {
    let message = `Request failed: ${response.status}`;
    try {
      const payload: { error?: string } = await response.json();
      if (payload.error) {
        message = payload.error;
      }
    } catch {
      // noop
    }
    throw new Error(message);
  }
  if (response.status === 204) {
    return null;
  }
  return response.json();
}

function isValidTimeInput(value: string): boolean {
  return /^\d{2}:\d{2}$/.test(value);
}

const WeekPlanner: React.FC = () => {
  const [weekOffset, setWeekOffset] = useState<number>(0);
  const weekTemplate = useMemo(() => buildWeekRows(weekOffset), [weekOffset]);
  const [rows, setRows] = useState<PlannerRow[]>(weekTemplate);
  const [allEntries, setAllEntries] = useState<TimeEntry[]>([]);
  const [busy, setBusy] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [message, setMessage] = useState<string>("");

  const weekTitle = useMemo(() => {
    const start = rows[0]?.date;
    const end = rows[6]?.date;
    if (!start || !end) {
      return "Week Planner";
    }
    return `Week Planner (${formatDateUTC(`${start}T00:00:00Z`)} - ${formatDateUTC(
      `${end}T00:00:00Z`
    )})`;
  }, [rows]);

  const weekNumberLabel = useMemo(() => {
    const monday = new Date(`${weekTemplate[0].date}T00:00:00Z`);
    const info = isoWeekInfo(monday);
    return ` Week ${info.week}`;
  }, [weekTemplate]);

  const totals = useMemo(() => {
    return buildMonthlyTotals(allEntries);
  }, [allEntries]);

  const loadWeekEntries = async (): Promise<void> => {
    setBusy(true);
    setError("");
    setMessage("");
    try {
      const response = await fetch("/api/time-entries");
      const allEntries: TimeEntry[] = await parseJSON(response);
      setAllEntries(allEntries);

      const weekDates = new Set(weekTemplate.map((row) => row.date));
      const latestByDate = new Map<string, TimeEntry>();

      allEntries.forEach((entry) => {
        const dayKey = isoToUTCDateInput(entry.start_time);
        if (!weekDates.has(dayKey)) {
          return;
        }
        const current = latestByDate.get(dayKey);
        if (!current) {
          latestByDate.set(dayKey, entry);
          return;
        }
        const currentUpdated = new Date(current.updated_at).getTime();
        const candidateUpdated = new Date(entry.updated_at).getTime();
        if (candidateUpdated > currentUpdated) {
          latestByDate.set(dayKey, entry);
        }
      });

      setRows(() =>
        weekTemplate.map((row) => {
          const entry = latestByDate.get(row.date);
          if (!entry) {
            return { ...row, entryId: null, startTime: "", endTime: "", lunchDuration: 30 };
          }
          return {
            ...row,
            entryId: entry.id,
            startTime: isoToUTCTimeInput(entry.start_time),
            endTime: entry.is_open ? "" : isoToUTCTimeInput(entry.end_time),
            lunchDuration: entry.lunch_duration,
          };
        })
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load week entries");
    } finally {
      setBusy(false);
    }
  };

  useEffect(() => {
    loadWeekEntries();
  }, [weekTemplate]);

  const updateRow = (date: string, patch: Partial<PlannerRow>): void => {
    setRows((prevRows) => prevRows.map((row) => (row.date === date ? { ...row, ...patch } : row)));
  };

  const handleStart = (date: string): void => {
    setMessage("");
    setError("");
    const nowTime = nowUTCTimeInput();
    const existingRow = rows.find((row) => row.date === date);
    if (!existingRow) {
      return;
    }
    const updatedRow = { ...existingRow, startTime: nowTime, endTime: "" };
    updateRow(date, { startTime: nowTime, endTime: "" });
    void saveRow(updatedRow, { auto: true });
  };

  const handleStop = (date: string): void => {
    setMessage("");
    setError("");
    const nowTime = nowUTCTimeInput();
    const existingRow = rows.find((row) => row.date === date);
    if (existingRow) {
      const updatedRow = { ...existingRow, endTime: nowTime };
      updateRow(date, { endTime: nowTime });
      void saveRow(updatedRow, { auto: true });
    }
  };

  const saveRow = async (row: PlannerRow, options?: SaveOptions): Promise<void> => {
    setBusy(true);
    setError("");
    setMessage("");
    try {
      if (!row.startTime) {
        throw new Error("Start time is required to save a day.");
      }
      if (!isValidTimeInput(row.startTime)) {
        throw new Error("Start time must use HH:mm.");
      }
      if (row.endTime && !isValidTimeInput(row.endTime)) {
        throw new Error("End time must use HH:mm.");
      }

      const payload: TimeEntryPayload = {
        date: row.date,
        start_time: row.startTime,
        lunch_duration: row.lunchDuration,
      };
      if (row.endTime) {
        payload.end_time = row.endTime;
      }

      if (row.entryId) {
        const response = await fetch(`/api/time-entries/${row.entryId}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
        await parseJSON(response);
      } else {
        const response = await fetch("/api/time-entries", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
        await parseJSON(response);
      }

      await loadWeekEntries();
      setMessage(options?.auto ? `Auto-saved ${row.weekday}.` : `Saved ${row.weekday}.`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save day.");
    } finally {
      setBusy(false);
    }
  };

  const maybeAutoSave = (row: PlannerRow): void => {
    if (busy) {
      return;
    }
    if (!isValidTimeInput(row.startTime)) {
      return;
    }
    if (row.endTime && !isValidTimeInput(row.endTime)) {
      return;
    }
    void saveRow(row, { auto: true });
  };

  const deleteRow = async (row: PlannerRow): Promise<void> => {
    setBusy(true);
    setError("");
    setMessage("");
    try {
      if (row.entryId) {
        const response = await fetch(`/api/time-entries/${row.entryId}`, {
          method: "DELETE",
        });
        await parseJSON(response);
      }

      updateRow(row.date, { entryId: null, startTime: "", endTime: "", lunchDuration: 30 });
      await loadWeekEntries();
      setMessage(`Deleted ${row.weekday}.`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete day.");
    } finally {
      setBusy(false);
    }
  };

  return (
    <main style={{ maxWidth: "1100px", margin: "0 auto", padding: "16px" }}>
      <h1>{weekTitle}</h1>
      <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
        <button type="button" onClick={() => setWeekOffset((prev) => prev - 1)} disabled={busy}>
          {"< Previous Week"}
        </button>
        <button type="button" onClick={() => setWeekOffset((prev) => prev + 1)} disabled={busy}>
          {"Next Week >"}
        </button>
        <button type="button" onClick={() => setWeekOffset(0)} disabled={busy || weekOffset === 0}>
          Current Week
        </button>
        <small>{weekNumberLabel}</small>
      </div>
      <p>Only Monday-Sunday is shown. Start/end is auto-saved.</p>
      {error && <p style={{ color: "#b00020" }}>{error}</p>}
      {message && <p style={{ color: "#0b6f36" }}>{message}</p>}

      <table style={{ width: "100%", borderCollapse: "collapse", marginTop: "12px" }}>
        <thead>
          <tr style={{ borderBottom: "2px solid #ddd" }}>
            <th style={{ textAlign: "left", padding: "8px" }}>Day</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Date</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Start (UTC)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>End (UTC)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Lunch (min)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Timer</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Total</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Delete</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => {
            const hasStartTime = isValidTimeInput(row.startTime);
            const hasEndTime = isValidTimeInput(row.endTime);
            const showStop = hasStartTime && !hasEndTime;
            const total = rowRoundedTotal(row);
            return (
              <tr key={row.date} style={{ borderBottom: "1px solid #eee" }}>
                <td style={{ padding: "8px" }}>{row.weekday}</td>
                <td style={{ padding: "8px" }}>{formatDateUTC(`${row.date}T00:00:00Z`)}</td>
                <td style={{ padding: "8px" }}>
                  <input
                    type="time"
                    value={row.startTime}
                    onChange={(event) => {
                      const nextRow = { ...row, startTime: event.target.value };
                      updateRow(row.date, { startTime: event.target.value });
                      maybeAutoSave(nextRow);
                    }}
                  />
                </td>
                <td style={{ padding: "8px" }}>
                  <input
                    type="time"
                    value={row.endTime}
                    onChange={(event) => {
                      const nextRow = { ...row, endTime: event.target.value };
                      updateRow(row.date, { endTime: event.target.value });
                      maybeAutoSave(nextRow);
                    }}
                  />
                </td>
                <td style={{ padding: "8px" }}>
                  <input
                    type="number"
                    min={0}
                    step={1}
                    value={row.lunchDuration}
                    onChange={(event) =>
                      updateRow(row.date, { lunchDuration: Number(event.target.value) })
                    }
                    style={{ width: "80px" }}
                  />
                </td>
                <td style={{ padding: "8px" }}>
                  {showStop ? (
                    <button type="button" onClick={() => handleStop(row.date)} disabled={busy}>
                      Stop
                    </button>
                  ) : (
                    <button
                      type="button"
                      onClick={() => handleStart(row.date)}
                      disabled={busy}
                    >
                      Start
                    </button>
                  )}
                </td>
                <td style={{ padding: "8px" }}>
                  {total === null ? "-" : `${formatDuration(total)} (${(total / 60).toFixed(2)}h)`}
                </td>
                <td style={{ padding: "8px" }}>
                  <button type="button" onClick={() => deleteRow(row)} disabled={busy}>
                    Delete
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>

      <div style={{ marginTop: "16px", display: "flex", justifyContent: "flex-start" }}>
        <aside style={{ border: "1px solid #ddd", borderRadius: "6px", padding: "12px", minWidth: "320px" }}>
          <h3 style={{ marginTop: 0 }}>Totals ({totals.monthLabel})</h3>
          {totals.weeks.map((week) => (
            <p key={week.key} style={{ margin: "4px 0" }}>
              {week.weekNumber}: {formatHoursFromMinutes(week.minutes)}
            </p>
          ))}
          <p style={{ margin: "8px 0 0 0", fontWeight: 600 }}>
            Month total: {formatHoursFromMinutes(totals.monthTotal)}
          </p>
        </aside>
      </div>
    </main>
  );
};

export default WeekPlanner;
