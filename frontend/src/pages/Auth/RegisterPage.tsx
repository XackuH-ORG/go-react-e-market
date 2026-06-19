import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../../features/auth/api';
import { useAuthStore } from '../../features/auth/store';

export function RegisterPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();
  const { setUser } = useAuthStore();

  const mutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setUser(data.user);
      navigate('/');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate({ email, password });
  };

  return (
    <div className="max-w-md mx-auto mt-10 p-8 bg-white/20 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/30 text-white">
      <h2 className="text-3xl font-bold mb-6 text-center drop-shadow-md">Регистрация</h2>
      {mutation.isError && <div className="text-red-200 bg-red-500/20 px-3 py-2 rounded-lg mb-4 text-sm text-center">Ошибка регистрации. Email может быть занят.</div>}
      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label className="block text-sm font-medium mb-1 drop-shadow-sm">Email</label>
          <input 
            type="email" 
            className="w-full bg-white/20 border border-white/30 rounded-xl px-4 py-3 focus:ring-2 focus:ring-white/50 outline-none text-white placeholder-white/60 transition-all shadow-inner"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1 drop-shadow-sm">Пароль</label>
          <input 
            type="password" 
            className="w-full bg-white/20 border border-white/30 rounded-xl px-4 py-3 focus:ring-2 focus:ring-white/50 outline-none text-white placeholder-white/60 transition-all shadow-inner"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
        </div>
        <button 
          type="submit" 
          disabled={mutation.isPending}
          className="w-full bg-white text-indigo-600 rounded-xl py-3 font-bold hover:bg-white/90 disabled:opacity-50 transition-all shadow-md transform hover:-translate-y-0.5"
        >
          {mutation.isPending ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
      </form>
      <div className="mt-6 text-center text-sm text-white/80">
        Уже есть аккаунт? <Link to="/login" className="text-white font-bold hover:underline drop-shadow-md">Войти</Link>
      </div>
    </div>
  );
}
