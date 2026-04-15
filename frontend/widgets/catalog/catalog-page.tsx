import { CatalogResults } from "@/widgets/catalog/catalog-results";
import { CatalogToolbar } from "@/widgets/catalog/catalog-toolbar";

export function CatalogPage() {
  return (
    <main className="page-frame">
      <div className="container space-y-8 py-10">
        <div className="max-w-3xl space-y-3">
          <p className="text-sm uppercase tracking-[0.2em] text-slate-500">Публичный каталог</p>
          <h1>Музыкальные предложения и запросы в одной витрине.</h1>
          <p>
            Находите готовые биты, услуги звукорежиссеров и задачи от заказчиков. Используйте поиск и фильтры, чтобы быстрее перейти к подходящей
            сделке.
          </p>
        </div>

        <CatalogToolbar />
        <CatalogResults />
      </div>
    </main>
  );
}
