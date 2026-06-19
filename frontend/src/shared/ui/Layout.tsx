import { Outlet, Link } from 'react-router-dom';
import { useAuthStore } from '../../features/auth/store';

export function Layout() {
  const { user, logout } = useAuthStore();

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 text-white">
      <header className="bg-white/20 backdrop-blur-md border-b border-white/30 sticky top-0 z-50">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="text-xl font-bold tracking-tight text-white drop-shadow-md">
            E-Market
          </Link>
          <nav className="flex items-center gap-4">
            {user ? (
              <>
                <span className="text-sm text-white/80">{user.email} ({user.role === 'ADMIN' ? 'Админ' : 'Пользователь'})</span>
                {user.role === 'ADMIN' && (
                  <Link to="/admin" className="text-sm font-medium hover:text-white/70 transition-colors">Админ</Link>
                )}
                <Link to="/cart" className="text-sm font-medium hover:text-white/70 transition-colors">Корзина</Link>
                <Link to="/orders" className="text-sm font-medium hover:text-white/70 transition-colors">Заказы</Link>
                <button 
                  onClick={logout}
                  className="text-sm font-medium text-red-200 hover:text-red-100 transition-colors"
                >
                  Выйти
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm font-medium hover:text-white/70 transition-colors">Войти</Link>
                <Link to="/register" className="text-sm font-medium hover:text-white/70 transition-colors">Регистрация</Link>
              </>
            )}
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4 py-8 relative z-10">
        <Outlet />
      </main>
    </div>
  );
}
