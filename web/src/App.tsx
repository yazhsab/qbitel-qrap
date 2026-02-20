import React from "react";
import { BrowserRouter, Routes, Route, NavLink } from "react-router-dom";
import { DashboardPage } from "./pages/DashboardPage.js";
import { AssessmentsPage } from "./pages/AssessmentsPage.js";
import { FindingsPage } from "./pages/FindingsPage.js";
import { HndlPage } from "./pages/HndlPage.js";

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <div className="qtn-app">
        <nav className="qtn-nav">
          <div className="qtn-nav__brand">QRAP Dashboard</div>
          <div className="qtn-nav__links">
            <NavLink to="/" end>
              Dashboard
            </NavLink>
            <NavLink to="/assessments">Assessments</NavLink>
            <NavLink to="/findings">Findings</NavLink>
            <NavLink to="/hndl">HNDL Analysis</NavLink>
          </div>
        </nav>
        <main className="qtn-main">
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/assessments" element={<AssessmentsPage />} />
            <Route path="/findings" element={<FindingsPage />} />
            <Route path="/hndl" element={<HndlPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
};
