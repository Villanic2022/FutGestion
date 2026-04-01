import React, { useState, useEffect } from 'react';
import { FaPlus, FaEdit, FaTrash } from 'react-icons/fa';
import * as api from '../api';

function PlayersPage() {
  const [players, setPlayers] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({ id: null, first_name: '', last_name: '', dni: '', birth_date: '' });

  const loadPlayers = () => {
    api.getPlayers().then(setPlayers).catch(console.error);
  };

  useEffect(() => {
    loadPlayers();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (formData.id) {
        await api.updatePlayer(formData.id, formData);
      } else {
        await api.createPlayer(formData);
      }
      setIsModalOpen(false);
      loadPlayers();
    } catch (error) {
      alert("Error guardando jugador. Verifique los datos.");
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm("¿Seguro que quiere eliminar este jugador? Se perderán sus registros de pago.")) {
      await api.deletePlayer(id);
      loadPlayers();
    }
  };

  const openModal = (player = null) => {
    if (player) {
      setFormData({ ...player, birth_date: player.birth_date.split('T')[0] });
    } else {
      setFormData({ id: null, first_name: '', last_name: '', dni: '', birth_date: '' });
    }
    setIsModalOpen(true);
  };

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-8">
        <h2 className="text-3xl font-bold">Jugadores</h2>
        <button 
          onClick={() => openModal()}
          className="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-lg shadow-blue-500/30"
        >
          <FaPlus /> Nuevo Jugador
        </button>
      </div>

      <div className="bg-slate-900 rounded-xl border border-slate-800 overflow-hidden shadow-xl">
        <table className="w-full text-left">
          <thead className="bg-slate-800/50">
            <tr>
              <th className="p-4 font-semibold text-slate-300">Nombre</th>
              <th className="p-4 font-semibold text-slate-300">Apellido</th>
              <th className="p-4 font-semibold text-slate-300">DNI</th>
              <th className="p-4 font-semibold text-slate-300 text-right">Acciones</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/50">
            {players.map(p => (
              <tr key={p.id} className="hover:bg-slate-800/30 transition-colors">
                <td className="p-4">{p.first_name}</td>
                <td className="p-4 font-medium">{p.last_name}</td>
                <td className="p-4 text-slate-400">{p.dni}</td>
                <td className="p-4 flex gap-3 justify-end">
                  <button onClick={() => openModal(p)} className="p-2 text-blue-400 hover:bg-blue-400/10 rounded transition-colors"><FaEdit /></button>
                  <button onClick={() => handleDelete(p.id)} className="p-2 text-rose-400 hover:bg-rose-400/10 rounded transition-colors"><FaTrash /></button>
                </td>
              </tr>
            ))}
            {players.length === 0 && (
              <tr>
                <td colSpan="4" className="p-8 text-center text-slate-500 italic">No hay jugadores registrados.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-xl w-full max-w-md shadow-2xl overflow-hidden">
            <div className="p-6 border-b border-slate-800">
              <h3 className="text-xl font-bold">{formData.id ? 'Editar Jugador' : 'Crear Jugador'}</h3>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-slate-400 mb-1">Nombre</label>
                  <input required type="text" value={formData.first_name} onChange={e => setFormData({...formData, first_name: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
                </div>
                <div>
                  <label className="block text-sm text-slate-400 mb-1">Apellido</label>
                  <input required type="text" value={formData.last_name} onChange={e => setFormData({...formData, last_name: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
                </div>
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1">DNI</label>
                <input required type="text" value={formData.dni} onChange={e => setFormData({...formData, dni: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1">Fecha de Nacimiento</label>
                <input required type="date" value={formData.birth_date} onChange={e => setFormData({...formData, birth_date: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 placeholder-slate-500" />
              </div>
              
              <div className="pt-4 flex justify-end gap-3">
                <button type="button" onClick={() => setIsModalOpen(false)} className="px-4 py-2 rounded-lg text-slate-400 hover:bg-slate-800 transition-colors">Cancelar</button>
                <button type="submit" className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-500 text-white font-medium transition-colors shadow-lg shadow-blue-500/20">Guardar</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default PlayersPage;
