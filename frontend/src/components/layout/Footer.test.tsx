import React from "react";
import { render, screen } from "@testing-library/react";

import { Footer } from "./Footer";

describe("Footer", () => {
  it("renders company, year, version, and links", () => {
    const year = new Date().getFullYear();

    render(
      <Footer
        companyName="Example Corp"
        version="1.2.3"
        links={[
          { label: "利用規約", href: "/terms" },
          { label: "プライバシー", href: "/privacy" },
        ]}
      />
    );

    expect(screen.getByText("Example Corp")).toBeInTheDocument();
    expect(screen.getByText(`© ${year}`)).toBeInTheDocument();
    expect(screen.getByText("v1.2.3")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "利用規約" })).toHaveAttribute("href", "/terms");
  });
});
