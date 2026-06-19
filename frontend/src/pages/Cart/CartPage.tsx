import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { cartApi } from '../../features/cart/api';
import { useNavigate, Link } from 'react-router-dom';

export function CartPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { data: cartItems, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: cartApi.getCart,
  });

  const removeMutation = useMutation({
    mutationFn: cartApi.removeFromCart,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['cart'] })
  });

  const clearMutation = useMutation({
    mutationFn: cartApi.clearCart,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['cart'] })
  });

  const checkoutMutation = useMutation({
    mutationFn: cartApi.checkout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      navigate('/orders');
    }
  });

  const total = cartItems?.reduce((acc, item) => acc + (item.product?.price || 0) * item.quantity, 0) || 0;

  return (
    <div className="max-w-4xl mx-auto space-y-8 text-white">
      <h1 className="text-4xl font-extrabold tracking-tight drop-shadow-md">Корзина</h1>
      
      {isLoading ? (
        <div className="animate-pulse space-y-4">
          <div className="h-24 bg-white/20 backdrop-blur-md rounded-2xl border border-white/20"></div>
          <div className="h-24 bg-white/20 backdrop-blur-md rounded-2xl border border-white/20"></div>
        </div>
      ) : !cartItems || cartItems.length === 0 ? (
        <div className="text-center py-20 bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl shadow-xl">
          <h2 className="text-2xl font-bold mb-2 drop-shadow-sm">Ваша корзина пуста</h2>
          <p className="text-white/80 mb-6 drop-shadow-sm">Похоже, вы еще ничего не добавили.</p>
          <Link to="/" className="inline-block bg-white text-indigo-600 px-8 py-3 rounded-xl font-bold hover:bg-white/90 transition-all shadow-md transform hover:-translate-y-0.5">
            Начать покупки
          </Link>
        </div>
      ) : (
        <div className="space-y-6">
          <div className="bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl overflow-hidden shadow-xl">
            <ul className="divide-y divide-white/20">
              {cartItems?.map(item => (
                <li key={item.id} className="p-5 flex items-center justify-between hover:bg-white/10 transition-colors">
                  <div className="flex items-center space-x-5">
                    {item.product?.imageUrl ? (
                      <img src={item.product.imageUrl} alt={item.product.name} className="w-24 h-24 rounded-xl object-cover border border-white/20 shadow-sm" />
                    ) : (
                      <div className="w-24 h-24 bg-black/10 rounded-xl flex items-center justify-center text-sm text-white/60 shadow-inner">Нет фото</div>
                    )}
                    <div>
                      <h3 className="font-bold text-xl drop-shadow-sm">{item.product?.name}</h3>
                      <p className="text-sm text-white/80 mt-1">Количество: <span className="font-semibold text-white">{item.quantity}</span></p>
                    </div>
                  </div>
                  <div className="flex flex-col items-end space-y-3">
                    <span className="font-extrabold text-2xl drop-shadow-sm">${((item.product?.price || 0) * item.quantity).toFixed(2)}</span>
                    <button 
                      onClick={() => removeMutation.mutate(item.id)}
                      className="text-red-300 text-sm font-bold hover:text-red-200 transition-colors drop-shadow-sm"
                    >
                      Удалить
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
          
          <div className="flex flex-col sm:flex-row justify-between items-center bg-white/20 backdrop-blur-lg p-6 rounded-2xl border border-white/30 shadow-xl gap-6">
            <div className="text-2xl drop-shadow-sm">
              Итого: <span className="font-extrabold ml-2">${total.toFixed(2)}</span>
            </div>
            <div className="flex items-center space-x-6 w-full sm:w-auto">
              <button 
                onClick={() => clearMutation.mutate()}
                className="text-white/80 font-bold hover:text-white transition-colors"
              >
                Очистить корзину
              </button>
              <button 
                onClick={() => checkoutMutation.mutate()}
                disabled={checkoutMutation.isPending}
                className="flex-1 sm:flex-none bg-white text-indigo-600 px-8 py-3 rounded-xl font-bold hover:bg-white/90 disabled:opacity-50 transition-all shadow-md transform hover:-translate-y-0.5"
              >
                {checkoutMutation.isPending ? 'Обработка...' : 'Оформить заказ'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
