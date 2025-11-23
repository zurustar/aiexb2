export type Role =
  | "GENERAL"
  | "SECRETARY"
  | "MANAGER"
  | "ADMIN"
  | "AUDITOR";

export type ResourceType = "MEETING_ROOM" | "EQUIPMENT";

export type ReservationApprovalStatus = "PENDING" | "CONFIRMED" | "REJECTED";

export type ReservationStatus =
  | "CONFIRMED"
  | "CANCELLED"
  | "CHECKED_IN"
  | "COMPLETED"
  | "NO_SHOW";

export type User = {
  id: string;
  sub: string;
  email: string;
  name: string;
  role: Role;
  managerId?: string | null;
  penaltyScore: number;
  penaltyScoreExpireAt?: string | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;
};

export type Resource = {
  id: string;
  name: string;
  type: ResourceType;
  capacity?: number | null;
  location?: string | null;
  equipment?: Record<string, unknown>;
  requiredRole?: Role | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Reservation = {
  id: string;
  organizerId: string;
  title: string;
  description: string;
  startAt: string;
  endAt: string;
  rrule?: string;
  isPrivate: boolean;
  timezone: string;
  approvalStatus: ReservationApprovalStatus;
  updatedBy?: string | null;
  version: number;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;
  organizer?: User;
};

export type ReservationInstance = {
  id: string;
  reservationId: string;
  reservationStartAt: string;
  startAt: string;
  endAt: string;
  originalStartAt?: string | null;
  status: ReservationStatus;
  checkedInAt?: string | null;
  createdAt: string;
  updatedAt: string;
  reservation?: Reservation;
  resources?: Resource[];
  participants?: User[];
};
