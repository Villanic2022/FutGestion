import React from 'react';
import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import { FaUsers, FaMoneyBillWave, FaList } from 'react-icons/fa';
import PlayersPage from './pages/PlayersPage';
import ConceptsPage from './pages/ConceptsPage';
import PaymentsPage from './pages/PaymentsPage';

function NavItem({ to, icon: Icon, children }) {
  const location = useLocation();
  const isActive = location.pathname === to;
  return (
    <Link
      to={to}
      className={`flex items-center space-x-2 px-4 py-3 rounded-lg transition-all duration-200 ${
        isActive 
          ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/30' 
          : 'text-slate-400 hover:bg-slate-800 hover:text-white'
      }`}
    >
      <Icon className="text-xl" />
      <span className="font-medium">{children}</span>
    </Link>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="flex bg-slate-950 text-slate-100 min-h-screen font-sans">
        {/* Sidebar */}
        <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col">
          <div className="p-6">
            <h1 className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-emerald-400 flex items-center gap-2">
              ⚽ FutGesti&oacute;n
            </h1>
          </div>
          <nav className="flex-1 px-4 space-y-2">
            <NavItem to="/" icon={FaMoneyBillWave}>Control de Pagos</NavItem>
            <NavItem to="/jugadores" icon={FaUsers}>Jugadores</NavItem>
            <NavItem to="/conceptos" icon={FaList}>Conceptos de Pago</NavItem>
          </nav>
          <div className="p-4 text-xs text-slate-500 text-center">
            v1.0.0
          </div>
        </aside>

        {/* Main Content */}
        <main className="flex-1 overflow-x-hidden overflow-y-auto">
          <Routes>
            <Route path="/" element={<PaymentsPage />} />
            <Route path="/jugadores" element={<PlayersPage />} />
            <Route path="/conceptos" element={<ConceptsPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
