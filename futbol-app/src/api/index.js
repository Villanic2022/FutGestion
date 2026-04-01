import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
});

export const getPlayers = () => api.get('/players').then(res => res.data);
export const createPlayer = (data) => api.post('/players', data).then(res => res.data);
export const updatePlayer = (id, data) => api.put(`/players/${id}`, data).then(res => res.data);
export const deletePlayer = (id) => api.delete(`/players/${id}`);

export const getConcepts = () => api.get('/concepts').then(res => res.data);
export const createConcept = (data) => api.post('/concepts', data).then(res => res.data);
export const updateConcept = (id, data) => api.put(`/concepts/${id}`, data).then(res => res.data);
export const deleteConcept = (id) => api.delete(`/concepts/${id}`);

export const getMatrix = () => api.get('/payments/matrix').then(res => res.data);
export const updatePayment = (data) => api.put('/payments', data).then(res => res.data);
export const exportPDFUrl = (conceptId = '') => `${API_URL}/export/pdf${conceptId ? `?concept_id=${conceptId}` : ''}`;
