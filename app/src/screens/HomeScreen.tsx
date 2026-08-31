import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, ActivityIndicator, Pressable, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, radius } from '../theme/colors';
import { api } from '../api/client';
import type { Match } from '../api/types';

type Tab = 'All' | 'Preview' | 'Highlight' | 'News';

const TABS: Tab[] = ['All', 'Preview', 'Highlight', 'News'];

function TopTabs({ active, onChange }: { active: Tab; onChange: (t: Tab) => void }) {
  return (
    <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={s.tabsRow}>
      {TABS.map((t) => (
        <Pressable key={t} onPress={() => onChange(t)} style={[s.tab, active === t && s.tabActive]}>
          <Text style={[s.tabText, active === t && s.tabTextActive]}>{t}</Text>
        </Pressable>
      ))}
    </ScrollView>
  );
}

function HighlightCard({ match, onPress }: { match: Match; onPress?: () => void }) {
  const isLive = match.status === 'live';
  return (
    <Pressable onPress={onPress} style={s.hlCard}>
      <View style={s.hlImageWrap}>
        <View style={s.hlImagePh}>
          <Ionicons name="football-outline" size={28} color="rgba(255,255,255,0.3)" />
        </View>
        <View style={s.hlOverlay} />
        <View style={s.hlBadge}>
          <Text style={s.hlBadgeText}>{isLive ? 'LIVE' : 'HIGHLIGHT'}</Text>
        </View>
      </View>
      <View style={s.hlBody}>
        <Text style={s.hlTitle} numberOfLines={2}>
          {match.home_club?.name ?? 'Home'} vs {match.away_club?.name ?? 'Away'}
        </Text>
        <Text style={s.hlMeta} numberOfLines={1}>
          {match.competition?.name ?? ''} {match.venue ? `· ${match.venue}` : ''}
        </Text>
      </View>
    </Pressable>
  );
}

function NewsRow({ title, time }: { title: string; time: string }) {
  return (
    <View style={s.newsRow}>
      <View style={s.newsThumb}>
        <Ionicons name="newspaper-outline" size={16} color={colors.textMuted} />
      </View>
      <View style={{ flex: 1 }}>
        <Text style={s.newsTitle} numberOfLines={2}>{title}</Text>
        <Text style={s.newsMeta}>{time}</Text>
      </View>
    </View>
  );
}

