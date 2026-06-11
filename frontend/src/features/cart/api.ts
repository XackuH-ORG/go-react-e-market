import { api } from '../../shared/api/axios';
import { Product } from '../catalog/api';

export interface CartItem {
  id: string;
  productId: string;
  quantity: number;
  product?: Product;
}

export const cartApi = {
  getCart: async () => {
    const { data } = await api.get<CartItem[]>('/cart');
    return data;
  },
  addToCart: async (productId: string, quantity: number) => {
    const { data } = await api.post<CartItem>('/cart', { productId, quantity });
    return data;
  },
  updateQuantity: async (itemId: string, quantity: number) => {
    const { data } = await api.put<CartItem>(`/cart/${itemId}`, { quantity });
    return data;
  },
  removeFromCart: async (itemId: string) => {
    await api.delete(`/cart/${itemId}`);
  },
  clearCart: async () => {
    await api.delete('/cart/clear');
  },
  checkout: async () => {
    const { data } = await api.post<{ orderId: string }>('/cart/checkout');
    return data;
  }
};
