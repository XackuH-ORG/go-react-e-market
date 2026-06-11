import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { catalogApi } from '../../features/catalog/api';
import { useAuthStore } from '../../features/auth/store';
import { cartApi } from '../../features/cart/api';
import { useState } from 'react';

export function CatalogPage() {
  const [search, setSearch] = useState('');
  const { user } = useAuthStore();
  const queryClient = useQueryClient();

  const { data: products, isLoading } = useQuery({
    queryKey: ['products', search],
    queryFn: () => catalogApi.getProducts({ search }),
  });

  const addToCart = useMutation({
    mutationFn: (productId: string) => cartApi.addToCart(productId, 1),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      alert('Added to cart');
    }
  });

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">Discover Products</h1>
        <input 
          type="text" 
          placeholder="Search products..." 
          className="border border-zinc-300 rounded-lg px-4 py-2.5 w-full sm:w-80 focus:ring-2 focus:ring-zinc-900 outline-none transition-shadow"
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="animate-pulse bg-zinc-200 h-80 rounded-xl"></div>
          ))}
        </div>
      ) : products?.length === 0 ? (
        <div className="text-center py-20 text-zinc-500">No products found.</div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products?.map(product => (
            <div key={product.id} className="group border border-zinc-200 rounded-xl p-5 bg-white shadow-sm hover:shadow-md transition-all flex flex-col">
              <div className="aspect-square w-full mb-5 overflow-hidden rounded-lg bg-zinc-100 relative">
                {product.imageUrl ? (
                  <img src={product.imageUrl} alt={product.name} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-zinc-400">No Image</div>
                )}
              </div>
              <h3 className="font-semibold text-lg text-zinc-900 line-clamp-1">{product.name}</h3>
              <p className="text-zinc-500 text-sm mt-1 mb-4 line-clamp-2 flex-1">{product.description}</p>
              <div className="flex items-center justify-between mt-auto pt-4 border-t border-zinc-100">
                <span className="font-bold text-xl text-zinc-900">${product.price.toFixed(2)}</span>
                {user && user.role !== 'Admin' && (
                  <button 
                    onClick={() => addToCart.mutate(product.id)}
                    disabled={addToCart.isPending}
                    className="bg-zinc-900 text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-zinc-800 disabled:opacity-50 transition-colors shadow-sm"
                  >
                    Add
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
