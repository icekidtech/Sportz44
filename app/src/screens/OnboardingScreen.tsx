import React, { useState } from 'react';
import { View, Text, StyleSheet, Pressable, ImageBackground } from 'react-native';
import { colors } from '../theme/colors';

const SLIDES = [
  {
    title: 'Easy Streaming',
    subtitle: 'Choose your plan to watch live match\nyour favourite club.',
    image: require('../../assets/onboarding/slide1.jpg'),
  },
  {
    title: 'Always Uptodate',
    subtitle: 'Stay updated with match highlight,\npreview and hot news',
    image: require('../../assets/onboarding/slide2.jpg'),
  },
  {
    title: 'Get Amazing Reward',
    subtitle: 'Redeem your points to get an\namazing reward',
    image: require('../../assets/onboarding/slide3.jpg'),
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
      <ImageBackground source={slide.image} style={styles.imageArea} resizeMode="cover">
        <View style={styles.gradient} />
      </ImageBackground>

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
  imageArea: { flex: 1, justifyContent: 'flex-end' },
  gradient: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    height: 280,
    backgroundColor: 'rgba(18,18,18,0.75)',
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
