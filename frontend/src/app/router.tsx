import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom';
import { Layout } from '../shared/ui/Layout';
import { CatalogPage } from '../pages/Catalog/CatalogPage';
import { LoginPage } from '../pages/Auth/LoginPage';
import { RegisterPage } from '../pages/Auth/RegisterPage';
import { CartPage } from '../pages/Cart/CartPage';
import { OrdersPage } from '../pages/Orders/OrdersPage';
import { AdminPage } from '../pages/Admin/AdminPage';
import { useAuthStore } from '../features/auth/store';

function ProtectedRoute({ children, role }: { children: React.ReactNode, role?: 'ADMIN' }) {
  const { user, isAuthenticated } = useAuthStore();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (role && user?.role !== role) return <Navigate to="/" replace />;
  return <>{children}</>;
}

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <CatalogPage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'register', element: <RegisterPage /> },
      { 
        path: 'cart', 
        element: <ProtectedRoute><CartPage /></ProtectedRoute> 
      },
      { 
        path: 'orders', 
        element: <ProtectedRoute><OrdersPage /></ProtectedRoute> 
      },
      { 
        path: 'admin', 
        element: <ProtectedRoute role="ADMIN"><AdminPage /></ProtectedRoute> 
      },
    ],
  },
]);

export function AppRouter() {
  return <RouterProvider router={router} />;
}
