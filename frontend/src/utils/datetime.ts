export function pad2(value: number): string {
  return String(value).padStart(2, "0");
}

export function todayUTCDateInput(): string {
  const now = new Date();
  return `${now.getUTCFullYear()}-${pad2(now.getUTCMonth() + 1)}-${pad2(
    now.getUTCDate()
  )}`;
}

export function nowUTCTimeInput(): string {
  const now = new Date();
  return `${pad2(now.getUTCHours())}:${pad2(now.getUTCMinutes())}`;
}

export function isoToUTCDateInput(iso: string): string {
  const date = new Date(iso);
  return `${date.getUTCFullYear()}-${pad2(date.getUTCMonth() + 1)}-${pad2(
    date.getUTCDate()
  )}`;
}

export function isoToUTCTimeInput(iso: string): string {
  const date = new Date(iso);
  return `${pad2(date.getUTCHours())}:${pad2(date.getUTCMinutes())}`;
}

export function formatDateUTC(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    timeZone: "UTC",
  });
}

export function formatTimeUTC(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleTimeString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: "UTC",
  });
}

export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `${hours}h ${pad2(mins)}m`;
}

export function formatTimerElapsed(startedAtISO: string): string {
  const started = new Date(startedAtISO).getTime();
  const now = Date.now();
  const diffMs = Math.max(0, now - started);
  const totalSeconds = Math.floor(diffMs / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return `${pad2(hours)}:${pad2(minutes)}:${pad2(seconds)}`;
}

export function formatElapsedMs(ms: number): string {
  const totalSeconds = Math.floor(Math.max(0, ms) / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return `${pad2(hours)}:${pad2(minutes)}:${pad2(seconds)}`;
}
