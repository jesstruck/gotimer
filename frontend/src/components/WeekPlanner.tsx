import React, { useEffect, useMemo, useState } from "react";
import { TimeEntry, TimeEntryPayload } from "../types";
import {
  formatDateUTC,
  formatElapsedMs,
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

type RunningTimer = {
  date: string;
  startedAtMs: number;
} | null;

type PlannerTheme = {
  id: string;
  name: string;
  logo: string;
  background: string;
  panel: string;
  accent: string;
  accentSoft: string;
  border: string;
  text: string;
  mutedText: string;
  inputBg: string;
  success: string;
  error: string;
};

const THEMES: PlannerTheme[] = [
  {
    id: "sea-turtle",
    name: "Sea Turtle",
    logo: "🐢",
    background: "linear-gradient(145deg, #0f3d2e 0%, #1f6f54 55%, #d8c58b 100%)",
    panel: "#f5f1e4",
    accent: "#1f6f54",
    accentSoft: "#d6e6dd",
    border: "#8ca690",
    text: "#1e2f28",
    mutedText: "#4f6459",
    inputBg: "#fffef8",
    success: "#11643a",
    error: "#9f2f2f",
  },
  {
    id: "blue-whale",
    name: "Blue Whale",
    logo: "🐋",
    background: "linear-gradient(150deg, #05233f 0%, #0e4f7e 60%, #9ed0ea 100%)",
    panel: "#eef7fc",
    accent: "#0d4f7d",
    accentSoft: "#d3e7f3",
    border: "#89acc2",
    text: "#12324a",
    mutedText: "#3f6076",
    inputBg: "#f8fdff",
    success: "#0b5b87",
    error: "#8c2333",
  },
];

function isKnownTheme(themeID: string): boolean {
  return THEMES.some((theme) => theme.id === themeID);
}

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

const WeekPlanner: React.FC = () => {
  const [weekOffset, setWeekOffset] = useState<number>(0);
  const weekTemplate = useMemo(() => buildWeekRows(weekOffset), [weekOffset]);
  const [rows, setRows] = useState<PlannerRow[]>(weekTemplate);
  const [allEntries, setAllEntries] = useState<TimeEntry[]>([]);
  const [runningTimer, setRunningTimer] = useState<RunningTimer>(null);
  const [tick, setTick] = useState<number>(Date.now());
  const [busy, setBusy] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [message, setMessage] = useState<string>("");
  const [themeId, setThemeId] = useState<string>(THEMES[0].id);
  const theme = useMemo(
    () => THEMES.find((candidate) => candidate.id === themeId) || THEMES[0],
    [themeId]
  );

  const buttonStyle: React.CSSProperties = {
    borderRadius: "10px",
    border: `1px solid ${theme.border}`,
    background: theme.accentSoft,
    color: theme.text,
    fontWeight: 600,
    cursor: "pointer",
    padding: "8px 12px",
  };

  const startStopButtonStyle: React.CSSProperties = {
    ...buttonStyle,
    background: theme.accent,
    color: "#ffffff",
  };

  const deleteButtonStyle: React.CSSProperties = {
    ...buttonStyle,
    background: "#f7d9d9",
    border: "1px solid #d28a8a",
    color: "#6d1f1f",
  };

  const inputStyle: React.CSSProperties = {
    borderRadius: "8px",
    border: `1px solid ${theme.border}`,
    background: theme.inputBg,
    color: theme.text,
    padding: "6px 8px",
  };

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

  useEffect(() => {
    if (!runningTimer) {
      return;
    }
    const id = window.setInterval(() => setTick(Date.now()), 1000);
    return () => window.clearInterval(id);
  }, [runningTimer]);

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
            endTime: isoToUTCTimeInput(entry.end_time),
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

  useEffect(() => {
    let active = true;

    const loadThemePreference = async (): Promise<void> => {
      try {
        const response = await fetch("/api/preferences/theme");
        const payload: { theme?: string } = await parseJSON(response);
        if (!active) {
          return;
        }

        const savedTheme = typeof payload.theme === "string" ? payload.theme.trim() : "";
        if (savedTheme && isKnownTheme(savedTheme)) {
          setThemeId(savedTheme);
        }
      } catch {
        // Theme defaults locally if preference endpoint is unavailable.
      }
    };

    void loadThemePreference();

    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    // Week changes reset active in-memory timer assignment for clarity.
    setRunningTimer(null);
  }, [weekOffset]);

  const updateRow = (date: string, patch: Partial<PlannerRow>): void => {
    setRows((prevRows) => prevRows.map((row) => (row.date === date ? { ...row, ...patch } : row)));
  };

  const persistThemePreference = async (nextThemeID: string): Promise<void> => {
    if (!isKnownTheme(nextThemeID)) {
      return;
    }
    try {
      const response = await fetch("/api/preferences/theme", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ theme: nextThemeID }),
      });
      await parseJSON(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save theme.");
    }
  };

  const handleThemeChange = (nextThemeID: string): void => {
    setThemeId(nextThemeID);
    void persistThemePreference(nextThemeID);
  };

  const handleStart = (date: string): void => {
    setMessage("");
    setError("");
    const nowTime = nowUTCTimeInput();
    setRunningTimer({ date, startedAtMs: Date.now() });
    updateRow(date, { startTime: nowTime, endTime: "" });
  };

  const handleStop = (date: string): void => {
    setMessage("");
    setError("");
    const nowTime = nowUTCTimeInput();
    updateRow(date, { endTime: nowTime });
    setRunningTimer((prev) => (prev?.date === date ? null : prev));
  };

  const saveRow = async (row: PlannerRow): Promise<void> => {
    setBusy(true);
    setError("");
    setMessage("");
    try {
      if (!row.startTime || !row.endTime) {
        throw new Error("Start and end time are required to save a day.");
      }

      const payload: TimeEntryPayload = {
        date: row.date,
        start_time: row.startTime,
        end_time: row.endTime,
        lunch_duration: row.lunchDuration,
      };

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
      setMessage(`Saved ${row.weekday}.`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save day.");
    } finally {
      setBusy(false);
    }
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

  const runningElapsed = runningTimer ? formatElapsedMs(tick - runningTimer.startedAtMs) : "00:00:00";

  return (
    <main
      style={{
        maxWidth: "1100px",
        margin: "16px auto",
        padding: "20px",
        borderRadius: "20px",
        background: theme.panel,
        color: theme.text,
        border: `1px solid ${theme.border}`,
        boxShadow: "0 16px 44px rgba(0, 0, 0, 0.2)",
        fontFamily: "'Trebuchet MS', 'Gill Sans', sans-serif",
      }}
    >
      <div
        style={{
          background: theme.background,
          borderRadius: "16px",
          padding: "14px 16px",
          marginBottom: "16px",
          color: "#ffffff",
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          gap: "10px",
          flexWrap: "wrap",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <span style={{ fontSize: "32px", lineHeight: 1 }} aria-hidden>
            {theme.logo}
          </span>
          <strong style={{ fontSize: "20px" }}>{theme.name}</strong>
        </div>
        <label style={{ display: "flex", alignItems: "center", gap: "8px", fontWeight: 700 }}>
          Theme
          <select
            value={themeId}
            onChange={(event) => handleThemeChange(event.target.value)}
            style={{ ...inputStyle, minWidth: "140px", color: theme.text }}
          >
            {THEMES.map((option) => (
              <option key={option.id} value={option.id}>
                {option.logo} {option.name}
              </option>
            ))}
          </select>
        </label>
      </div>

      <h1 style={{ marginTop: 0 }}>{weekTitle}</h1>
      <div style={{ display: "flex", gap: "8px", alignItems: "center", flexWrap: "wrap" }}>
        <button
          type="button"
          onClick={() => setWeekOffset((prev) => prev - 1)}
          disabled={busy}
          style={buttonStyle}
        >
          {"< Previous Week"}
        </button>
        <button
          type="button"
          onClick={() => setWeekOffset((prev) => prev + 1)}
          disabled={busy}
          style={buttonStyle}
        >
          {"Next Week >"}
        </button>
        <button
          type="button"
          onClick={() => setWeekOffset(0)}
          disabled={busy || weekOffset === 0}
          style={buttonStyle}
        >
          Current Week
        </button>
        <small style={{ color: theme.mutedText }}>{weekNumberLabel}</small>
      </div>
      <p style={{ color: theme.mutedText }}>
        Only Monday-Sunday is shown. Enter start/end/lunch or use per-day start/stop, then save.
      </p>
      {error && <p style={{ color: theme.error }}>{error}</p>}
      {message && <p style={{ color: theme.success }}>{message}</p>}

      <table
        style={{
          width: "100%",
          borderCollapse: "collapse",
          marginTop: "12px",
          background: theme.inputBg,
          borderRadius: "12px",
          overflow: "hidden",
        }}
      >
        <thead>
          <tr style={{ borderBottom: `2px solid ${theme.border}`, background: theme.accentSoft }}>
            <th style={{ textAlign: "left", padding: "8px" }}>Day</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Date</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Start (UTC)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>End (UTC)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Lunch (min)</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Timer</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Total</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Save</th>
            <th style={{ textAlign: "left", padding: "8px" }}>Delete</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => {
            const isRunning = runningTimer?.date === row.date;
            const total = rowRoundedTotal(row);
            return (
              <tr key={row.date} style={{ borderBottom: `1px solid ${theme.border}` }}>
                <td style={{ padding: "8px" }}>{row.weekday}</td>
                <td style={{ padding: "8px" }}>{formatDateUTC(`${row.date}T00:00:00Z`)}</td>
                <td style={{ padding: "8px" }}>
                  <input
                    type="time"
                    value={row.startTime}
                    onChange={(event) => updateRow(row.date, { startTime: event.target.value })}
                    style={inputStyle}
                  />
                </td>
                <td style={{ padding: "8px" }}>
                  <input
                    type="time"
                    value={row.endTime}
                    onChange={(event) => updateRow(row.date, { endTime: event.target.value })}
                    style={inputStyle}
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
                    style={{ ...inputStyle, width: "80px" }}
                  />
                </td>
                <td style={{ padding: "8px" }}>
                  {isRunning ? (
                    <div style={{ display: "flex", gap: "6px", alignItems: "center" }}>
                      <button
                        type="button"
                        onClick={() => handleStop(row.date)}
                        disabled={busy}
                        style={startStopButtonStyle}
                      >
                        Stop
                      </button>
                      <span>{runningElapsed}</span>
                    </div>
                  ) : (
                    <button
                      type="button"
                      onClick={() => handleStart(row.date)}
                      disabled={busy || Boolean(runningTimer)}
                      style={startStopButtonStyle}
                    >
                      Start
                    </button>
                  )}
                </td>
                <td style={{ padding: "8px" }}>
                  {total === null ? "-" : `${formatDuration(total)} (${(total / 60).toFixed(2)}h)`}
                </td>
                <td style={{ padding: "8px" }}>
                  <button type="button" onClick={() => saveRow(row)} disabled={busy} style={buttonStyle}>
                    Save
                  </button>
                </td>
                <td style={{ padding: "8px" }}>
                  <button
                    type="button"
                    onClick={() => deleteRow(row)}
                    disabled={busy}
                    style={deleteButtonStyle}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>

      <div style={{ marginTop: "16px", display: "flex", justifyContent: "flex-start" }}>
        <aside
          style={{
            border: `1px solid ${theme.border}`,
            borderRadius: "10px",
            padding: "12px",
            minWidth: "320px",
            background: theme.accentSoft,
          }}
        >
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
