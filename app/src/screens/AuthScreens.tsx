import React, { useState } from 'react';
import { View, Text, StyleSheet, TextInput, Pressable, ActivityIndicator, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, radius } from '../theme/colors';
import { useAuth } from '../context/AuthContext';

function Input({ placeholder, value, onChangeText, secureTextEntry, autoCapitalize, keyboardType }: any) {
  return (
    <TextInput
      placeholder={placeholder}
      placeholderTextColor={colors.textMuted}
      value={value}
      onChangeText={onChangeText}
      secureTextEntry={secureTextEntry}
      autoCapitalize={autoCapitalize ?? 'none'}
      keyboardType={keyboardType}
      style={styles.input}
    />
  );
}

export function LoginScreen({ navigation }: any) {
  const { login } = useAuth();
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const onLogin = async () => {
    if (!identifier || !password) return Alert.alert('Error', 'Fill in all fields');
    setLoading(true);
    try { await login(identifier, password); } catch (e: any) { Alert.alert('Login failed', e.message); }
    setLoading(false);
  };

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <View style={styles.content}>
        <Text style={styles.title}>Welcome back</Text>
        <Text style={styles.subtitle}>Sign in to Sportz44</Text>
        <Input placeholder="Email or username" value={identifier} onChangeText={setIdentifier} keyboardType="email-address" />
        <Input placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry />
        <Pressable style={styles.btn} onPress={onLogin} disabled={loading}>
          {loading ? <ActivityIndicator color={colors.textOnPrimary} /> : <Text style={styles.btnText}>Sign In</Text>}
        </Pressable>
        <Pressable onPress={() => navigation.navigate('Register')} style={styles.link}>
          <Text style={styles.linkText}>Don't have an account? <Text style={styles.linkAccent}>Sign up</Text></Text>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

export function RegisterScreen({ navigation }: any) {
  const { register } = useAuth();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const onRegister = async () => {
    if (!username || !email || !password) return Alert.alert('Error', 'Fill in all fields');
    setLoading(true);
    try { await register(username, email, password); } catch (e: any) { Alert.alert('Registration failed', e.message); }
    setLoading(false);
  };

  return (
    <SafeAreaView style={styles.safe} edges={['top']}>
      <View style={styles.content}>
        <Text style={styles.title}>Create account</Text>
        <Text style={styles.subtitle}>Join Sportz44</Text>
        <Input placeholder="Username" value={username} onChangeText={setUsername} />
        <Input placeholder="Email" value={email} onChangeText={setEmail} keyboardType="email-address" />
        <Input placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry />
        <Pressable style={styles.btn} onPress={onRegister} disabled={loading}>
          {loading ? <ActivityIndicator color={colors.textOnPrimary} /> : <Text style={styles.btnText}>Create Account</Text>}
        </Pressable>
        <Pressable onPress={() => navigation.navigate('Login')} style={styles.link}>
          <Text style={styles.linkText}>Already have an account? <Text style={styles.linkAccent}>Sign in</Text></Text>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: colors.background },
  content: { flex: 1, padding: spacing.lg, justifyContent: 'center' },
  title: { color: colors.text, fontSize: 26, fontWeight: '800' },
  subtitle: { color: colors.textMuted, fontSize: 13, marginTop: 4, marginBottom: 24 },
  input: {
    backgroundColor: colors.surface, borderWidth: 1, borderColor: colors.border,
    borderRadius: radius.md, paddingHorizontal: 14, paddingVertical: 13,
    color: colors.text, fontSize: 14, marginBottom: 12,
  },
  btn: { backgroundColor: colors.primary, borderRadius: radius.md, paddingVertical: 14, alignItems: 'center', marginTop: 8 },
  btnText: { color: colors.textOnPrimary, fontSize: 14, fontWeight: '700' },
  link: { marginTop: 16, alignItems: 'center' },
  linkText: { color: colors.textMuted, fontSize: 13 },
  linkAccent: { color: colors.primary, fontWeight: '600' },
});
