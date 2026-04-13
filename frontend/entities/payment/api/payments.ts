import { apiClient } from "@/lib/api/client";
import { BalanceResponse, Payment, PaymentSyncResponse } from "@/shared/types/api";

export const paymentsApi = {
  getBalance: () => apiClient.get<BalanceResponse>("/payments/balance"),
  createDeposit: (amount: number) => apiClient.post<Payment>("/payments/deposits", { amount }),
  sync: (externalId: string) => apiClient.post<PaymentSyncResponse>("/payments/sync", { external_id: externalId })
};
