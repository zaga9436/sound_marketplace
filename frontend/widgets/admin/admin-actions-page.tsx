"use client";

import { useQuery } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { getErrorMessage } from "@/lib/api/errors";
import { AdminPageHeader, AdminSection, StatusBadge } from "@/widgets/admin/admin-ui";

export function AdminActionsPage() {
  const actionsQuery = useQuery({
    queryKey: ["admin", "actions", "page"],
    queryFn: () => adminApi.listActions({ limit: "100" })
  });

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Журнал действий" description="История moderation actions: кто, что и когда сделал в административной части продукта." />

      <AdminSection className="overflow-hidden p-0">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-slate-50 text-left text-slate-500">
              <tr>
                <th className="px-5 py-4 font-medium">Когда</th>
                <th className="px-5 py-4 font-medium">Администратор</th>
                <th className="px-5 py-4 font-medium">Цель</th>
                <th className="px-5 py-4 font-medium">Действие</th>
                <th className="px-5 py-4 font-medium">Причина</th>
              </tr>
            </thead>
            <tbody>
              {(actionsQuery.data ?? []).map((action) => (
                <tr key={action.id} className="border-t border-slate-200">
                  <td className="px-5 py-4 text-slate-700">{new Date(action.created_at).toLocaleString("ru-RU")}</td>
                  <td className="px-5 py-4 font-mono text-xs text-slate-600">{action.admin_user_id}</td>
                  <td className="px-5 py-4">
                    <div className="space-y-1">
                      <StatusBadge tone="slate">{action.target_type}</StatusBadge>
                      <p className="font-mono text-xs text-slate-500">{action.target_id}</p>
                    </div>
                  </td>
                  <td className="px-5 py-4 text-slate-900">{action.action}</td>
                  <td className="px-5 py-4 text-slate-700">{action.reason || "Без причины"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {!actionsQuery.isLoading && (actionsQuery.data?.length ?? 0) === 0 ? (
          <div className="p-6 text-sm text-slate-500">Журнал пока пуст.</div>
        ) : null}
      </AdminSection>

      {actionsQuery.isError ? <p className="text-sm text-destructive">{getErrorMessage(actionsQuery.error)}</p> : null}
    </div>
  );
}
