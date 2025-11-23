export type BusinessHours = {
  startHour: number;
  endHour: number;
  timeZone?: string;
};

export const formatDateTime = (input: string | Date, locale = "ja-JP", timeZone?: string): string => {
  const date = typeof input === "string" ? new Date(input) : input;
  return new Intl.DateTimeFormat(locale, {
    dateStyle: "medium",
    timeStyle: "short",
    timeZone,
  }).format(date);
};

export const toTimeZone = (input: string | Date, timeZone: string): Date => {
  const date = typeof input === "string" ? new Date(input) : input;
  const parts = new Intl.DateTimeFormat("en-US", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  })
    .formatToParts(date)
    .reduce<Record<string, string>>((acc, part) => {
      if (part.type !== "literal") acc[part.type] = part.value;
      return acc;
    }, {});

  return new Date(
    `${parts.year}-${parts.month}-${parts.day}T${parts.hour}:${parts.minute}:${parts.second}.000${offsetFromTimeZone(date, timeZone)}`
  );
};

const offsetFromTimeZone = (date: Date, timeZone: string): string => {
  const localDate = new Date(date.toLocaleString("en-US", { timeZone }));
  const diff = (localDate.getTime() - date.getTime()) / 60000;
  const sign = diff >= 0 ? "+" : "-";
  const abs = Math.abs(diff);
  const hours = String(Math.floor(abs / 60)).padStart(2, "0");
  const minutes = String(abs % 60).padStart(2, "0");
  return `${sign}${hours}:${minutes}`;
};

export const validateEmail = (email: string): boolean => {
  const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return pattern.test(email);
};

export const isWithinBusinessHours = (date: string | Date, hours: BusinessHours): boolean => {
  const tzDate = hours.timeZone ? toTimeZone(date, hours.timeZone) : typeof date === "string" ? new Date(date) : date;
  const hour = tzDate.getHours();
  return hour >= hours.startHour && hour < hours.endHour;
};
