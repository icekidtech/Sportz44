export const colors = {
  // NowTV-inspired dark theme
  background: '#0A0E1A',
  surface: '#141A2A',
  surfaceLight: '#1C2333',
  surfaceElevated: '#242E44',
  card: '#1A2235',
  border: '#2A3447',
  borderLight: '#334155',

  // Brand
  primary: '#00D084',
  primaryDark: '#00B070',
  primaryLight: '#1AE89A',
  accent: '#3B82F6',
  accentLight: '#60A5FA',

  // Text
  text: '#FFFFFF',
  textSecondary: '#94A3B8',
  textMuted: '#64748B',
  textOnPrimary: '#0A0E1A',

  // Status
  live: '#EF4444',
  liveBg: 'rgba(239,68,68,0.15)',
  success: '#00D084',
  warning: '#F59E0B',
  error: '#EF4444',

  // Overlays
  overlay: 'rgba(10,14,26,0.85)',
  cardOverlay: 'rgba(0,0,0,0.4)',
} as const;

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
} as const;

export const radius = {
  sm: 8,
  md: 12,
  lg: 16,
  xl: 20,
  full: 9999,
} as const;

export const typography = {
  h1: { fontSize: 28, fontWeight: '800' as const, lineHeight: 34 },
  h2: { fontSize: 22, fontWeight: '700' as const, lineHeight: 28 },
  h3: { fontSize: 18, fontWeight: '700' as const, lineHeight: 24 },
  h4: { fontSize: 16, fontWeight: '600' as const, lineHeight: 22 },
  body: { fontSize: 14, fontWeight: '400' as const, lineHeight: 20 },
  bodyBold: { fontSize: 14, fontWeight: '600' as const, lineHeight: 20 },
  caption: { fontSize: 12, fontWeight: '500' as const, lineHeight: 16 },
  small: { fontSize: 11, fontWeight: '500' as const, lineHeight: 14 },
} as const;
