import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, FlatList, RefreshControl, ActivityIndicator, Pressable } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { MatchCard } from '../components/MatchCard';
import { api } from '../api/client';
import type { Match } from '../api/types';

const FILTERS = [
  { label: 'All', value: '' },
  { label: 'Live', value: 'live' },
  { label: 'Finished', value: 'finished' },
  { label: 'Upcoming', value: 'scheduled' },
] as const;

export function FixturesScreen({ navigation }: any) {
  const [filter, setFilter] = useState('');
  const [matches, setMatches] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const load = useCallback(async () => {
    try {
      const q = filter ? `?status=${filter}` : '';
      const res = await api.get<{ data: Match[] }>(`/api/matches${q}`);
      setMatches(res.data ?? []);
    } catch { setMatches([]); }
    setLoading(false);
    setRefreshing(false);
  }, [filter]);

  useEffect(() => { setLoading(true); load(); }, [load]);

  if (loading) return <View style={styles.center}><ActivityIndicator color={colors.primary} /></View>;

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <View style={styles.header}>
        <Text style={styles.title}>Fixtures</Text>
        <View style={styles.filters}>
          {FILTERS.map((f) => (
            <Pressable
              key={f.value}
              onPress={() => setFilter(f.value)}
              style={[styles.chip, filter === f.value && styles.chipActive]}
            >
              <Text style={[styles.chipText, filter === f.value && styles.chipTextActive]}>{f.label}</Text>
            </Pressable>
          ))}
        </View>
      </View>
      <FlatList
        data={matches}
        keyExtractor={(m) => String(m.id)}
        contentContainerStyle={styles.list}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={() => { setRefreshing(true); load(); }} tintColor={colors.primary} />}
        ListEmptyComponent={<View style={styles.empty}><Text style={styles.emptyText}>No fixtures</Text></View>}
        renderItem={({ item }) => <MatchCard match={item} onPress={() => navigation.navigate('MatchDetail', { id: item.id })} />}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  header: { padding: spacing.md, paddingBottom: spacing.sm },
  title: { color: colors.text, fontSize: 22, fontWeight: '800' },
  filters: { flexDirection: 'row', gap: 8, marginTop: 12 },
  chip: { paddingHorizontal: 14, paddingVertical: 7, borderRadius: radius.full, backgroundColor: colors.surface, borderWidth: 1, borderColor: colors.border },
  chipActive: { backgroundColor: colors.primary, borderColor: colors.primary },
  chipText: { color: colors.textSecondary, fontSize: 12, fontWeight: '600' },
  chipTextActive: { color: colors.textOnPrimary },
  list: { padding: spacing.md, paddingTop: 0 },
  empty: { alignItems: 'center', paddingTop: 48 },
  emptyText: { color: colors.textMuted, fontSize: 13 },
});
