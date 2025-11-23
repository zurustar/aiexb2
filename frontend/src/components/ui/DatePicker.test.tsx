import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { DatePicker } from "./DatePicker";

describe("DatePicker", () => {
  it("renders label and helper", () => {
    render(<DatePicker label="開始" helperText="開始日時を選択" name="start" />);

    expect(screen.getByText("開始")).toBeInTheDocument();
    expect(screen.getByText("開始日時を選択")).toBeInTheDocument();
  });

  it("fires onChange with the chosen value", () => {
    const handleChange = jest.fn();
    render(<DatePicker label="開始日時" onChange={handleChange} name="start" />);

    const input = screen.getByLabelText("開始日時");
    fireEvent.change(input, { target: { value: "2025-11-24T10:30" } });

    expect(handleChange).toHaveBeenCalledWith("2025-11-24T10:30");
  });

  it("shows error message and sets border color", () => {
    render(<DatePicker label="終了日時" error="必須項目です" name="end" />);

    expect(screen.getByRole("alert")).toHaveTextContent("必須項目です");
    const input = screen.getByLabelText("終了日時");
    expect(input.className).toContain("border-red-500");
  });
});
