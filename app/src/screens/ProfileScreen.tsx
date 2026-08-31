import React from 'react';
import { View, Text, StyleSheet, Pressable, ScrollView, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { useAuth } from '../context/AuthContext';

export function ProfileScreen({ navigation }: any) {
  const { user, loading, logout } = useAuth();

  if (loading) return <View style={styles.center}><ActivityIndicator color={colors.primary} /></View>;

  if (!user) {
    return (
      <SafeAreaView style={styles.safe} edges={['top']}>
        <View style={styles.center}>
          <Text style={styles.title}>Welcome to Sportz44</Text>
          <Text style={styles.subtitle}>Sign in to follow clubs and get live updates</Text>
          <Pressable style={styles.primaryBtn} onPress={() => navigation.navigate('Login')}>
            <Text style={styles.primaryBtnText}>Sign In</Text>
          </Pressable>
          <Pressable style={styles.secondaryBtn} onPress={() => navigation.navigate('Register')}>
            <Text style={styles.secondaryBtnText}>Create Account</Text>
          </Pressable>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <ScrollView contentContainerStyle={styles.content}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{user.username.charAt(0).toUpperCase()}</Text>
        </View>
        <Text style={styles.name}>{user.username}</Text>
        <Text style={styles.email}>{user.email}</Text>
        <View style={styles.badge}><Text style={styles.badgeText}>{user.role}</Text></View>

        <View style={styles.card}>
          <Pressable style={styles.row} onPress={() => {}}>
            <Text style={styles.rowLabel}>My Clubs</Text>
            <Text style={styles.rowValue}>›</Text>
          </Pressable>
          <View style={styles.divider} />
          <Pressable style={styles.row} onPress={() => {}}>
            <Text style={styles.rowLabel}>Notifications</Text>
            <Text style={styles.rowValue}>›</Text>
          </Pressable>
          <View style={styles.divider} />
          <Pressable style={styles.row} onPress={() => {}}>
            <Text style={styles.rowLabel}>Settings</Text>
            <Text style={styles.rowValue}>›</Text>
          </Pressable>
        </View>

        <Pressable style={styles.logoutBtn} onPress={logout}>
          <Text style={styles.logoutText}>Sign Out</Text>
        </Pressable>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, alignItems: 'center', justifyContent: 'center', padding: spacing.lg },
  content: { padding: spacing.md, alignItems: 'center', paddingTop: spacing.xl },
  title: { color: colors.text, fontSize: 22, fontWeight: '800', textAlign: 'center' },
  subtitle: { color: colors.textMuted, fontSize: 13, marginTop: 6, textAlign: 'center' },
  primaryBtn: { backgroundColor: colors.primary, borderRadius: radius.md, paddingVertical: 14, paddingHorizontal: 48, marginTop: 24, width: '100%', alignItems: 'center' },
  primaryBtnText: { color: colors.textOnPrimary, fontSize: 14, fontWeight: '700' },
  secondaryBtn: { borderWidth: 1, borderColor: colors.border, borderRadius: radius.md, paddingVertical: 14, paddingHorizontal: 48, marginTop: 12, width: '100%', alignItems: 'center' },
  secondaryBtnText: { color: colors.text, fontSize: 14, fontWeight: '600' },
  avatar: { width: 72, height: 72, borderRadius: 36, backgroundColor: colors.primary, alignItems: 'center', justifyContent: 'center' },
  avatarText: { color: colors.textOnPrimary, fontSize: 28, fontWeight: '800' },
  name: { color: colors.text, fontSize: 18, fontWeight: '700', marginTop: 12 },
  email: { color: colors.textMuted, fontSize: 13, marginTop: 2 },
  badge: { backgroundColor: colors.surfaceElevated, paddingHorizontal: 10, paddingVertical: 4, borderRadius: radius.full, marginTop: 8 },
  badgeText: { color: colors.textSecondary, fontSize: 11, fontWeight: '600', textTransform: 'uppercase' },
  card: { backgroundColor: colors.card, borderRadius: radius.lg, borderWidth: 1, borderColor: colors.border, width: '100%', marginTop: 24, overflow: 'hidden' },
  row: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: spacing.md },
  rowLabel: { color: colors.text, fontSize: 14, fontWeight: '500' },
  rowValue: { color: colors.textMuted, fontSize: 18 },
  divider: { height: 1, backgroundColor: colors.border, marginHorizontal: spacing.md },
  logoutBtn: { marginTop: 24, paddingVertical: 12 },
  logoutText: { color: colors.error, fontSize: 14, fontWeight: '600' },
});
