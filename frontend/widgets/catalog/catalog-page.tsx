import { CatalogResults } from "@/widgets/catalog/catalog-results";
import { CatalogToolbar } from "@/widgets/catalog/catalog-toolbar";

export function CatalogPage() {
  return (
    <main className="page-frame">
      <div className="container space-y-8 py-10">
        <div className="max-w-3xl space-y-3">
          <p className="text-sm uppercase tracking-[0.2em] text-primary">Публичный каталог</p>
          <h1>Предложения и запросы в единой музыкальной витрине.</h1>
          <p>Ищите карточки по реальному backend контракту, фильтруйте результаты и переходите к сделке без лишних экранов.</p>
        </div>

        <CatalogToolbar />
        <CatalogResults />
      </div>
    </main>
  );
}
