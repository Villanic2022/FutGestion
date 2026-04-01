import React from 'react';
import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import { FaUsers, FaMoneyBillWave, FaList, FaBars, FaTimes } from 'react-icons/fa';
import PlayersPage from './pages/PlayersPage';
import ConceptsPage from './pages/ConceptsPage';
import PaymentsPage from './pages/PaymentsPage';

function NavItem({ to, icon: Icon, children, onClick }) {
  const location = useLocation();
  const isActive = location.pathname === to;
  return (
    <Link
      to={to}
      onClick={onClick}
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
  const [isSidebarOpen, setIsSidebarOpen] = React.useState(false);

  const closeSidebar = () => setIsSidebarOpen(false);

  return (
    <BrowserRouter>
      <div className="flex bg-slate-950 text-slate-100 min-h-screen font-sans relative">
        {/* Mobile Overlay */}
        {isSidebarOpen && (
          <div 
            className="fixed inset-0 bg-black/50 z-30 lg:hidden backdrop-blur-sm"
            onClick={closeSidebar}
          />
        )}

        {/* Sidebar */}
        <aside className={`
          fixed lg:static inset-y-0 left-0 w-64 bg-slate-900 border-r border-slate-800 flex flex-col z-40
          transition-transform duration-300 ease-in-out
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        `}>
          <div className="p-6 flex justify-between items-center">
            <h1 className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-emerald-400 flex items-center gap-2">
              ⚽ FutGesti&oacute;n
            </h1>
            <button className="lg:hidden text-slate-400 hover:text-white" onClick={closeSidebar}>
              <FaTimes className="text-xl" />
            </button>
          </div>
          <nav className="flex-1 px-4 space-y-2">
            <NavItem to="/" icon={FaMoneyBillWave} onClick={closeSidebar}>Control de Pagos</NavItem>
            <NavItem to="/jugadores" icon={FaUsers} onClick={closeSidebar}>Jugadores</NavItem>
            <NavItem to="/conceptos" icon={FaList} onClick={closeSidebar}>Conceptos de Pago</NavItem>
          </nav>
          <div className="p-4 text-xs text-slate-500 text-center">
            v1.0.0
          </div>
        </aside>

        {/* Main Content */}
        <div className="flex-1 flex flex-col min-w-0 h-screen overflow-hidden">
          {/* Mobile Header */}
          <header className="lg:hidden bg-slate-900 border-b border-slate-800 p-4 flex items-center justify-between z-20">
            <h1 className="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-emerald-400">
              ⚽ FutGesti&oacute;n
            </h1>
            <button 
              className="p-2 text-slate-400 hover:text-white"
              onClick={() => setIsSidebarOpen(true)}
            >
              <FaBars className="text-xl" />
            </button>
          </header>

          <main className="flex-1 overflow-x-hidden overflow-y-auto">
            <Routes>
              <Route path="/" element={<PaymentsPage />} />
              <Route path="/jugadores" element={<PlayersPage />} />
              <Route path="/conceptos" element={<ConceptsPage />} />
            </Routes>
          </main>
        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
