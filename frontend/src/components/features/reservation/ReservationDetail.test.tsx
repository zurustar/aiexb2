import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { ReservationDetail } from "./ReservationDetail";
import { Reservation } from "@/types/models";

const reservation: Reservation = {
  id: "r1",
  organizerId: "u1",
  title: "企画レビュー",
  description: "成果物の確認",
  startAt: "2025-11-24T09:00:00Z",
  endAt: "2025-11-24T10:00:00Z",
  isPrivate: true,
  timezone: "Asia/Tokyo",
  approvalStatus: "PENDING",
  version: 1,
  createdAt: "2025-11-23T00:00:00Z",
  updatedAt: "2025-11-23T00:00:00Z",
};

describe("ReservationDetail", () => {
  it("renders reservation information", () => {
    render(
      <ReservationDetail reservation={reservation} isOpen onClose={jest.fn()} />
    );

    expect(screen.getByText("企画レビュー")).toBeInTheDocument();
    expect(screen.getByText("成果物の確認")).toBeInTheDocument();
    expect(screen.getByText("非公開")).toBeInTheDocument();
    expect(screen.getByTestId("start-at").textContent).toContain("2025");
    expect(screen.getByTestId("end-at").textContent).toContain("2025");
  });

  it("calls edit and cancel handlers", () => {
    const handleEdit = jest.fn();
    const handleCancel = jest.fn();

    render(
      <ReservationDetail
        reservation={reservation}
        isOpen
        onClose={jest.fn()}
        onEdit={handleEdit}
        onCancel={handleCancel}
      />
    );

    fireEvent.click(screen.getByRole("button", { name: "編集" }));
    fireEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(handleEdit).toHaveBeenCalledWith(reservation);
    expect(handleCancel).toHaveBeenCalledWith(reservation);
  });
});
