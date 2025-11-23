import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { Button } from "./Button";

describe("Button", () => {
  it("renders label and handles click", () => {
    const handleClick = jest.fn();

    render(<Button onClick={handleClick}>保存</Button>);

    const button = screen.getByRole("button", { name: "保存" });
    fireEvent.click(button);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("applies variant and size classes", () => {
    render(
      <Button variant="secondary" size="lg">
        詳細
      </Button>
    );

    const button = screen.getByRole("button", { name: "詳細" });
    expect(button.className).toContain("border-gray-300");
    expect(button.className).toContain("text-base");
  });

  it("shows loading state and disables interaction", () => {
    const handleClick = jest.fn();

    render(
      <Button isLoading onClick={handleClick}>
        読み込み
      </Button>
    );

    const button = screen.getByRole("button", { name: "処理中..." });
    expect(button).toBeDisabled();
    fireEvent.click(button);
    expect(handleClick).not.toHaveBeenCalled();
  });
});
