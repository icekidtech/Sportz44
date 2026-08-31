import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { colors, radius } from '../theme/colors';

export function LiveBadge({ minute }: { minute?: number }) {
  return (
    <View style={styles.badge}>
      <View style={styles.dot} />
      <Text style={styles.text}>LIVE{minute ? ` ${minute}'` : ''}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.liveBg,
    borderWidth: 1,
    borderColor: 'rgba(239,68,68,0.3)',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: radius.full,
    gap: 5,
  },
  dot: { width: 6, height: 6, borderRadius: 3, backgroundColor: colors.live },
  text: { color: colors.live, fontSize: 10, fontWeight: '800', letterSpacing: 0.5 },
});
