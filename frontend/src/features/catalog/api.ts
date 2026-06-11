import { api } from '../../shared/api/axios';

export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  imageUrl?: string;
  categoryId?: string;
}

export const catalogApi = {
  getProducts: async (params?: { search?: string; sortBy?: string; order?: string; categoryId?: string }) => {
    const { data } = await api.get<Product[]>('/products', { params });
    return data;
  },
  getProduct: async (id: string) => {
    const { data } = await api.get<Product>(`/products/${id}`);
    return data;
  }
};
