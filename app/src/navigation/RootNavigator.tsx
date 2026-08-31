import React from 'react';
import { NavigationContainer, DefaultTheme } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Ionicons } from '@expo/vector-icons';
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
      <Tab.Screen
        name="Home"
        component={HomeScreen}
        options={{
          tabBarLabel: 'Home',
          tabBarIcon: ({ color, size }) => <Ionicons name="home-outline" size={20} color={color} />,
        }}
      />
      <Tab.Screen
        name="Live"
        component={LiveScreen}
        options={{
          tabBarLabel: 'Live',
          tabBarIcon: ({ color, size }) => <Ionicons name="radio-outline" size={20} color={color} />,
        }}
      />
      <Tab.Screen
        name="Fixtures"
        component={FixturesScreen}
        options={{
          tabBarLabel: 'Fixtures',
          tabBarIcon: ({ color, size }) => <Ionicons name="calendar-outline" size={20} color={color} />,
        }}
      />
      <Tab.Screen
        name="Standings"
        component={StandingsScreen}
        options={{
          tabBarLabel: 'Table',
          tabBarIcon: ({ color, size }) => <Ionicons name="trophy-outline" size={20} color={color} />,
        }}
      />
      <Tab.Screen
        name="Profile"
        component={ProfileScreen}
        options={{
          tabBarLabel: 'Profile',
          tabBarIcon: ({ color, size }) => <Ionicons name="person-outline" size={20} color={color} />,
        }}
      />
    </Tab.Navigator>
  );
}

function AuthStack() {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, contentStyle: { backgroundColor: colors.background } }}>
      <Stack.Screen name="Login" component={LoginScreen} />
      <Stack.Screen name="Register" component={RegisterScreen} />
    </Stack.Navigator>
  );
}

function AppStack() {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, contentStyle: { backgroundColor: colors.background } }}>
      <Stack.Screen name="Main" component={Tabs} />
      <Stack.Screen name="MatchDetail" component={MatchDetailScreen} />
    </Stack.Navigator>
  );
}

export function RootNavigator() {
  const { user, loading } = useAuth();

  if (loading) return null;

  return (
    <NavigationContainer theme={navTheme}>
      {user ? <AppStack /> : <AuthStack />}
    </NavigationContainer>
  );
}
