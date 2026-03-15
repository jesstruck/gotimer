import { formatDuration, formatElapsedMs, pad2 } from "./datetime";

describe("datetime utilities", () => {
  test("pad2 left pads single digit numbers", () => {
    expect(pad2(3)).toBe("03");
    expect(pad2(12)).toBe("12");
  });

  test("formatDuration returns hour/minute format", () => {
    expect(formatDuration(0)).toBe("0h 00m");
    expect(formatDuration(75)).toBe("1h 15m");
  });

  test("formatElapsedMs renders HH:MM:SS", () => {
    expect(formatElapsedMs(0)).toBe("00:00:00");
    expect(formatElapsedMs(3723000)).toBe("01:02:03");
  });
});
