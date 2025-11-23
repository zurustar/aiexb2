import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import { Sidebar, SidebarItem } from "./Sidebar";
import { Role } from "@/types/models";

describe("Sidebar", () => {
  const items: SidebarItem[] = [
    { label: "ダッシュボード", href: "/dashboard", roles: ["ADMIN"] },
    { label: "予約", href: "/events" },
  ];

  it("filters items based on role and renders badges", () => {
    const auth = {
      isAuthenticated: true,
      hasRole: (role: Role | Role[]) => {
        if (Array.isArray(role)) return role.includes("ADMIN");
        return role === "ADMIN";
      },
    };

    render(
      <Sidebar
        items={[
          ...items,
          { label: "承認", href: "/approvals", badge: 3, roles: ["MANAGER"] },
          { label: "通知", href: "/notifications", badge: 5 },
        ]}
        auth={auth}
      />
    );

    expect(screen.getByText("ダッシュボード")).toBeInTheDocument();
    expect(screen.getByText("予約")).toBeInTheDocument();
    expect(screen.queryByText("承認")).not.toBeInTheDocument();
    expect(screen.getByLabelText("badge")).toHaveTextContent("5");
  });

  it("invokes selection handler when item is clicked", () => {
    const handleSelect = jest.fn();

    render(<Sidebar items={items} auth={{ isAuthenticated: true, hasRole: () => true }} onItemSelect={handleSelect} />);

    fireEvent.click(screen.getByText("予約"));
    expect(handleSelect).toHaveBeenCalledWith(items[1]);
  });
});
