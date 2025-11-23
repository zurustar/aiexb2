import { formatDateTime, isWithinBusinessHours, toTimeZone, validateEmail } from "./utils";

describe("utils", () => {
  it("formats date string with timezone", () => {
    const iso = "2025-01-01T09:00:00Z";
    const formatted = formatDateTime(iso, "en-US", "Asia/Tokyo");
    expect(formatted).toContain("2025");
  });

  it("converts time to specified timezone", () => {
    const source = new Date("2025-01-01T00:00:00Z");
    const converted = toTimeZone(source, "Asia/Tokyo");
    expect(converted.getHours()).not.toBe(source.getUTCHours());
  });

  it("validates email addresses", () => {
    expect(validateEmail("user@example.com")).toBe(true);
    expect(validateEmail("invalid-email")).toBe(false);
  });

  it("checks business hours respecting timezone", () => {
    const date = "2025-01-01T09:30:00Z";
    expect(
      isWithinBusinessHours(date, {
        startHour: 17,
        endHour: 20,
        timeZone: "Asia/Tokyo",
      })
    ).toBe(true);
    expect(
      isWithinBusinessHours(date, {
        startHour: 10,
        endHour: 12,
        timeZone: "UTC",
      })
    ).toBe(false);
  });
});
