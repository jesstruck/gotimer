export interface TimeEntry {
  id: number;
  start_time: string;
  end_time: string;
  lunch_duration: number;
  source: "manual" | "timer" | string;
  created_at: string;
  updated_at: string;
  worked_minutes: number;
  rounded_worked_minutes: number;
  rounded_worked_hours: number;
}

export interface TimerState {
  active: boolean;
  started_at?: string;
}

export interface SummaryResponse {
  period: "weekly" | "monthly" | string;
  anchor_date: string;
  range_start: string;
  range_end: string;
  total_entries: number;
  total_worked_minutes: number;
  total_rounded_minutes: number;
  total_rounded_hours: number;
}

export interface TimeEntryPayload {
  date: string;
  start_time: string;
  end_time: string;
  lunch_duration: number;
}
