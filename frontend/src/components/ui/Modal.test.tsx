import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { Modal } from "./Modal";

describe("Modal", () => {
  it("renders content when open", () => {
    render(
      <Modal isOpen title="モーダルタイトル" onClose={jest.fn()}>
        <p>内容</p>
      </Modal>
    );

    expect(screen.getByText("モーダルタイトル")).toBeInTheDocument();
    expect(screen.getByText("内容")).toBeInTheDocument();
  });

  it("invokes onClose when overlay clicked", () => {
    const handleClose = jest.fn();

    render(
      <Modal isOpen title="タイトル" onClose={handleClose}>
        <p>body</p>
      </Modal>
    );

    fireEvent.click(screen.getByTestId("modal-overlay"));

    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it("can disable overlay close and still close via button", () => {
    const handleClose = jest.fn();

    render(
      <Modal isOpen title="タイトル" onClose={handleClose} closeOnOverlayClick={false}>
        <p>body</p>
      </Modal>
    );

    fireEvent.click(screen.getByTestId("modal-overlay"));
    expect(handleClose).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "閉じる" }));
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
