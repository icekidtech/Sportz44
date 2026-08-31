import React from 'react';
import { View, Text, StyleSheet, Image, Pressable } from 'react-native';
import { colors, radius, spacing } from '../theme/colors';
import { LiveBadge } from './LiveBadge';
import type { Match } from '../api/types';

function formatTime(dateStr: string) {
  try {
    const d = new Date(dateStr);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } catch { return dateStr; }
}

function formatDate(dateStr: string) {
  try {
    const d = new Date(dateStr);
    return d.toLocaleDateString([], { month: 'short', day: 'numeric' });
  } catch { return ''; }
}

export function MatchCard({ match, onPress }: { match: Match; onPress?: () => void }) {
  const isLive = match.status === 'live';
  const isFinished = match.status === 'finished';

  return (
    <Pressable onPress={onPress} style={({ pressed }) => [styles.card, pressed && { opacity: 0.85 }]}>
      {/* Top row: competition + status */}
      <View style={styles.topRow}>
        <Text style={styles.competition} numberOfLines={1}>
          {match.competition?.name ?? '—'} · {formatDate(match.match_date)}
        </Text>
        {isLive ? (
          <LiveBadge minute={match.minute} />
        ) : isFinished ? (
          <View style={styles.ftBadge}><Text style={styles.ftText}>FT</Text></View>
        ) : (
          <Text style={styles.time}>{formatTime(match.match_date)}</Text>
        )}
      </View>

      {/* Teams + score */}
      <View style={styles.teamsRow}>
        <View style={styles.team}>
          {match.home_club?.logo_url ? (
            <Image source={{ uri: match.home_club.logo_url }} style={styles.logo} />
          ) : (
            <View style={[styles.logo, styles.logoPlaceholder]} />
          )}
          <Text style={styles.teamName} numberOfLines={1}>{match.home_club?.name ?? 'Home'}</Text>
        </View>

        <View style={styles.scoreBox}>
          {isLive || isFinished ? (
            <Text style={styles.score}>{match.home_score} — {match.away_score}</Text>
          ) : (
            <Text style={styles.vs}>VS</Text>
          )}
          {match.venue ? <Text style={styles.venue} numberOfLines={1}>{match.venue}</Text> : null}
        </View>

        <View style={styles.team}>
          {match.away_club?.logo_url ? (
            <Image source={{ uri: match.away_club.logo_url }} style={styles.logo} />
          ) : (
            <View style={[styles.logo, styles.logoPlaceholder]} />
          )}
          <Text style={styles.teamName} numberOfLines={1}>{match.away_club?.name ?? 'Away'}</Text>
        </View>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
    marginBottom: spacing.sm,
  },
  topRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 },
  competition: { color: colors.textMuted, fontSize: 11, fontWeight: '600', letterSpacing: 0.3, flex: 1, marginRight: 8 },
  time: { color: colors.textSecondary, fontSize: 12, fontWeight: '600' },
  ftBadge: { backgroundColor: colors.surfaceElevated, paddingHorizontal: 8, paddingVertical: 3, borderRadius: radius.full },
  ftText: { color: colors.textSecondary, fontSize: 10, fontWeight: '800' },
  teamsRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
  team: { flex: 1, alignItems: 'center', gap: 6 },
  logo: { width: 36, height: 36, borderRadius: 18, backgroundColor: colors.surfaceLight },
  logoPlaceholder: { borderWidth: 1, borderColor: colors.border },
  teamName: { color: colors.text, fontSize: 12, fontWeight: '600', textAlign: 'center' },
  scoreBox: { alignItems: 'center', paddingHorizontal: 12, minWidth: 80 },
  score: { color: colors.text, fontSize: 20, fontWeight: '800' },
  vs: { color: colors.textMuted, fontSize: 13, fontWeight: '700', letterSpacing: 1 },
  venue: { color: colors.textMuted, fontSize: 10, marginTop: 2, textAlign: 'center' },
});
