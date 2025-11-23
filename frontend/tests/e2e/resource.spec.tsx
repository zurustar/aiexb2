import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ResourcesPage } from "@/app/resources/page";
import { UseResourcesResult } from "@/hooks/useResources";
import { Resource } from "@/types/models";

const resources: Resource[] = [
  {
    id: "room-1",
    name: "第1会議室",
    type: "MEETING_ROOM",
    capacity: 10,
    location: "東京オフィス 10F",
    requiredRole: "GENERAL",
    isActive: true,
    createdAt: "2025-11-01T00:00:00Z",
    updatedAt: "2025-11-01T00:00:00Z",
  },
];

describe("E2E: リソース検索フロー", () => {
  it("submits search conditions and shows list", async () => {
    const search = jest.fn();
    const mockState: UseResourcesResult = {
      resources,
      isLoading: false,
      error: null,
      search,
      checkAvailability: jest.fn(),
    };

    const useResourcesHook = () => mockState;

    render(<ResourcesPage useResourcesHook={useResourcesHook} />);

    fireEvent.change(screen.getByLabelText("キーワード"), { target: { value: "会議室" } });
    fireEvent.change(screen.getByLabelText("種別"), { target: { value: "MEETING_ROOM" } });
    fireEvent.change(screen.getByLabelText("収容人数"), { target: { value: "8" } });
    fireEvent.change(screen.getByLabelText("必要ロール"), { target: { value: "GENERAL" } });
    fireEvent.click(screen.getByText("検索"));

    await waitFor(() => expect(search).toHaveBeenCalledWith({
      keyword: "会議室",
      type: "MEETING_ROOM",
      capacity: 8,
      requiredRole: "GENERAL",
    }));

    expect(screen.getByTestId("resource-list")).toBeInTheDocument();
    expect(screen.getByText("第1会議室")).toBeInTheDocument();
  });
});
