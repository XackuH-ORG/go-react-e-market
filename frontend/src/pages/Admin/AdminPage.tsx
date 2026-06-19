import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminApi } from '../../features/admin/api';
import { catalogApi } from '../../features/catalog/api';

export function AdminPage() {
  const [tab, setTab] = useState<'products' | 'orders' | 'users'>('products');
  
  const tabTranslations: Record<string, string> = {
    products: 'Товары',
    orders: 'Заказы',
    users: 'Пользователи'
  };

  return (
    <div className="space-y-8 max-w-6xl mx-auto text-white">
      <h1 className="text-4xl font-extrabold tracking-tight drop-shadow-md">Панель администратора</h1>
      
      <div className="flex space-x-6 border-b border-white/30">
        {(['products', 'orders', 'users'] as const).map(t => (
          <button 
            key={t}
            onClick={() => setTab(t)}
            className={`py-3 px-1 font-bold text-sm uppercase tracking-wider border-b-4 transition-all ${
              tab === t 
                ? 'border-white text-white drop-shadow-sm' 
                : 'border-transparent text-white/60 hover:text-white/90'
            }`}
          >
            {tabTranslations[t]}
          </button>
        ))}
      </div>

      <div className="bg-white/20 backdrop-blur-lg border border-white/30 rounded-2xl shadow-2xl overflow-hidden">
        {tab === 'products' && <AdminProducts />}
        {tab === 'orders' && <AdminOrders />}
        {tab === 'users' && <AdminUsers />}
      </div>
    </div>
  );
}

function AdminProducts() {
  const { data: products } = useQuery({ queryKey: ['admin-products'], queryFn: () => catalogApi.getProducts() });
  
  return (
    <div className="p-8 space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold drop-shadow-sm">Товары</h2>
        <button className="bg-white text-indigo-600 px-5 py-2.5 rounded-xl text-sm font-bold hover:bg-white/90 transition-all shadow-md transform hover:-translate-y-0.5">Добавить товар</button>
      </div>
      <div className="divide-y divide-white/20">
        {products?.map(p => (
          <div key={p.id} className="py-5 flex justify-between items-center hover:bg-white/10 px-4 rounded-xl transition-colors -mx-4">
            <div className="flex items-center space-x-5">
              <div className="w-14 h-14 bg-black/10 rounded-lg border border-white/20 shadow-inner flex items-center justify-center overflow-hidden">
                {p.imageUrl ? <img src={p.imageUrl} alt="" className="w-full h-full object-cover" /> : <span className="text-[10px] text-white/50">Нет фото</span>}
              </div>
              <div>
                <div className="font-bold text-lg drop-shadow-sm">{p.name}</div>
                <div className="text-sm text-white/80">${p.price} • В наличии: <span className="font-semibold text-white">{p.stock}</span></div>
              </div>
            </div>
            <div className="space-x-4">
              <button className="text-blue-200 font-bold hover:text-white transition-colors drop-shadow-sm">Изменить</button>
              <button className="text-red-300 font-bold hover:text-red-100 transition-colors drop-shadow-sm">Удалить</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function AdminOrders() {
  const [search, setSearch] = useState('');
  const { data: orders } = useQuery({ queryKey: ['admin-orders', search], queryFn: () => adminApi.getOrders(search) });
  
  return (
    <div className="p-8 space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h2 className="text-2xl font-bold drop-shadow-sm">Заказы</h2>
        <input 
          type="text" 
          placeholder="Поиск по 4 символам..." 
          className="bg-white/20 border border-white/30 rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-white/50 outline-none transition-shadow text-white placeholder-white/60 shadow-inner w-full sm:w-auto"
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
      </div>
      <div className="divide-y divide-white/20">
        {orders?.map(o => (
          <div key={o.id} className="py-5 flex justify-between items-center hover:bg-white/10 px-4 rounded-xl transition-colors -mx-4">
            <div>
              <div className="font-mono text-sm drop-shadow-sm">{o.id}</div>
              <div className="text-sm text-white/80">${o.total} • {new Date(o.createdAt).toLocaleDateString('ru-RU')}</div>
            </div>
            <div className="flex items-center space-x-3">
              <select className="bg-white/20 border border-white/30 rounded-lg text-sm px-2 py-1.5 focus:ring-2 focus:ring-white/50 outline-none shadow-inner cursor-pointer">
                <option value={o.status} className="text-zinc-900">{o.status}</option>
                <option value="Processing" className="text-zinc-900">В обработке</option>
                <option value="Shipped" className="text-zinc-900">Отправлен</option>
                <option value="Delivered" className="text-zinc-900">Доставлен</option>
                <option value="Cancelled" className="text-zinc-900">Отменен</option>
              </select>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function AdminUsers() {
  const { data: users } = useQuery({ queryKey: ['admin-users'], queryFn: () => adminApi.getUsers() });
  
  return (
    <div className="p-8 space-y-6">
      <h2 className="text-2xl font-bold drop-shadow-sm">Пользователи</h2>
      <div className="divide-y divide-white/20">
        {users?.map(u => (
          <div key={u.id} className="py-5 flex justify-between items-center hover:bg-white/10 px-4 rounded-xl transition-colors -mx-4">
            <div>
              <div className="font-bold text-lg drop-shadow-sm">{u.email}</div>
              <div className="text-sm font-mono text-white/60">{u.id}</div>
            </div>
            <div>
              <select className="bg-white/20 border border-white/30 rounded-lg text-sm px-2 py-1.5 focus:ring-2 focus:ring-white/50 outline-none shadow-inner cursor-pointer">
                <option value={u.role} className="text-zinc-900">{u.role === 'ADMIN' ? 'Админ' : 'Покупатель'}</option>
                <option value="CUSTOMER" className="text-zinc-900">Покупатель</option>
                <option value="ADMIN" className="text-zinc-900">Админ</option>
              </select>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
