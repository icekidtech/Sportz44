import React, { useState, useCallback } from 'react';
import 'react-native-gesture-handler';
import 'react-native-reanimated';
import { StatusBar } from 'expo-status-bar';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import * as Splash from 'expo-splash-screen';
import { AuthProvider } from './src/context/AuthContext';
import { RootNavigator } from './src/navigation/RootNavigator';
import { SplashScreen } from './src/screens/SplashScreen';

Splash.preventAutoHideAsync().catch(() => {});

export default function App() {
  const [showSplash, setShowSplash] = useState(true);

  const onSplashFinish = useCallback(async () => {
    setShowSplash(false);
    await Splash.hideAsync().catch(() => {});
  }, []);

  if (showSplash) {
    return (
      <>
        <SplashScreen onFinish={onSplashFinish} />
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
