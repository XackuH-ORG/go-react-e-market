import { api } from '../../shared/api/axios';
import { User } from './store';

export const authApi = {
  login: async (credentials: any) => {
    const { data } = await api.post<{user: User}>('/auth/login', credentials);
    return data;
  },
  register: async (credentials: any) => {
    const { data } = await api.post<{user: User}>('/auth/register', credentials);
    return data;
  },
  me: async () => {
    const { data } = await api.get<{user: User}>('/auth/me');
    return data;
  },
  logout: async () => {
    await api.post('/auth/logout');
  }
};
