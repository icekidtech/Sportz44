import React, { useState, useRef } from 'react';
import { View, Text, StyleSheet, Pressable, Dimensions, ImageBackground } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, radius } from '../theme/colors';

const { width } = Dimensions.get('window');

const SLIDES = [
  {
    title: 'Easy Streaming',
    subtitle: 'Choose your plan to watch live match\nyour favourite club.',
    image: null as string | null,
  },
  {
    title: 'Always Uptodate',
    subtitle: 'Stay updated with match highlight,\npreview and hot news',
    image: null as string | null,
  },
  {
    title: 'Get Amazing Reward',
    subtitle: 'Redeem your points to get an\namazing reward',
    image: null as string | null,
  },
];

export function OnboardingScreen({ navigation, onComplete }: { navigation: any; onComplete: () => void }) {
  const [index, setIndex] = useState(0);
  const isLast = index === SLIDES.length - 1;
  const slide = SLIDES[index];

  const handleNext = () => {
    if (isLast) onComplete();
    else setIndex((i) => i + 1);
  };

  const handleBack = () => {
    if (index > 0) setIndex((i) => i - 1);
  };

  const handleSkip = () => onComplete();

  return (
    <View style={styles.container}>
      {/* Full-bleed image area with gradient overlay */}
      <View style={styles.imageArea}>
        {/* Placeholder for hero image — replace with actual ImageBackground when assets are ready */}
        <View style={styles.imagePlaceholder}>
          <Text style={styles.placeholderText}>🏟️</Text>
        </View>
        <View style={styles.gradient} />
      </View>

      {/* Bottom content */}
      <View style={styles.bottom}>
        <Text style={styles.title}>{slide.title}</Text>
        <Text style={styles.subtitle}>{slide.subtitle}</Text>

        {/* Dots */}
        <View style={styles.dots}>
          {SLIDES.map((_, i) => (
            <View key={i} style={[styles.dot, i === index && styles.dotActive]} />
          ))}
        </View>

        {/* CTA */}
        <Pressable style={styles.cta} onPress={handleNext}>
          <Text style={styles.ctaText}>{isLast ? "Let's Get Started" : 'Next'}</Text>
        </Pressable>

        {/* Back / Skip */}
        <Pressable onPress={isLast ? handleBack : handleSkip} style={styles.secondaryBtn}>
          <Text style={styles.secondaryText}>{isLast ? 'Back' : 'Skip'}</Text>
        </Pressable>
        {!isLast && index > 0 && (
          <Pressable onPress={handleBack} style={styles.backLink}>
            <Text style={styles.backLinkText}>Back</Text>
          </Pressable>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#121212' },
  imageArea: { flex: 1, position: 'relative' },
  imagePlaceholder: {
    flex: 1,
    backgroundColor: '#1A1A1A',
    alignItems: 'center',
    justifyContent: 'center',
  },
  placeholderText: { fontSize: 48, opacity: 0.3 },
  gradient: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    height: 220,
    backgroundColor: 'rgba(18,18,18,0.85)',
  },
  bottom: {
    backgroundColor: '#121212',
    paddingHorizontal: 24,
    paddingTop: 16,
    paddingBottom: 32,
    alignItems: 'center',
  },
  title: { color: '#FFFFFF', fontSize: 18, fontWeight: '700', textAlign: 'center' },
  subtitle: { color: '#9CA3AF', fontSize: 12, textAlign: 'center', marginTop: 8, lineHeight: 18 },
  dots: { flexDirection: 'row', gap: 6, marginTop: 20, marginBottom: 20 },
  dot: { width: 6, height: 6, borderRadius: 3, backgroundColor: '#3A3A3A' },
  dotActive: { backgroundColor: '#FFFFFF', width: 18 },
  cta: {
    backgroundColor: '#0EA5B5', // Figma teal
    borderRadius: 24,
    paddingVertical: 14,
    alignItems: 'center',
    width: '100%',
  },
  ctaText: { color: '#FFFFFF', fontSize: 13, fontWeight: '700' },
  secondaryBtn: { marginTop: 12, paddingVertical: 6 },
  secondaryText: { color: '#6B7280', fontSize: 12, fontWeight: '500' },
  backLink: { marginTop: 4 },
  backLinkText: { color: '#6B7280', fontSize: 12 },
});
