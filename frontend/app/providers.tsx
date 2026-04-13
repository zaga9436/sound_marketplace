"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { PropsWithChildren, useState } from "react";

import { AuthBootstrap } from "@/lib/auth/auth-bootstrap";
import { createQueryClient } from "@/lib/query/query-client";

export function Providers({ children }: PropsWithChildren) {
  const [queryClient] = useState(() => createQueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <AuthBootstrap />
      {children}
    </QueryClientProvider>
  );
}
