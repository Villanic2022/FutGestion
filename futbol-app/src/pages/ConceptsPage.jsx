import React, { useState, useEffect } from 'react';
import { FaPlus, FaEdit, FaTrash } from 'react-icons/fa';
import * as api from '../api';

function ConceptsPage() {
  const [concepts, setConcepts] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({ id: null, name: '', description: '' });

  const loadConcepts = () => {
    api.getConcepts().then(setConcepts).catch(console.error);
  };

  useEffect(() => {
    loadConcepts();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (formData.id) {
        await api.updateConcept(formData.id, formData);
      } else {
        await api.createConcept(formData);
      }
      setIsModalOpen(false);
      loadConcepts();
    } catch (error) {
      alert("Error guardando concepto.");
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm("¿Seguro que quiere eliminar este concepto? Se perderán todos los registros de pago asociados.")) {
      await api.deleteConcept(id);
      loadConcepts();
    }
  };

  const openModal = (concept = null) => {
    if (concept) {
      setFormData(concept);
    } else {
      setFormData({ id: null, name: '', description: '' });
    }
    setIsModalOpen(true);
  };

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-8">
        <h2 className="text-3xl font-bold">Conceptos de Pago</h2>
        <button 
          onClick={() => openModal()}
          className="bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-lg shadow-emerald-500/30"
        >
          <FaPlus /> Nuevo Concepto
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {concepts.map(c => (
          <div key={c.id} className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-lg hover:border-slate-700 transition-colors">
            <div className="flex justify-between items-start mb-4">
              <h3 className="text-xl font-bold text-slate-100">{c.name}</h3>
              <div className="flex gap-2">
                <button onClick={() => openModal(c)} className="text-slate-400 hover:text-blue-400 transition-colors"><FaEdit /></button>
                <button onClick={() => handleDelete(c.id)} className="text-slate-400 hover:text-rose-400 transition-colors"><FaTrash /></button>
              </div>
            </div>
            <p className="text-slate-400 text-sm">
              {c.description || <span className="italic text-slate-600">Sin descripción</span>}
            </p>
          </div>
        ))}
        {concepts.length === 0 && (
          <div className="col-span-full text-center p-12 bg-slate-900 border border-slate-800 rounded-xl">
            <p className="text-slate-500 italic">No hay conceptos de pago creados.</p>
            <p className="text-sm mt-2 text-slate-600">Creá conceptos como "Seguro Anuál", "Cuota Camiseta 1", etc.</p>
          </div>
        )}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-xl w-full max-w-md shadow-2xl overflow-hidden">
            <div className="p-6 border-b border-slate-800">
              <h3 className="text-xl font-bold">{formData.id ? 'Editar Concepto' : 'Crear Concepto'}</h3>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-slate-400 mb-1">Nombre</label>
                <input required type="text" placeholder="Ej: Seguro 2026" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500" />
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1">Descripción</label>
                <textarea rows="3" value={formData.description} onChange={e => setFormData({...formData, description: e.target.value})} className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 placeholder-slate-600" placeholder="Opcional"></textarea>
              </div>
              
              <div className="pt-4 flex justify-end gap-3">
                <button type="button" onClick={() => setIsModalOpen(false)} className="px-4 py-2 rounded-lg text-slate-400 hover:bg-slate-800 transition-colors">Cancelar</button>
                <button type="submit" className="px-4 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-medium transition-colors shadow-lg shadow-emerald-500/20">Guardar</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default ConceptsPage;
