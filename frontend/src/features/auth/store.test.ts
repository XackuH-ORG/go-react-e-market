import { describe, it, expect, beforeEach } from 'vitest';
import { useAuthStore, User, Role } from './store';

describe('Auth Store', () => {
  beforeEach(() => {
    // Reset store before each test
    useAuthStore.setState({ user: null, isAuthenticated: false });
  });

  it('should have initial state with null user and not authenticated', () => {
    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });

  it('should set user and update isAuthenticated to true', () => {
    const mockUser: User = { id: '1', email: 'test@example.com', role: 'CUSTOMER' };
    
    useAuthStore.getState().setUser(mockUser);
    
    const state = useAuthStore.getState();
    expect(state.user).toEqual(mockUser);
    expect(state.isAuthenticated).toBe(true);
  });

  it('should logout by setting user to null and isAuthenticated to false', () => {
    const mockUser: User = { id: '1', email: 'test@example.com', role: 'CUSTOMER' };
    useAuthStore.getState().setUser(mockUser);
    
    // Perform logout
    useAuthStore.getState().logout();
    
    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });
  
  it('should handle setUser with null appropriately', () => {
    const mockUser: User = { id: '1', email: 'test@example.com', role: 'CUSTOMER' };
    useAuthStore.getState().setUser(mockUser);
    expect(useAuthStore.getState().isAuthenticated).toBe(true);
    
    // Set to null
    useAuthStore.getState().setUser(null);
    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });
});
