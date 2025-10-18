import axios from 'axios';
import { Trade, CreateTradeRequest, UpdateTradeRequest } from '../types/trade';

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const tradesApi = {
  // Get all trades
  getAll: async (): Promise<Trade[]> => {
    const response = await api.get<Trade[]>('/trades');
    return response.data || [];
  },

  // Get a single trade by ID
  getById: async (id: number): Promise<Trade> => {
    const response = await api.get<Trade>(`/trades/${id}`);
    return response.data;
  },

  // Create a new trade
  create: async (data: CreateTradeRequest): Promise<Trade> => {
    const response = await api.post<Trade>('/trades', data);
    return response.data;
  },

  // Update an existing trade
  update: async (id: number, data: UpdateTradeRequest): Promise<Trade> => {
    const response = await api.put<Trade>(`/trades/${id}`, data);
    return response.data;
  },

  // Delete a trade
  delete: async (id: number): Promise<void> => {
    await api.delete(`/trades/${id}`);
  },
};

export default api;
