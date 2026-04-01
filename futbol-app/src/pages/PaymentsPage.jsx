import React, { useState, useEffect } from 'react';
import { FaFilePdf, FaCheck, FaTimes, FaSpinner, FaSearch } from 'react-icons/fa';
import * as api from '../api';

function PaymentsPage() {
  const [data, setData] = useState({ concepts: [], rows: [] });
  const [isLoading, setIsLoading] = useState(true);
  const [filterConcept, setFilterConcept] = useState('');

  const loadMatrix = () => {
    setIsLoading(true);
    api.getMatrix()
      .then(res => setData(res))
      .catch(console.error)
      .finally(() => setIsLoading(false));
  };

  useEffect(() => {
    loadMatrix();
  }, []);

  const togglePayment = async (playerId, conceptId, currentStatus) => {
    try {
      // Optimistic update
      const newStatus = !currentStatus;
      const newData = { ...data };
      const row = newData.rows.find(r => r.player.id === playerId);
      if (row) {
        if (!row.payments[conceptId]) {
          row.payments[conceptId] = { paid: newStatus };
        } else {
          row.payments[conceptId].paid = newStatus;
        }
        setData(newData);
      }

      // API call
      await api.updatePayment({
        player_id: playerId,
        concept_id: conceptId,
        paid: newStatus,
        amount: 0,
        notes: ''
      });
    } catch (error) {
      console.error(error);
      loadMatrix(); // Revert on error
    }
  };

  const handleExport = () => {
    window.open(api.exportPDFUrl(filterConcept), '_blank');
  };

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex flex-col items-center gap-4 text-slate-400">
          <FaSpinner className="animate-spin text-4xl text-blue-500" />
          <p>Cargando matriz de pagos...</p>
        </div>
      </div>
    );
  }

  // Filter concepts to show
  const displayConcepts = filterConcept 
    ? data.concepts.filter(c => c.id.toString() === filterConcept)
    : data.concepts;

  return (
    <div className="p-8 h-full flex flex-col">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <h2 className="text-3xl font-bold mb-1">Control de Pagos</h2>
          <p className="text-slate-400 text-sm">Hacé click en las celdas para alternar entre Pagado / No Pagado</p>
        </div>
        
        <div className="flex gap-4 items-center">
          <div className="relative">
            <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
            <select 
              className="bg-slate-900 border border-slate-700 outline-none rounded-lg pl-10 pr-4 py-2 appearance-none cursor-pointer hover:border-slate-500 transition-colors"
              value={filterConcept}
              onChange={(e) => setFilterConcept(e.target.value)}
            >
              <option value="">Todos los conceptos</option>
              {data.concepts.map(c => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
          
          <button 
            onClick={handleExport}
            className="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-lg shadow-indigo-500/30 font-medium whitespace-nowrap"
          >
            <FaFilePdf /> Exportar PDF
          </button>
        </div>
      </div>

      <div className="flex-1 bg-slate-900 rounded-xl border border-slate-800 overflow-hidden shadow-2xl relative">
        {data.concepts.length === 0 || data.rows.length === 0 ? (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-slate-500">
            <p>Faltan datos para mostrar la matriz.</p>
            <p className="text-sm mt-2">Asegurate de tener al menos un jugador y un concepto creados.</p>
          </div>
        ) : (
          <div className="absolute inset-0 overflow-auto">
            <table className="w-full text-left border-collapse min-w-max">
              <thead className="bg-slate-800/80 backdrop-blur-md sticky top-0 z-10 shadow-sm">
                <tr>
                  <th className="p-4 font-semibold text-slate-200 border-b border-r border-slate-700/50 sticky left-0 bg-slate-800 z-20 shadow-[2px_0_5px_-2px_rgba(0,0,0,0.5)] w-64">
                    Jugador
                  </th>
                  {displayConcepts.map(c => (
                    <th key={c.id} className="p-4 font-semibold text-slate-300 border-b border-r border-slate-700/50 text-center w-32 min-w-[120px] max-w-[200px] truncate" title={c.name}>
                      {c.name}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {data.rows.map((row, i) => (
                  <tr key={row.player.id} className="hover:bg-slate-800/30 transition-colors group">
                    <td className="p-4 border-b border-r border-slate-800/50 sticky left-0 font-medium text-slate-300 bg-slate-900 group-hover:bg-slate-800/80 transition-colors z-10 shadow-[2px_0_5px_-2px_rgba(0,0,0,0.5)] whitespace-nowrap">
                      {row.player.last_name}, {row.player.first_name}
                    </td>
                    {displayConcepts.map(c => {
                      const payment = row.payments[c.id];
                      const isPaid = payment?.paid || false;
                      return (
                        <td 
                          key={c.id} 
                          className="p-1 border-b border-r border-slate-800/50 text-center cursor-pointer transition-colors"
                          onClick={() => togglePayment(row.player.id, c.id, isPaid)}
                        >
                          <div className={`mx-auto w-10 h-10 rounded-lg flex items-center justify-center transition-all duration-200 ${
                            isPaid 
                              ? 'bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30' 
                              : 'bg-rose-500/10 text-rose-500 hover:bg-rose-500/20'
                          }`}>
                            {isPaid ? <FaCheck /> : <FaTimes className="text-rose-500/80" />}
                          </div>
                        </td>
                      );
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default PaymentsPage;
