import type { Service } from "./types";

type ServiceListProps = {
  services: Service[];
  loading: boolean;
  error: string;
  onRefresh: () => void;
  onEdit: (svc: Service) => void;
  onDelete: (svc: Service) => void;
};

export function ServiceList({
  services,
  loading,
  error,
  onRefresh,
  onEdit,
  onDelete,
}: ServiceListProps) {
  return (
    <section className="card">
      <div className="row">
        <h2>Services</h2>
        <button type="button" onClick={onRefresh} disabled={loading}>
          Refresh
        </button>
      </div>

      {error && <p className="error">{error}</p>}
      {loading && <p>Memuat...</p>}

      {!loading && services.length === 0 && (
        <p className="muted">Belum ada service. Buat yang pertama di atas.</p>
      )}

      <ul className="list">
        {services.map((svc) => (
          <li key={svc.id}>
            <div className="item-main">
              <strong>{svc.name}</strong>
              <span className="muted">
                {svc.owner || "no owner"} · {svc.description || "no description"}
              </span>
            </div>
            <div className="item-actions">
              <button
                type="button"
                className="button-secondary"
                onClick={() => onEdit(svc)}
              >
                Edit
              </button>
              <button
                type="button"
                className="button-danger"
                onClick={() => onDelete(svc)}
              >
                Hapus
              </button>
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}
