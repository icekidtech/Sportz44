import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, ActivityIndicator, Pressable } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { MatchCard } from '../components/MatchCard';
import { SectionHeader } from '../components/SectionHeader';
import { LiveBadge } from '../components/LiveBadge';
import { api } from '../api/client';
import type { Match } from '../api/types';

export function HomeScreen({ navigation }: any) {
  const [live, setLive] = useState<Match[]>([]);
  const [upcoming, setUpcoming] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const load = useCallback(async () => {
    try {
      const [liveRes, upcomingRes] = await Promise.all([
        api.get<{ data: Match[] }>('/api/matches?status=live').catch(() => ({ data: [] as Match[] })),
        api.get<{ data: Match[] }>('/api/matches?status=scheduled&limit=6').catch(() => ({ data: [] as Match[] })),
      ]);
      setLive(liveRes.data ?? []);
      setUpcoming(upcomingRes.data ?? []);
    } catch {}
    setLoading(false);
    setRefreshing(false);
  }, []);

  useEffect(() => { load(); }, [load]);

  const onRefresh = useCallback(() => { setRefreshing(true); load(); }, [load]);

  if (loading) {
    return (
      <View style={styles.center}><ActivityIndicator color={colors.primary} size="large" /></View>
    );
  }

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={colors.primary} />}
      >
        {/* Header */}
        <View style={styles.header}>
          <Text style={styles.greeting}>Sportz44</Text>
          <Text style={styles.subtitle}>Your football intelligence hub</Text>
        </View>

        {/* Live now */}
        <SectionHeader title="Live Now" action={live.length > 0 ? `${live.length} matches` : undefined} />
        {live.length === 0 ? (
          <View style={styles.emptyCard}>
            <Text style={styles.emptyText}>No live matches right now</Text>
            <Text style={styles.emptySub}>Check back during match hours</Text>
          </View>
        ) : (
          live.map((m) => (
            <MatchCard key={m.id} match={m} onPress={() => navigation.navigate('MatchDetail', { id: m.id })} />
          ))
        )}

        {/* Upcoming */}
        <SectionHeader title="Upcoming" action="See all" onAction={() => navigation.navigate('Fixtures')} />
        {upcoming.length === 0 ? (
          <View style={styles.emptyCard}><Text style={styles.emptyText}>No upcoming fixtures</Text></View>
        ) : (
          upcoming.map((m) => (
            <MatchCard key={m.id} match={m} onPress={() => navigation.navigate('MatchDetail', { id: m.id })} />
          ))
        )}

        {/* Quick stats teaser */}
        <View style={styles.promoCard}>
          <Text style={styles.promoTitle}>Follow your clubs</Text>
          <Text style={styles.promoText}>Get live updates, stats & community for every match.</Text>
          <Pressable style={styles.promoBtn} onPress={() => navigation.navigate('Standings')}>
            <Text style={styles.promoBtnText}>View Standings</Text>
          </Pressable>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  scroll: { flex: 1 },
  content: { padding: spacing.md, paddingBottom: 32 },
  center: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  header: { marginBottom: spacing.lg, marginTop: spacing.sm },
  greeting: { color: colors.text, fontSize: 26, fontWeight: '800' },
  subtitle: { color: colors.textMuted, fontSize: 13, marginTop: 2 },
  emptyCard: {
    backgroundColor: colors.card, borderRadius: radius.lg, borderWidth: 1, borderColor: colors.border,
    padding: spacing.lg, alignItems: 'center', marginBottom: spacing.sm,
  },
  emptyText: { color: colors.textSecondary, fontSize: 13, fontWeight: '600' },
  emptySub: { color: colors.textMuted, fontSize: 11, marginTop: 4 },
  promoCard: {
    backgroundColor: colors.surfaceElevated, borderRadius: radius.lg, padding: spacing.lg,
    marginTop: spacing.lg, borderWidth: 1, borderColor: colors.border,
  },
  promoTitle: { color: colors.text, fontSize: 15, fontWeight: '700' },
  promoText: { color: colors.textSecondary, fontSize: 12, marginTop: 4, lineHeight: 18 },
  promoBtn: { backgroundColor: colors.primary, borderRadius: radius.md, paddingVertical: 10, alignItems: 'center', marginTop: 12 },
  promoBtnText: { color: colors.textOnPrimary, fontSize: 13, fontWeight: '700' },
});
