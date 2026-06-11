import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminApi } from '../../features/admin/api';
import { catalogApi } from '../../features/catalog/api';

export function AdminPage() {
  const [tab, setTab] = useState<'products' | 'orders' | 'users'>('products');
  
  return (
    <div className="space-y-8 max-w-6xl mx-auto">
      <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">Admin Dashboard</h1>
      
      <div className="flex space-x-6 border-b border-zinc-200">
        {(['products', 'orders', 'users'] as const).map(t => (
          <button 
            key={t}
            onClick={() => setTab(t)}
            className={`py-3 px-1 font-semibold text-sm uppercase tracking-wider border-b-2 transition-colors ${
              tab === t 
                ? 'border-zinc-900 text-zinc-900' 
                : 'border-transparent text-zinc-500 hover:text-zinc-800'
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      <div className="bg-white border border-zinc-200 rounded-xl shadow-sm overflow-hidden">
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
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Products</h2>
        <button className="bg-zinc-900 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-zinc-800">Add Product</button>
      </div>
      <div className="divide-y divide-zinc-100">
        {products?.map(p => (
          <div key={p.id} className="py-4 flex justify-between items-center">
            <div className="flex items-center space-x-4">
              <div className="w-12 h-12 bg-zinc-100 rounded-md">
                {p.imageUrl && <img src={p.imageUrl} alt="" className="w-full h-full object-cover rounded-md" />}
              </div>
              <div>
                <div className="font-medium text-zinc-900">{p.name}</div>
                <div className="text-sm text-zinc-500">${p.price} • Stock: {p.stock}</div>
              </div>
            </div>
            <div className="space-x-3">
              <button className="text-blue-600 text-sm font-medium hover:underline">Edit</button>
              <button className="text-red-600 text-sm font-medium hover:underline">Delete</button>
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
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Orders</h2>
        <input 
          type="text" 
          placeholder="Search by last 4 chars..." 
          className="border border-zinc-300 rounded-lg px-3 py-1.5 text-sm"
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
      </div>
      <div className="divide-y divide-zinc-100">
        {orders?.map(o => (
          <div key={o.id} className="py-4 flex justify-between items-center">
            <div>
              <div className="font-mono text-sm">{o.id}</div>
              <div className="text-sm text-zinc-500">${o.total} • {new Date(o.createdAt).toLocaleDateString()}</div>
            </div>
            <div className="flex items-center space-x-3">
              <select className="border border-zinc-300 rounded-md text-sm p-1">
                <option value={o.status}>{o.status}</option>
                <option value="Processing">Processing</option>
                <option value="Shipped">Shipped</option>
                <option value="Delivered">Delivered</option>
                <option value="Cancelled">Cancelled</option>
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
    <div className="p-6 space-y-6">
      <h2 className="text-xl font-bold">Users</h2>
      <div className="divide-y divide-zinc-100">
        {users?.map(u => (
          <div key={u.id} className="py-4 flex justify-between items-center">
            <div>
              <div className="font-medium">{u.email}</div>
              <div className="text-sm text-zinc-500">{u.id}</div>
            </div>
            <div>
              <select className="border border-zinc-300 rounded-md text-sm p-1">
                <option value={u.role}>{u.role}</option>
                <option value="User">User</option>
                <option value="Admin">Admin</option>
              </select>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
