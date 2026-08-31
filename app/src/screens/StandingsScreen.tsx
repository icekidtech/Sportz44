import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, ActivityIndicator, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { api } from '../api/client';
import type { Standing, TopScorer } from '../api/types';

export function StandingsScreen() {
  const [standings, setStandings] = useState<Standing[]>([]);
  const [scorers, setScorers] = useState<TopScorer[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const currentSeason = String(new Date().getFullYear());

  const load = useCallback(async () => {
    try {
      const [s, t] = await Promise.all([
        api.get<{ standings: Standing[] }>(`/api/standings?league=1&season=${currentSeason}`).catch(() => ({ standings: [] as Standing[] })),
        api.get<{ top_scorers: TopScorer[] }>(`/api/standings/1/top-scorers?season=${currentSeason}`).catch(() => ({ top_scorers: [] as TopScorer[] })),
      ]);
      // Fall back to all-time if current season has no data yet
      if ((s.standings ?? []).length === 0) {
        const fallback = await api.get<{ standings: Standing[] }>('/api/standings?league=1').catch(() => ({ standings: [] as Standing[] }));
        setStandings(fallback.standings ?? []);
      } else {
        setStandings(s.standings ?? []);
      }
      setScorers(t.top_scorers ?? []);
    } catch {}
    setLoading(false);
    setRefreshing(false);
  }, [currentSeason]);

  useEffect(() => { load(); }, [load]);

  if (loading) return <View style={styles.center}><ActivityIndicator color={colors.primary} /></View>;

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <ScrollView
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={() => { setRefreshing(true); load(); }} tintColor={colors.primary} />}
      >
        <Text style={styles.title}>Standings</Text>

        {/* Table header */}
        <View style={[styles.row, styles.headerRow]}>
          <Text style={[styles.cell, styles.pos]}>#</Text>
          <Text style={[styles.cell, styles.club]}>Club</Text>
          <Text style={styles.num}>P</Text>
          <Text style={styles.num}>W</Text>
          <Text style={styles.num}>D</Text>
          <Text style={styles.num}>L</Text>
          <Text style={styles.num}>GD</Text>
          <Text style={[styles.num, styles.pts]}>Pts</Text>
        </View>

        {standings.length === 0 ? (
          <View style={styles.empty}><Text style={styles.emptyText}>No standings yet</Text></View>
        ) : (
          standings.map((s) => (
            <View key={s.club_id} style={styles.row}>
              <Text style={[styles.cell, styles.pos]}>{s.position}</Text>
              <View style={{ flexDirection: 'row', alignItems: 'center', gap: 6, flex: 1 }}>
                {s.logo_url ? <Image source={{ uri: s.logo_url }} style={styles.logo} /> : null}
                <Text style={styles.clubName} numberOfLines={1}>{s.club_name}</Text>
              </View>
              <Text style={styles.num}>{s.played}</Text>
              <Text style={styles.num}>{s.won}</Text>
              <Text style={styles.num}>{s.drawn}</Text>
              <Text style={styles.num}>{s.lost}</Text>
              <Text style={styles.num}>{s.goal_difference}</Text>
              <Text style={[styles.num, styles.pts]}>{s.points}</Text>
            </View>
          ))
        )}

        {/* Top scorers */}
        {scorers.length > 0 && (
          <>
            <Text style={[styles.title, { marginTop: spacing.lg }]}>Top Scorers</Text>
            {scorers.slice(0, 10).map((p, i) => (
              <View key={p.player_id} style={styles.scorerRow}>
                <Text style={styles.scorerRank}>{i + 1}</Text>
                <View style={{ flex: 1 }}>
                  <Text style={styles.scorerName}>{p.player_name}</Text>
                  <Text style={styles.scorerClub}>{p.club_name}</Text>
                </View>
                <Text style={styles.scorerGoals}>{p.goals}</Text>
              </View>
            ))}
          </>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  content: { padding: spacing.md, paddingBottom: 32 },
  title: { color: colors.text, fontSize: 22, fontWeight: '800', marginBottom: 12 },
  headerRow: { borderBottomWidth: 1, borderBottomColor: colors.border, paddingBottom: 8, marginBottom: 4 },
  row: { flexDirection: 'row', alignItems: 'center', paddingVertical: 8, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.5)' },
  cell: { color: colors.textSecondary, fontSize: 12 },
  pos: { width: 28, textAlign: 'center', color: colors.textMuted, fontWeight: '700' },
  club: { flex: 1 },
  clubName: { color: colors.text, fontSize: 12, fontWeight: '600', flex: 1 },
  logo: { width: 18, height: 18, borderRadius: 9 },
  num: { width: 28, textAlign: 'center', color: colors.textSecondary, fontSize: 12, fontWeight: '500' },
  pts: { color: colors.text, fontWeight: '700' },
  empty: { alignItems: 'center', paddingTop: 24 },
  emptyText: { color: colors.textMuted, fontSize: 13 },
  scorerRow: { flexDirection: 'row', alignItems: 'center', paddingVertical: 10, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.5)', gap: 12 },
  scorerRank: { width: 24, textAlign: 'center', color: colors.textMuted, fontSize: 12, fontWeight: '700' },
  scorerName: { color: colors.text, fontSize: 13, fontWeight: '600' },
  scorerClub: { color: colors.textMuted, fontSize: 11, marginTop: 1 },
  scorerGoals: { color: colors.primary, fontSize: 16, fontWeight: '800', minWidth: 24, textAlign: 'right' },
});
