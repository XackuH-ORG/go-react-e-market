import { api } from '../../shared/api/axios';
import { Product } from '../catalog/api';
import { Order } from '../orders/api';
import { User } from '../auth/store';

export const adminApi = {
  // Products
  createProduct: async (product: Partial<Product> | FormData) => {
    const { data } = await api.post<Product>('/admin/products', product);
    return data;
  },
  updateProduct: async (id: string, product: Partial<Product> | FormData) => {
    const { data } = await api.put<Product>(`/admin/products/${id}`, product);
    return data;
  },
  deleteProduct: async (id: string) => {
    await api.delete(`/admin/products/${id}`);
  },

  // Orders
  getOrders: async (searchQuery?: string) => {
    const { data } = await api.get<Order[]>('/admin/orders', { params: { search: searchQuery } });
    return data;
  },
  updateOrderStatus: async (id: string, status: string) => {
    const { data } = await api.patch<Order>(`/admin/orders/${id}/status`, { status });
    return data;
  },

  // Users
  getUsers: async () => {
    const { data } = await api.get<User[]>('/admin/users');
    return data;
  },
  updateUserRole: async (id: string, role: string) => {
    const { data } = await api.patch<User>(`/admin/users/${id}/role`, { role });
    return data;
  }
};
