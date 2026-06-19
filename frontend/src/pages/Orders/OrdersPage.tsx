import { useQuery } from '@tanstack/react-query';
import { ordersApi } from '../../features/orders/api';

export function OrdersPage() {
  const { data: orders, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: ordersApi.getOrders,
  });

  const translateStatus = (status: string) => {
    const statuses: Record<string, string> = {
      Pending: 'В ожидании',
      Processing: 'В обработке',
      Shipped: 'Отправлен',
      Delivered: 'Доставлен',
      Cancelled: 'Отменен'
    };
    return statuses[status] || status;
  };

  return (
    <div className="max-w-4xl mx-auto space-y-8 text-white">
      <h1 className="text-4xl font-extrabold tracking-tight drop-shadow-md">Ваши заказы</h1>

      {isLoading ? (
        <div className="animate-pulse space-y-4">
          <div className="h-32 bg-white/20 backdrop-blur-md rounded-2xl border border-white/20"></div>
          <div className="h-32 bg-white/20 backdrop-blur-md rounded-2xl border border-white/20"></div>
        </div>
      ) : !orders || orders.length === 0 ? (
        <div className="text-center py-20 bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl shadow-xl">
          <h2 className="text-2xl font-bold mb-2 drop-shadow-sm">Заказы не найдены</h2>
          <p className="text-white/80 drop-shadow-sm">Вы еще не сделали ни одного заказа.</p>
        </div>
      ) : (
        <div className="space-y-6">
          {orders.map(order => (
            <div key={order.id} className="bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl overflow-hidden shadow-xl">
              <div className="bg-white/10 p-6 border-b border-white/20 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <div className="text-xs uppercase tracking-wider text-white/60 font-bold mb-1">ID Заказа</div>
                  <div className="font-mono text-sm drop-shadow-sm">{order.id}</div>
                </div>
                <div>
                  <div className="text-xs uppercase tracking-wider text-white/60 font-bold mb-1">Дата</div>
                  <div className="text-sm drop-shadow-sm">{new Date(order.createdAt).toLocaleDateString('ru-RU')}</div>
                </div>
                <div>
                  <div className="text-xs uppercase tracking-wider text-white/60 font-bold mb-1">Сумма</div>
                  <div className="font-extrabold text-lg drop-shadow-sm">${order.total.toFixed(2)}</div>
                </div>
                <div>
                  <span className={`inline-flex items-center px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-sm border
                    ${order.status === 'Delivered' ? 'bg-green-500/20 text-green-100 border-green-500/30' : 
                      order.status === 'Cancelled' ? 'bg-red-500/20 text-red-100 border-red-500/30' : 
                      'bg-blue-500/20 text-blue-100 border-blue-500/30'}`}>
                    {translateStatus(order.status)}
                  </span>
                </div>
              </div>
              <ul className="divide-y divide-white/20 p-2">
                {order.items?.map(item => (
                  <li key={item.id} className="p-4 flex items-center justify-between hover:bg-white/10 rounded-xl transition-colors">
                    <div className="flex items-center space-x-5">
                      {item.product?.imageUrl ? (
                        <img src={item.product.imageUrl} alt={item.product.name} className="w-14 h-14 rounded-lg object-cover border border-white/20 shadow-sm" />
                      ) : (
                        <div className="w-14 h-14 bg-black/10 rounded-lg flex items-center justify-center text-[10px] text-white/60 shadow-inner">Нет фото</div>
                      )}
                      <div>
                        <div className="font-bold text-lg drop-shadow-sm">{item.product?.name}</div>
                        <div className="text-sm text-white/80">Кол-во: <span className="font-semibold text-white">{item.quantity}</span></div>
                      </div>
                    </div>
                    <div className="font-extrabold text-lg drop-shadow-sm">
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
