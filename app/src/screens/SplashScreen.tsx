import React, { useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withTiming,
  withDelay,
  Easing,
} from 'react-native-reanimated';
import { colors } from '../theme/colors';

const WORD = 'Sportz44';
const LETTER_DELAY = 120;
const LETTER_DURATION = 280;

function AnimatedLetter({ char, index }: { char: string; index: number }) {
  const opacity = useSharedValue(0.18);
  const scale = useSharedValue(0.96);

  useEffect(() => {
    opacity.value = withDelay(index * LETTER_DELAY, withTiming(1, { duration: LETTER_DURATION, easing: Easing.out(Easing.cubic) }));
    scale.value = withDelay(index * LETTER_DELAY, withTiming(1, { duration: LETTER_DURATION, easing: Easing.out(Easing.cubic) }));
  }, []);

  const style = useAnimatedStyle(() => ({
    opacity: opacity.value,
    transform: [{ scale: scale.value }],
  }));

  return (
    <Animated.Text style={[styles.letter, style]}>{char}</Animated.Text>
  );
}

export function SplashScreen({ onFinish }: { onFinish: () => void }) {
  useEffect(() => {
    const total = WORD.length * LETTER_DELAY + LETTER_DURATION + 500;
    const t = setTimeout(onFinish, total);
    return () => clearTimeout(t);
  }, [onFinish]);

  return (
    <View style={styles.container}>
      <View style={styles.wordRow}>
        {WORD.split('').map((ch, i) => (
          <AnimatedLetter key={`${ch}-${i}`} char={ch} index={i} />
        ))}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#121212', // main-bg from Figma
    alignItems: 'center',
    justifyContent: 'center',
  },
  wordRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  letter: {
    color: '#FFFFFF',
    fontSize: 28,
    fontWeight: '800',
    fontStyle: 'italic',
    letterSpacing: -0.5,
  },
});
