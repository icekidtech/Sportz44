import React from 'react';
import { NavigationContainer, DefaultTheme } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Text } from 'react-native';
import { colors } from '../theme/colors';
import { useAuth } from '../context/AuthContext';

import { HomeScreen } from '../screens/HomeScreen';
import { LiveScreen } from '../screens/LiveScreen';
import { FixturesScreen } from '../screens/FixturesScreen';
import { StandingsScreen } from '../screens/StandingsScreen';
import { ProfileScreen } from '../screens/ProfileScreen';
import { MatchDetailScreen } from '../screens/MatchDetailScreen';
import { LoginScreen, RegisterScreen } from '../screens/AuthScreens';

const Stack = createNativeStackNavigator();
const Tab = createBottomTabNavigator();

const navTheme = {
  ...DefaultTheme,
  colors: { ...DefaultTheme.colors, background: colors.background, card: colors.surface, text: colors.text, border: colors.border, primary: colors.primary },
};

function TabIcon({ label, focused }: { label: string; focused: boolean }) {
  return <Text style={{ fontSize: 11, fontWeight: focused ? '800' : '500', color: focused ? colors.primary : colors.textMuted }}>{label}</Text>;
}

function Tabs() {
  return (
    <Tab.Navigator
      screenOptions={{
        headerShown: false,
        tabBarStyle: { backgroundColor: colors.surface, borderTopColor: colors.border, height: 58, paddingBottom: 6, paddingTop: 4 },
        tabBarActiveTintColor: colors.primary,
        tabBarInactiveTintColor: colors.textMuted,
      }}
    >
      <Tab.Screen name="Home" component={HomeScreen} options={{ tabBarIcon: ({ focused }) => <TabIcon label="⌂" focused={focused} />, tabBarLabel: 'Home' }} />
      <Tab.Screen name="Live" component={LiveScreen} options={{ tabBarIcon: ({ focused }) => <TabIcon label="●" focused={focused} />, tabBarLabel: 'Live' }} />
      <Tab.Screen name="Fixtures" component={FixturesScreen} options={{ tabBarLabel: 'Fixtures' }} />
      <Tab.Screen name="Standings" component={StandingsScreen} options={{ tabBarLabel: 'Table' }} />
      <Tab.Screen name="Profile" component={ProfileScreen} options={{ tabBarLabel: 'Profile' }} />
    </Tab.Navigator>
  );
}

export function RootNavigator() {
  const { user } = useAuth();
  return (
    <NavigationContainer theme={navTheme}>
      <Stack.Navigator screenOptions={{ headerShown: false, contentStyle: { backgroundColor: colors.background } }}>
        <Stack.Screen name="Main" component={Tabs} />
        <Stack.Screen name="MatchDetail" component={MatchDetailScreen} />
        <Stack.Screen name="Login" component={LoginScreen} options={{ presentation: 'modal' }} />
        <Stack.Screen name="Register" component={RegisterScreen} options={{ presentation: 'modal' }} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
