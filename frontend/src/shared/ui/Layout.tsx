import { Outlet, Link } from 'react-router-dom';
import { useAuthStore } from '../../features/auth/store';

export function Layout() {
  const { user, logout } = useAuthStore();

  return (
    <div className="min-h-screen flex flex-col bg-zinc-50">
      <header className="bg-white border-b border-zinc-200">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="text-xl font-bold tracking-tight text-zinc-900">
            E-Market
          </Link>
          <nav className="flex items-center gap-4">
            {user ? (
              <>
                <span className="text-sm text-zinc-600">{user.email} ({user.role})</span>
                {user.role === 'Admin' && (
                  <Link to="/admin" className="text-sm font-medium hover:text-primary">Admin</Link>
                )}
                <Link to="/cart" className="text-sm font-medium hover:text-primary">Cart</Link>
                <Link to="/orders" className="text-sm font-medium hover:text-primary">Orders</Link>
                <button 
                  onClick={logout}
                  className="text-sm font-medium text-destructive hover:text-destructive/80"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm font-medium hover:text-primary">Login</Link>
                <Link to="/register" className="text-sm font-medium hover:text-primary">Register</Link>
              </>
            )}
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4 py-8">
        <Outlet />
      </main>
    </div>
  );
}
