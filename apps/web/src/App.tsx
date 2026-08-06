import { useState } from "react";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { ScaffoldPanel } from "./features/scaffold/ScaffoldPanel";
import "./styles/app.css";

function App() {
  const [catalogTick, setCatalogTick] = useState(0);

  return (
    <div className="page">
      <header>
        <h1>Sailorport</h1>
        <p>Software catalog & golden path</p>
      </header>
      <ScaffoldPanel onSuccess={() => setCatalogTick((n) => n + 1)} />
      <CatalogPage refreshToken={catalogTick} />
    </div>
  );
}

export default App;
