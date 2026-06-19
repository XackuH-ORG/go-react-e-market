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
      alert('Добавлено в корзину');
    }
  });

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-4xl font-extrabold tracking-tight text-white drop-shadow-md">Каталог товаров</h1>
        <input 
          type="text" 
          placeholder="Поиск товаров..." 
          className="bg-white/20 border border-white/30 rounded-xl px-4 py-3 w-full sm:w-80 focus:ring-2 focus:ring-white/50 outline-none transition-shadow text-white placeholder-white/60 shadow-inner"
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="animate-pulse bg-white/20 backdrop-blur-md border border-white/20 h-80 rounded-2xl"></div>
          ))}
        </div>
      ) : products?.length === 0 ? (
        <div className="text-center py-20 text-white/80 drop-shadow-sm text-lg">Товары не найдены.</div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products?.map(product => (
            <div key={product.id} className="group bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl p-5 shadow-xl hover:shadow-2xl transition-all flex flex-col text-white">
              <div className="aspect-square w-full mb-5 overflow-hidden rounded-xl bg-black/10 relative">
                {product.imageUrl ? (
                  <img src={product.imageUrl} alt={product.name} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-white/60 drop-shadow-sm">Нет фото</div>
                )}
              </div>
              <h3 className="font-bold text-xl drop-shadow-sm line-clamp-1">{product.name}</h3>
              <p className="text-white/80 text-sm mt-2 mb-5 line-clamp-2 flex-1 drop-shadow-sm">{product.description}</p>
              <div className="flex items-center justify-between mt-auto pt-4 border-t border-white/20">
                <span className="font-extrabold text-2xl drop-shadow-sm">${product.price.toFixed(2)}</span>
                {user && user.role !== 'ADMIN' && (
                  <button 
                    onClick={() => addToCart.mutate(product.id)}
                    disabled={addToCart.isPending}
                    className="bg-white text-indigo-600 px-5 py-2.5 rounded-xl text-sm font-bold hover:bg-white/90 disabled:opacity-50 transition-all shadow-md transform hover:-translate-y-0.5"
                  >
                    В корзину
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
