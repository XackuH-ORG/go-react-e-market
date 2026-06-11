import { describe, it, expect } from 'vitest';
import { cn } from './utils';

describe('utils', () => {
  describe('cn', () => {
    it('should merge tailwind classes correctly', () => {
      expect(cn('bg-red-500', 'bg-blue-500')).toBe('bg-blue-500');
    });

    it('should conditionally apply classes', () => {
      expect(cn('p-4', { 'm-4': true, 'm-2': false })).toBe('p-4 m-4');
    });

    it('should handle undefined and null values', () => {
      expect(cn('text-center', null, undefined, 'font-bold')).toBe('text-center font-bold');
    });

    it('should handle arrays of classes', () => {
      expect(cn(['text-sm', 'leading-none'], 'font-medium')).toBe('text-sm leading-none font-medium');
    });
  });
});
