import { apiClient } from "@/lib/api/client";
import { Deliverable, DownloadResponse } from "@/shared/types/api";

function createFileFormData(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return formData;
}

export const deliverablesApi = {
  list: (orderId: string) => apiClient.get<Deliverable[]>(`/orders/${orderId}/deliverables`),
  upload: (orderId: string, file: File) => apiClient.upload<Deliverable>(`/orders/${orderId}/deliverables`, createFileFormData(file)),
  getDownloadUrl: (orderId: string, deliverableId: string) =>
    apiClient.get<DownloadResponse>(`/orders/${orderId}/deliverables/${deliverableId}/download`)
};
