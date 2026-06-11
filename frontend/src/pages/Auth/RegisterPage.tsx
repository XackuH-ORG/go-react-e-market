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
    <div className="max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-sm border border-zinc-200">
      <h2 className="text-2xl font-bold mb-6 text-center">Register</h2>
      {mutation.isError && <div className="text-red-500 mb-4 text-sm">Registration failed. Email might be taken.</div>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input 
            type="email" 
            className="w-full border border-zinc-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-zinc-900 outline-none"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Password</label>
          <input 
            type="password" 
            className="w-full border border-zinc-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-zinc-900 outline-none"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
        </div>
        <button 
          type="submit" 
          disabled={mutation.isPending}
          className="w-full bg-zinc-900 text-white rounded-md py-2 font-medium hover:bg-zinc-800 disabled:opacity-50 transition-colors"
        >
          {mutation.isPending ? 'Registering...' : 'Register'}
        </button>
      </form>
      <div className="mt-4 text-center text-sm text-zinc-600">
        Already have an account? <Link to="/login" className="text-zinc-900 font-medium hover:underline">Login</Link>
      </div>
    </div>
  );
}
