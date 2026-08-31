import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, ActivityIndicator, Image, Pressable } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { LiveBadge } from '../components/LiveBadge';
import { api } from '../api/client';
import type { Match, MatchEvent, MatchLineup, MatchStat } from '../api/types';

type Tab = 'events' | 'lineup' | 'stats';

export function MatchDetailScreen({ route, navigation }: any) {
  const { id } = route.params;
  const [match, setMatch] = useState<Match | null>(null);
  const [events, setEvents] = useState<MatchEvent[]>([]);
  const [lineup, setLineup] = useState<MatchLineup[]>([]);
  const [stats, setStats] = useState<MatchStat[]>([]);
  const [tab, setTab] = useState<Tab>('events');
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      const m = await api.get<Match>(`/api/matches/${id}`);
      setMatch(m as any);
      const [ev, lu, st] = await Promise.all([
        api.get<{ events: MatchEvent[] }>(`/api/matches/${id}/events`).catch(() => ({ events: [] as MatchEvent[] })),
        api.get<{ lineup: MatchLineup[] }>(`/api/matches/${id}/lineup`).catch(() => ({ lineup: [] as MatchLineup[] })),
        api.get<{ stats: MatchStat[] }>(`/api/matches/${id}/stats`).catch(() => ({ stats: [] as MatchStat[] })),
      ]);
      setEvents(ev.events ?? []);
      setLineup(lu.lineup ?? []);
      setStats(st.stats ?? []);
    } catch {}
    setLoading(false);
  }, [id]);

  useEffect(() => { load(); }, [load]);

  if (loading) return <View style={styles.center}><ActivityIndicator color={colors.primary} /></View>;
  if (!match) return <View style={styles.center}><Text style={styles.muted}>Match not found</Text></View>;

  const isLive = match.status === 'live';

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <ScrollView contentContainerStyle={styles.content}>
        {/* Back + header */}
        <Pressable onPress={() => navigation.goBack()}><Text style={styles.back}>‹ Back</Text></Pressable>

        <View style={styles.scoreCard}>
          <View style={styles.teamsRow}>
            <View style={styles.team}>
              {match.home_club?.logo_url ? <Image source={{ uri: match.home_club.logo_url }} style={styles.logo} /> : <View style={[styles.logo, styles.logoPh]} />}
              <Text style={styles.teamName}>{match.home_club?.name ?? 'Home'}</Text>
            </View>
            <View style={styles.scoreBox}>
              <Text style={styles.score}>{match.home_score} — {match.away_score}</Text>
              {isLive ? <LiveBadge minute={match.minute} /> : <Text style={styles.status}>{match.status}</Text>}
            </View>
            <View style={styles.team}>
              {match.away_club?.logo_url ? <Image source={{ uri: match.away_club.logo_url }} style={styles.logo} /> : <View style={[styles.logo, styles.logoPh]} />}
              <Text style={styles.teamName}>{match.away_club?.name ?? 'Away'}</Text>
            </View>
          </View>
          {match.venue ? <Text style={styles.venue}>{match.venue}</Text> : null}
        </View>

        {/* Tabs */}
        <View style={styles.tabs}>
          {(['events', 'lineup', 'stats'] as Tab[]).map((t) => (
            <Pressable key={t} onPress={() => setTab(t)} style={[styles.tab, tab === t && styles.tabActive]}>
              <Text style={[styles.tabText, tab === t && styles.tabTextActive]}>{t.charAt(0).toUpperCase() + t.slice(1)}</Text>
            </Pressable>
          ))}
        </View>

        {/* Tab content */}
        {tab === 'events' && (
          <View>
            {events.length === 0 ? <Text style={styles.muted}>No events yet</Text> : events.map((e) => (
              <View key={e.id} style={styles.eventRow}>
                <Text style={styles.eventMinute}>{e.minute}'</Text>
                <View style={styles.eventDot} />
                <View style={{ flex: 1 }}>
                  <Text style={styles.eventType}>{e.event_type} — {e.player_name}</Text>
                  {e.detail ? <Text style={styles.eventDetail}>{e.detail}</Text> : null}
                  {e.assist_player_name ? <Text style={styles.eventAssist}>Assist: {e.assist_player_name}</Text> : null}
                </View>
              </View>
            ))}
          </View>
        )}

        {tab === 'lineup' && (
          <View>
            {lineup.length === 0 ? <Text style={styles.muted}>No lineup available</Text> : lineup.map((p) => (
              <View key={p.id} style={styles.lineupRow}>
                <Text style={styles.jersey}>{p.jersey_number}</Text>
                <Text style={styles.playerName}>{p.player_name}</Text>
                <Text style={styles.position}>{p.position}</Text>
                <Text style={styles.starter}>{p.is_starter ? 'XI' : 'Sub'}</Text>
              </View>
            ))}
          </View>
        )}

        {tab === 'stats' && (
          <View>
            {stats.length === 0 ? <Text style={styles.muted}>No stats available</Text> : stats.map((s) => (
              <View key={s.id} style={styles.statRow}>
                <Text style={styles.statType}>{s.stat_type}</Text>
                <Text style={styles.statValue}>{s.value}</Text>
              </View>
            ))}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  content: { padding: spacing.md, paddingBottom: 32 },
  back: { color: colors.primary, fontSize: 14, fontWeight: '600', marginBottom: 12 },
  scoreCard: { backgroundColor: colors.card, borderRadius: radius.lg, borderWidth: 1, borderColor: colors.border, padding: spacing.md, marginBottom: spacing.md },
  teamsRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
  team: { flex: 1, alignItems: 'center', gap: 6 },
  logo: { width: 40, height: 40, borderRadius: 20 },
  logoPh: { backgroundColor: colors.surfaceLight, borderWidth: 1, borderColor: colors.border },
  teamName: { color: colors.text, fontSize: 12, fontWeight: '600', textAlign: 'center' },
  scoreBox: { alignItems: 'center', gap: 6, paddingHorizontal: 12 },
  score: { color: colors.text, fontSize: 24, fontWeight: '800' },
  status: { color: colors.textMuted, fontSize: 11, fontWeight: '600', textTransform: 'uppercase' },
  venue: { color: colors.textMuted, fontSize: 11, textAlign: 'center', marginTop: 8 },
  tabs: { flexDirection: 'row', backgroundColor: colors.surface, borderRadius: radius.md, padding: 4, marginBottom: spacing.md, gap: 4 },
  tab: { flex: 1, paddingVertical: 8, alignItems: 'center', borderRadius: radius.sm },
  tabActive: { backgroundColor: colors.primary },
  tabText: { color: colors.textSecondary, fontSize: 12, fontWeight: '600' },
  tabTextActive: { color: colors.textOnPrimary },
  muted: { color: colors.textMuted, fontSize: 13, textAlign: 'center', marginTop: 16 },
  eventRow: { flexDirection: 'row', alignItems: 'flex-start', paddingVertical: 10, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.4)', gap: 10 },
  eventMinute: { color: colors.primary, fontSize: 12, fontWeight: '700', width: 32, textAlign: 'right' },
  eventDot: { width: 8, height: 8, borderRadius: 4, backgroundColor: colors.primary, marginTop: 4 },
  eventType: { color: colors.text, fontSize: 13, fontWeight: '600' },
  eventDetail: { color: colors.textSecondary, fontSize: 11, marginTop: 2 },
  eventAssist: { color: colors.textMuted, fontSize: 11, marginTop: 1, fontStyle: 'italic' },
  lineupRow: { flexDirection: 'row', alignItems: 'center', paddingVertical: 8, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.4)', gap: 10 },
  jersey: { color: colors.primary, fontSize: 13, fontWeight: '800', width: 28, textAlign: 'center' },
  playerName: { color: colors.text, fontSize: 13, fontWeight: '500', flex: 1 },
  position: { color: colors.textMuted, fontSize: 11, width: 60, textAlign: 'right' },
  starter: { color: colors.textSecondary, fontSize: 10, fontWeight: '700', width: 28, textAlign: 'center' },
  statRow: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: 8, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.4)' },
  statType: { color: colors.textSecondary, fontSize: 12 },
  statValue: { color: colors.text, fontSize: 12, fontWeight: '700' },
});
