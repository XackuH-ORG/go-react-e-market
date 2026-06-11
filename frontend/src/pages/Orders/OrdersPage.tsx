import { useQuery } from '@tanstack/react-query';
import { ordersApi } from '../../features/orders/api';

export function OrdersPage() {
  const { data: orders, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: ordersApi.getOrders,
  });

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">Your Orders</h1>

      {isLoading ? (
        <div className="animate-pulse space-y-4">
          <div className="h-32 bg-zinc-200 rounded-xl"></div>
          <div className="h-32 bg-zinc-200 rounded-xl"></div>
        </div>
      ) : !orders || orders.length === 0 ? (
        <div className="text-center py-20 bg-white border border-zinc-200 rounded-xl">
          <h2 className="text-xl font-medium text-zinc-900 mb-2">No orders found</h2>
          <p className="text-zinc-500">You haven't placed any orders yet.</p>
        </div>
      ) : (
        <div className="space-y-6">
          {orders.map(order => (
            <div key={order.id} className="bg-white border border-zinc-200 rounded-xl overflow-hidden shadow-sm">
              <div className="bg-zinc-50 p-5 border-b border-zinc-200 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <div className="text-xs uppercase tracking-wider text-zinc-500 font-medium mb-1">Order ID</div>
                  <div className="font-mono text-sm text-zinc-900">{order.id}</div>
                </div>
                <div>
                  <div className="text-xs uppercase tracking-wider text-zinc-500 font-medium mb-1">Date</div>
                  <div className="text-sm text-zinc-900">{new Date(order.createdAt).toLocaleDateString()}</div>
                </div>
                <div>
                  <div className="text-xs uppercase tracking-wider text-zinc-500 font-medium mb-1">Total</div>
                  <div className="font-bold text-zinc-900">${order.total.toFixed(2)}</div>
                </div>
                <div>
                  <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-bold capitalize tracking-wide
                    ${order.status === 'Delivered' ? 'bg-green-100 text-green-800' : 
                      order.status === 'Cancelled' ? 'bg-red-100 text-red-800' : 
                      'bg-blue-100 text-blue-800'}`}>
                    {order.status}
                  </span>
                </div>
              </div>
              <ul className="divide-y divide-zinc-100 p-2">
                {order.items?.map(item => (
                  <li key={item.id} className="p-3 flex items-center justify-between hover:bg-zinc-50 rounded-lg transition-colors">
                    <div className="flex items-center space-x-4">
                      {item.product?.imageUrl ? (
                        <img src={item.product.imageUrl} alt={item.product.name} className="w-12 h-12 rounded-md object-cover border border-zinc-100" />
                      ) : (
                        <div className="w-12 h-12 bg-zinc-100 rounded-md flex items-center justify-center text-[10px] text-zinc-400">No Img</div>
                      )}
                      <div>
                        <div className="font-medium text-zinc-900">{item.product?.name}</div>
                        <div className="text-sm text-zinc-500">Qty: {item.quantity}</div>
                      </div>
                    </div>
                    <div className="font-semibold text-zinc-900">
                      ${((item.product?.price || 0) * item.quantity).toFixed(2)}
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
