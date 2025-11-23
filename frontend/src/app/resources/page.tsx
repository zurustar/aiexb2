"use client";

import React, { useMemo, useState } from "react";

import { Button } from "@/components/ui/Button";
import { useResources, UseResourcesResult } from "@/hooks/useResources";
import { Resource, ResourceType, Role } from "@/types/models";

export type ResourcesPageProps = {
  useResourcesHook?: (params?: Parameters<typeof useResources>[0]) => UseResourcesResult;
};

export const ResourcesPage: React.FC<ResourcesPageProps> = ({ useResourcesHook = useResources }) => {
  const resourcesState = useResourcesHook({ autoLoad: true });
  const { resources, isLoading, error, search } = resourcesState;

  const [keyword, setKeyword] = useState("");
  const [resourceType, setResourceType] = useState<ResourceType | "">("");
  const [capacity, setCapacity] = useState("");
  const [role, setRole] = useState<Role | "">("");

  const visibleResources = useMemo(() => resources, [resources]);

  const handleSearch = async (event: React.FormEvent) => {
    event.preventDefault();
    await search({
      keyword: keyword || undefined,
      type: (resourceType as ResourceType) || undefined,
      capacity: capacity ? Number(capacity) : undefined,
      requiredRole: (role as Role) || undefined,
    });
  };

  return (
    <main className="flex flex-col gap-6" aria-label="リソース管理">
      <header className="flex flex-col gap-2">
        <p className="text-sm font-semibold text-blue-700">リソース管理</p>
        <h1 className="text-2xl font-bold text-gray-900">会議室・備品の検索</h1>
        <p className="text-sm text-gray-700">キーワードや種別で空きリソースを探します。</p>
      </header>

      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}

      <form onSubmit={handleSearch} className="grid grid-cols-1 gap-3 rounded-lg border border-gray-200 bg-white p-5 shadow-sm md:grid-cols-4">
        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="keyword">
          キーワード
          <input
            id="keyword"
            name="keyword"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
            placeholder="会議室、備品名など"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="type">
          種別
          <select
            id="type"
            name="type"
            value={resourceType}
            onChange={(e) => setResourceType(e.target.value as ResourceType | "")}
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          >
            <option value="">指定なし</option>
            <option value="MEETING_ROOM">会議室</option>
            <option value="EQUIPMENT">備品</option>
          </select>
        </label>

        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="capacity">
          収容人数
          <input
            id="capacity"
            name="capacity"
            type="number"
            min="1"
            value={capacity}
            onChange={(e) => setCapacity(e.target.value)}
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
            placeholder="10"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm font-medium text-gray-800" htmlFor="role">
          必要ロール
          <select
            id="role"
            name="role"
            value={role}
            onChange={(e) => setRole(e.target.value as Role | "")}
            className="rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200"
          >
            <option value="">指定なし</option>
            <option value="GENERAL">一般</option>
            <option value="SECRETARY">秘書</option>
            <option value="MANAGER">マネージャー</option>
            <option value="ADMIN">管理者</option>
            <option value="AUDITOR">監査</option>
          </select>
        </label>

        <div className="md:col-span-4 flex justify-end">
          <Button type="submit" isLoading={isLoading} disabled={isLoading}>
            検索
          </Button>
        </div>
      </form>

      <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm" aria-label="リソース一覧">
        {visibleResources.length === 0 ? (
          <p className="text-sm text-gray-700">条件に一致するリソースがありません。</p>
        ) : (
          <ul className="divide-y divide-gray-200" data-testid="resource-list">
            {visibleResources.map((resource: Resource) => (
              <li key={resource.id} className="flex items-start justify-between py-3">
                <div>
                  <p className="font-semibold text-gray-900">{resource.name}</p>
                  <p className="text-xs text-gray-500">{resource.type} {resource.capacity ? ` / ${resource.capacity}名` : ""}</p>
                  {resource.location && <p className="text-xs text-gray-500">{resource.location}</p>}
                </div>
                {resource.requiredRole && (
                  <span className="rounded-full bg-purple-50 px-2 py-1 text-xs font-semibold text-purple-700">{resource.requiredRole}</span>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
};

export default function Page() {
  return <ResourcesPage />;
}
