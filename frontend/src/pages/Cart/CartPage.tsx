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
    <div className="max-w-4xl mx-auto space-y-8">
      <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">Shopping Cart</h1>
      
      {isLoading ? (
        <div className="animate-pulse space-y-4">
          <div className="h-24 bg-zinc-200 rounded-lg"></div>
          <div className="h-24 bg-zinc-200 rounded-lg"></div>
        </div>
      ) : !cartItems || cartItems.length === 0 ? (
        <div className="text-center py-20 bg-white border border-zinc-200 rounded-xl">
          <h2 className="text-xl font-medium text-zinc-900 mb-2">Your cart is empty</h2>
          <p className="text-zinc-500 mb-6">Looks like you haven't added anything yet.</p>
          <Link to="/" className="bg-zinc-900 text-white px-6 py-2.5 rounded-lg font-medium hover:bg-zinc-800 transition-colors">
            Start Shopping
          </Link>
        </div>
      ) : (
        <div className="space-y-6">
          <div className="bg-white border border-zinc-200 rounded-xl overflow-hidden shadow-sm">
            <ul className="divide-y divide-zinc-100">
              {cartItems?.map(item => (
                <li key={item.id} className="p-5 flex items-center justify-between hover:bg-zinc-50 transition-colors">
                  <div className="flex items-center space-x-4">
                    {item.product?.imageUrl ? (
                      <img src={item.product.imageUrl} alt={item.product.name} className="w-20 h-20 rounded-md object-cover border border-zinc-100" />
                    ) : (
                      <div className="w-20 h-20 bg-zinc-100 rounded-md flex items-center justify-center text-xs text-zinc-400">No Img</div>
                    )}
                    <div>
                      <h3 className="font-semibold text-lg text-zinc-900">{item.product?.name}</h3>
                      <p className="text-sm text-zinc-500 mt-1">Quantity: {item.quantity}</p>
                    </div>
                  </div>
                  <div className="flex flex-col items-end space-y-2">
                    <span className="font-bold text-lg text-zinc-900">${((item.product?.price || 0) * item.quantity).toFixed(2)}</span>
                    <button 
                      onClick={() => removeMutation.mutate(item.id)}
                      className="text-red-500 text-sm font-medium hover:text-red-600 transition-colors"
                    >
                      Remove
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
          
          <div className="flex flex-col sm:flex-row justify-between items-center bg-white p-6 rounded-xl border border-zinc-200 shadow-sm gap-4">
            <div className="text-xl">
              Total: <span className="font-extrabold text-zinc-900 ml-2">${total.toFixed(2)}</span>
            </div>
            <div className="flex items-center space-x-4 w-full sm:w-auto">
              <button 
                onClick={() => clearMutation.mutate()}
                className="text-zinc-600 font-medium hover:text-zinc-900 transition-colors"
              >
                Clear Cart
              </button>
              <button 
                onClick={() => checkoutMutation.mutate()}
                disabled={checkoutMutation.isPending}
                className="flex-1 sm:flex-none bg-zinc-900 text-white px-8 py-3 rounded-lg font-bold hover:bg-zinc-800 disabled:opacity-50 shadow-sm transition-colors"
              >
                {checkoutMutation.isPending ? 'Processing...' : 'Checkout'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
