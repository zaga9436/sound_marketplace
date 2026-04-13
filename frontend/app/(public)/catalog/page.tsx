import { Suspense } from "react";

import { CatalogPage } from "@/widgets/catalog/catalog-page";

export default function PublicCatalogPage() {
  return (
    <Suspense fallback={null}>
      <CatalogPage />
    </Suspense>
  );
}