export function HomeScreen({ navigation }: any) {
  const [tab, setTab] = useState<Tab>('All');
  const [live, setLive] = useState<Match[]>([]);
  const [upcoming, setUpcoming] = useState<Match[]>([]);
  const [finished, setFinished] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const load = useCallback(async () => {
    try {
      const [liveRes, upcomingRes, finishedRes] = await Promise.all([
        api.get<{ matches: Match[] }>('/api/matches?status=live').catch(() => ({ matches: [] as Match[] })),
        api.get<{ matches: Match[] }>('/api/matches?status=scheduled&limit=6').catch(() => ({ matches: [] as Match[] })),
        api.get<{ matches: Match[] }>('/api/matches?status=finished&limit=6').catch(() => ({ matches: [] as Match[] })),
      ]);
      setLive(liveRes.matches ?? []);
      setUpcoming(upcomingRes.matches ?? []);
      setFinished(finishedRes.matches ?? []);
    } catch {}
    setLoading(false);
    setRefreshing(false);
  }, []);

  useEffect(() => { load(); }, [load]);

  const highlights = [...live, ...finished].slice(0, 6);
  const showHighlights = tab === 'All' || tab === 'Highlight';
  const showNews = tab === 'All' || tab === 'News';
  const showUpcoming = tab === 'All' || tab === 'Preview';

  if (loading) {
    return <View style={s.center}><ActivityIndicator color={colors.primary} size="large" /></View>;
  }

  return (
    <SafeAreaView style={s.safe} edges={['top']}>
      {/* Header */}
      <View style={s.header}>
        <Text style={s.logo}>Sportz44</Text>
        <View style={s.headerIcons}>
          <Pressable style={s.iconBtn}><Ionicons name="search-outline" size={18} color={colors.text} /></Pressable>
          <Pressable style={s.iconBtn}><Ionicons name="notifications-outline" size={18} color={colors.text} /></Pressable>
        </View>
      </View>

      <TopTabs active={tab} onChange={setTab} />

      <ScrollView
        style={s.scroll}
        contentContainerStyle={s.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={() => { setRefreshing(true); load(); }} tintColor={colors.primary} />}
      >
        {/* Match Highlight — horizontal cards */}
        {showHighlights && (
          <>
            <View style={s.sectionHead}>
              <Text style={s.sectionTitle}>Match Highlight</Text>
              <Pressable onPress={() => navigation.navigate('Fixtures')}><Text style={s.seeAll}>See All</Text></Pressable>
            </View>
            {highlights.length === 0 ? (
              <View style={s.emptyCard}><Text style={s.emptyText}>No highlights yet</Text></View>
            ) : (
              <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={s.hlRow}>
                {highlights.map((m) => (
                  <HighlightCard key={m.id} match={m} onPress={() => navigation.navigate('MatchDetail', { id: m.id })} />
                ))}
              </ScrollView>
            )}
          </>
        )}

        {/* Upcoming / Preview */}
        {showUpcoming && (
          <>
            <View style={s.sectionHead}>
              <Text style={s.sectionTitle}>Upcoming</Text>
              <Pressable onPress={() => navigation.navigate('Fixtures')}><Text style={s.seeAll}>See All</Text></Pressable>
            </View>
            {upcoming.length === 0 ? (
              <View style={s.emptyCard}><Text style={s.emptyText}>No upcoming fixtures</Text></View>
            ) : (
              upcoming.map((m) => (
                <Pressable key={m.id} onPress={() => navigation.navigate('MatchDetail', { id: m.id })} style={s.fixtureRow}>
                  <Text style={s.fixtureTeams} numberOfLines={1}>{m.home_club?.name ?? 'Home'} vs {m.away_club?.name ?? 'Away'}</Text>
                  <Text style={s.fixtureMeta}>{new Date(m.match_date).toLocaleDateString()} · {m.venue ?? ''}</Text>
                </Pressable>
              ))
            )}
          </>
        )}

        {/* Trending News — placeholder until news API */}
        {showNews && (
          <>
            <View style={s.sectionHead}>
              <Text style={s.sectionTitle}>Trending News</Text>
              <Pressable><Text style={s.seeAll}>See All</Text></Pressable>
            </View>
            <NewsRow title="Match highlights and analysis coming soon" time="Today" />
            <NewsRow title="Follow your favourite clubs for personalised updates" time="Today" />
            <NewsRow title="Live scores and standings — stay tuned" time="Today" />
          </>
        )}

        {/* Live indicator */}
        {live.length > 0 && (
          <View style={s.liveBanner}>
            <View style={s.liveDot} />
            <Text style={s.liveText}>{live.length} live match{live.length > 1 ? 'es' : ''} now</Text>
            <Pressable onPress={() => navigation.navigate('Live')}><Text style={s.liveLink}>Watch →</Text></Pressable>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const s = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingHorizontal: spacing.md, paddingTop: 8, paddingBottom: 8 },
  logo: { color: colors.text, fontSize: 18, fontWeight: '800', fontStyle: 'italic' },
  headerIcons: { flexDirection: 'row', gap: 8 },
  iconBtn: { width: 32, height: 32, borderRadius: 16, backgroundColor: colors.surface, alignItems: 'center', justifyContent: 'center' },
  tabsRow: { paddingHorizontal: spacing.md, gap: 8, paddingVertical: 8 },
  tab: { paddingHorizontal: 14, paddingVertical: 7, borderRadius: radius.full, backgroundColor: colors.surface, borderWidth: 1, borderColor: colors.border },
  tabActive: { backgroundColor: colors.text, borderColor: colors.text },
  tabText: { color: colors.textSecondary, fontSize: 12, fontWeight: '600' },
  tabTextActive: { color: colors.background },
  scroll: { flex: 1 },
  content: { padding: spacing.md, paddingTop: 8, paddingBottom: 32 },
  sectionHead: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginTop: 16, marginBottom: 10 },
  sectionTitle: { color: colors.text, fontSize: 14, fontWeight: '700' },
  seeAll: { color: colors.textMuted, fontSize: 11, fontWeight: '600' },
  emptyCard: { backgroundColor: colors.card, borderRadius: radius.lg, borderWidth: 1, borderColor: colors.border, padding: spacing.md, alignItems: 'center' },
  emptyText: { color: colors.textMuted, fontSize: 12 },
  hlRow: { gap: 12, paddingRight: spacing.md },
  hlCard: { width: 160, backgroundColor: colors.card, borderRadius: radius.lg, overflow: 'hidden', borderWidth: 1, borderColor: colors.border },
  hlImageWrap: { height: 90, backgroundColor: colors.surfaceLight, position: 'relative' },
  hlImagePh: { flex: 1, alignItems: 'center', justifyContent: 'center' },
  hlOverlay: { ...StyleSheet.absoluteFillObject, backgroundColor: 'rgba(0,0,0,0.2)' },
  hlBadge: { position: 'absolute', top: 8, left: 8, backgroundColor: 'rgba(0,0,0,0.6)', paddingHorizontal: 6, paddingVertical: 2, borderRadius: 4 },
  hlBadgeText: { color: '#fff', fontSize: 9, fontWeight: '700' },
  hlBody: { padding: 10 },
  hlTitle: { color: colors.text, fontSize: 11, fontWeight: '600', lineHeight: 14 },
  hlMeta: { color: colors.textMuted, fontSize: 10, marginTop: 4 },
  fixtureRow: { backgroundColor: colors.card, borderRadius: radius.md, borderWidth: 1, borderColor: colors.border, padding: 12, marginBottom: 8 },
  fixtureTeams: { color: colors.text, fontSize: 12, fontWeight: '600' },
  fixtureMeta: { color: colors.textMuted, fontSize: 10, marginTop: 4 },
  newsRow: { flexDirection: 'row', gap: 10, paddingVertical: 10, borderBottomWidth: 1, borderBottomColor: 'rgba(42,52,71,0.3)' },
  newsThumb: { width: 48, height: 48, borderRadius: 8, backgroundColor: colors.surface, alignItems: 'center', justifyContent: 'center' },
  newsTitle: { color: colors.text, fontSize: 12, fontWeight: '500', lineHeight: 16 },
  newsMeta: { color: colors.textMuted, fontSize: 10, marginTop: 4 },
  liveBanner: { flexDirection: 'row', alignItems: 'center', gap: 8, backgroundColor: 'rgba(239,68,68,0.12)', borderWidth: 1, borderColor: 'rgba(239,68,68,0.25)', borderRadius: radius.md, padding: 12, marginTop: 16 },
  liveDot: { width: 8, height: 8, borderRadius: 4, backgroundColor: '#EF4444' },
  liveText: { color: '#EF4444', fontSize: 12, fontWeight: '600', flex: 1 },
  liveLink: { color: '#EF4444', fontSize: 12, fontWeight: '700' },
});
