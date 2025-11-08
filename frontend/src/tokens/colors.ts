/**
 * Design Tokens: Colors
 * Based on Pinecone brand identity
 */

export const colors = {
  // Primary palette
  primary: '#3A7D44',        // Forest green
  primaryLight: '#4A9D54',
  primaryDark: '#2A5D34',

  // Background
  background: '#FAF9F6',     // Warm off-white
  backgroundDark: '#F0EFE9',

  // Text
  text: '#121212',           // Black
  textSecondary: '#4A4A4A',
  textMuted: '#6B6B6B',

  // UI elements
  border: '#D1D1D1',
  borderLight: '#E5E5E5',

  // Feedback
  error: '#D32F2F',
  errorLight: '#FFCDD2',
  success: '#388E3C',
  successLight: '#C8E6C9',
  warning: '#F57C00',
  warningLight: '#FFE0B2',
  info: '#1976D2',
  infoLight: '#BBDEFB',

  // Neutrals
  white: '#FFFFFF',
  black: '#000000',
  gray100: '#F5F5F5',
  gray200: '#EEEEEE',
  gray300: '#E0E0E0',
  gray400: '#BDBDBD',
  gray500: '#9E9E9E',
  gray600: '#757575',
  gray700: '#616161',
  gray800: '#424242',
  gray900: '#212121',
} as const

export type ColorKey = keyof typeof colors
