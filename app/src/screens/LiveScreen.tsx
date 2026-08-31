import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, FlatList, RefreshControl, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing } from '../theme/colors';
import { MatchCard } from '../components/MatchCard';
import { api } from '../api/client';
import type { Match } from '../api/types';

export function LiveScreen({ navigation }: any) {
  const [matches, setMatches] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const load = useCallback(async () => {
    try {
      const res = await api.get<{ data: Match[] }>('/api/matches?status=live');
      setMatches(res.data ?? []);
    } catch { setMatches([]); }
    setLoading(false);
    setRefreshing(false);
  }, []);

  useEffect(() => {
    load();
    const id = setInterval(load, 30000);
    return () => clearInterval(id);
  }, [load]);

  if (loading) return <View style={styles.center}><ActivityIndicator color={colors.primary} /></View>;

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <View style={styles.header}>
        <Text style={styles.title}>Live Matches</Text>
        <Text style={styles.subtitle}>Auto-refreshes every 30s</Text>
      </View>
      <FlatList
        data={matches}
        keyExtractor={(m) => String(m.id)}
        contentContainerStyle={styles.list}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={() => { setRefreshing(true); load(); }} tintColor={colors.primary} />}
        ListEmptyComponent={<View style={styles.empty}><Text style={styles.emptyText}>No live matches</Text></View>}
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
  subtitle: { color: colors.textMuted, fontSize: 12, marginTop: 2 },
  list: { padding: spacing.md, paddingTop: 0 },
  empty: { alignItems: 'center', paddingTop: 48 },
  emptyText: { color: colors.textMuted, fontSize: 13 },
});
