import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { Layout } from './Layout';
import { useAuthStore } from '../../features/auth/store';

describe('Layout Component', () => {
  beforeEach(() => {
    // Reset auth store before each test
    useAuthStore.setState({ user: null, isAuthenticated: false });
  });

  const renderLayout = () => {
    return render(
      <MemoryRouter>
        <Layout />
      </MemoryRouter>
    );
  };

  it('renders standard links when not authenticated', () => {
    renderLayout();
    
    expect(screen.getByText('E-Market')).toBeInTheDocument();
    expect(screen.getByText('Войти')).toBeInTheDocument();
    expect(screen.getByText('Регистрация')).toBeInTheDocument();
    expect(screen.queryByText('Выйти')).not.toBeInTheDocument();
    expect(screen.queryByText('Корзина')).not.toBeInTheDocument();
    expect(screen.queryByText('Logout')).not.toBeInTheDocument();
  });

  it('renders user links when authenticated as User', () => {
    useAuthStore.setState({ 
      user: { id: '1', email: 'user@test.com', role: 'CUSTOMER' },
      isAuthenticated: true 
    });

    renderLayout();

    expect(screen.getByText('user@test.com (Пользователь)')).toBeInTheDocument();
    expect(screen.getByText('Корзина')).toBeInTheDocument();
    expect(screen.getByText('Заказы')).toBeInTheDocument();
    expect(screen.getByText('Выйти')).toBeInTheDocument();
    
    expect(screen.queryByText('Админ')).not.toBeInTheDocument();
    expect(screen.queryByText('Войти')).not.toBeInTheDocument();
    expect(screen.queryByText('Регистрация')).not.toBeInTheDocument();
  });

  it('renders Admin link when authenticated as Admin', () => {
    useAuthStore.setState({ 
      user: { id: '2', email: 'admin@test.com', role: 'ADMIN' },
      isAuthenticated: true 
    });

    renderLayout();

    expect(screen.getByText('admin@test.com (Админ)')).toBeInTheDocument();
    expect(screen.getByText('Админ')).toBeInTheDocument();
    expect(screen.getByText('Корзина')).toBeInTheDocument();
    expect(screen.getByText('Заказы')).toBeInTheDocument();
    expect(screen.getByText('Выйти')).toBeInTheDocument();
  });

  it('calls logout when Logout button is clicked', () => {
    useAuthStore.setState({ 
      user: { id: '1', email: 'user@test.com', role: 'CUSTOMER' },
      isAuthenticated: true 
    });

    renderLayout();

    const logoutBtn = screen.getByText('Выйти');
    fireEvent.click(logoutBtn);

    // After clicking logout, store should be updated and layout should rerender
    // to show Login/Register links again
    expect(useAuthStore.getState().user).toBeNull();
  });
});
