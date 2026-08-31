import React, { useState, useCallback, useEffect } from 'react';
import 'react-native-gesture-handler';
import 'react-native-reanimated';
import { StatusBar } from 'expo-status-bar';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import * as Splash from 'expo-splash-screen';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { AuthProvider } from './src/context/AuthContext';
import { RootNavigator } from './src/navigation/RootNavigator';
import { SplashScreen } from './src/screens/SplashScreen';
import { OnboardingScreen } from './src/screens/OnboardingScreen';

Splash.preventAutoHideAsync().catch(() => {});

const ONBOARDING_KEY = 'sportz44_onboarding_done';

export default function App() {
  const [showSplash, setShowSplash] = useState(true);
  const [showOnboarding, setShowOnboarding] = useState<boolean | null>(null);

  useEffect(() => {
    AsyncStorage.getItem(ONBOARDING_KEY).then((v) => setShowOnboarding(v !== '1'));
  }, []);

  const onSplashFinish = useCallback(async () => {
    setShowSplash(false);
    await Splash.hideAsync().catch(() => {});
  }, []);

  const onOnboardingComplete = useCallback(async () => {
    await AsyncStorage.setItem(ONBOARDING_KEY, '1');
    setShowOnboarding(false);
  }, []);

  const onOnboardingSkip = useCallback(async () => {
    await AsyncStorage.setItem(ONBOARDING_KEY, '1');
    setShowOnboarding(false);
  }, []);

  if (showSplash) {
    return (
      <>
        <SplashScreen onFinish={onSplashFinish} />
        <StatusBar style="light" />
      </>
    );
  }

  if (showOnboarding === null) return null;

  if (showOnboarding) {
    return (
      <>
        <OnboardingScreen onComplete={onOnboardingComplete} navigation={null as any} />
        <StatusBar style="light" />
      </>
    );
  }

  return (
    <SafeAreaProvider>
      <AuthProvider>
        <RootNavigator />
        <StatusBar style="light" />
      </AuthProvider>
    </SafeAreaProvider>
  );
}
