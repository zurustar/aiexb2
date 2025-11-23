import React from "react";
import { act, render, screen } from "@testing-library/react";

import { Toast } from "./Toast";

describe("Toast", () => {
  jest.useFakeTimers();

  it("shows message and variant icon", () => {
    render(<Toast message="保存しました" variant="success" />);

    expect(screen.getByRole("status")).toHaveTextContent("保存しました");
    expect(screen.getByText("✔")).toBeInTheDocument();
  });

  it("auto closes after duration", () => {
    const handleClose = jest.fn();
    render(<Toast message="完了" onClose={handleClose} duration={2000} />);

    act(() => {
      jest.advanceTimersByTime(2000);
    });

    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it("does not render when not visible", () => {
    render(<Toast message="hidden" isVisible={false} />);

    expect(screen.queryByRole("status")).toBeNull();
  });
});
