import { CatalogPage } from "./features/catalog/CatalogPage";
import "./styles/app.css";

function App() {
  return (
    <div className="page">
      <header>
        <h1>Sailorport</h1>
        <p>Software catalog</p>
      </header>
      <CatalogPage />
    </div>
  );
}

export default App;
