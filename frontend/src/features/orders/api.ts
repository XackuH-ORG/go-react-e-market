import { api } from '../../shared/api/axios';
import { CartItem } from '../cart/api';

export interface Order {
  id: string;
  userId: string;
  total: number;
  status: 'Pending' | 'Processing' | 'Shipped' | 'Delivered' | 'Cancelled';
  createdAt: string;
  items: CartItem[];
}

export const ordersApi = {
  getOrders: async () => {
    const { data } = await api.get<Order[]>('/orders');
    return data;
  },
  getOrder: async (id: string) => {
    const { data } = await api.get<Order>(`/orders/${id}`);
    return data;
  }
};
