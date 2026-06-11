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
    expect(screen.getByText('Login')).toBeInTheDocument();
    expect(screen.getByText('Register')).toBeInTheDocument();
    
    expect(screen.queryByText('Cart')).not.toBeInTheDocument();
    expect(screen.queryByText('Orders')).not.toBeInTheDocument();
    expect(screen.queryByText('Logout')).not.toBeInTheDocument();
  });

  it('renders user links when authenticated as User', () => {
    useAuthStore.setState({ 
      user: { id: '1', email: 'user@test.com', role: 'User' },
      isAuthenticated: true 
    });

    renderLayout();

    expect(screen.getByText('user@test.com (User)')).toBeInTheDocument();
    expect(screen.getByText('Cart')).toBeInTheDocument();
    expect(screen.getByText('Orders')).toBeInTheDocument();
    expect(screen.getByText('Logout')).toBeInTheDocument();
    
    expect(screen.queryByText('Admin')).not.toBeInTheDocument();
    expect(screen.queryByText('Login')).not.toBeInTheDocument();
    expect(screen.queryByText('Register')).not.toBeInTheDocument();
  });

  it('renders Admin link when authenticated as Admin', () => {
    useAuthStore.setState({ 
      user: { id: '2', email: 'admin@test.com', role: 'Admin' },
      isAuthenticated: true 
    });

    renderLayout();

    expect(screen.getByText('admin@test.com (Admin)')).toBeInTheDocument();
    expect(screen.getByText('Admin')).toBeInTheDocument();
    expect(screen.getByText('Cart')).toBeInTheDocument();
    expect(screen.getByText('Orders')).toBeInTheDocument();
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });

  it('calls logout when Logout button is clicked', () => {
    useAuthStore.setState({ 
      user: { id: '1', email: 'user@test.com', role: 'User' },
      isAuthenticated: true 
    });

    renderLayout();

    const logoutBtn = screen.getByText('Logout');
    fireEvent.click(logoutBtn);

    // After clicking logout, store should be updated and layout should rerender
    // to show Login/Register links again
    expect(useAuthStore.getState().user).toBeNull();
  });
});
