import React, { useState } from 'react';
import { View, Text, StyleSheet, TextInput, Pressable, ActivityIndicator, Alert, ScrollView, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';
import { useAuth } from '../context/AuthContext';

const TEAL = '#0EA5B5';
const BG = '#121212';
const INPUT_BG = '#1E1E1E';
const INPUT_BORDER = '#2E2E2E';
const MUTED = '#6B7280';
const PLACEHOLDER = '#4B5563';

function Field({
  iconName,
  placeholder,
  value,
  onChangeText,
  secureTextEntry,
  keyboardType,
  autoCapitalize,
}: {
  iconName: keyof typeof Ionicons.glyphMap;
  placeholder: string;
  value: string;
  onChangeText: (v: string) => void;
  secureTextEntry?: boolean;
  keyboardType?: any;
  autoCapitalize?: any;
}) {
  const [hidden, setHidden] = useState(!!secureTextEntry);
  return (
    <View style={s.inputWrap}>
      <Ionicons name={iconName} size={16} color={MUTED} style={s.inputIcon} />
      <TextInput
        placeholder={placeholder}
        placeholderTextColor={PLACEHOLDER}
        value={value}
        onChangeText={onChangeText}
        secureTextEntry={hidden}
        autoCapitalize={autoCapitalize ?? 'none'}
        keyboardType={keyboardType}
        style={s.input}
      />
      {secureTextEntry && (
        <Pressable onPress={() => setHidden((h) => !h)} hitSlop={8}>
          <Ionicons name={hidden ? 'eye-off-outline' : 'eye-outline'} size={18} color={MUTED} />
        </Pressable>
      )}
    </View>
  );
}

function GoogleButton({ label }: { label: string }) {
  return (
    <Pressable style={s.googleBtn} onPress={() => Alert.alert('Coming soon', 'Google sign-in will be available soon.')}>
      <View style={s.googleCircle}>
        <Image source={require('../../assets/google.png')} style={s.googleImg} resizeMode="contain" />
      </View>
      <Text style={s.googleLabel}>{label}</Text>
    </Pressable>
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
    try {
      await login(identifier, password);
    } catch (e: any) {
      Alert.alert('Login failed', e.message);
    }
    setLoading(false);
  };

  return (
    <SafeAreaView style={s.safe} edges={['top', 'bottom']}>
      <View style={s.heroWrap}>
        <Image source={require('../../assets/onboarding/slide4.jpg')} style={s.heroImage} resizeMode="cover" />
        <View style={s.heroOverlay} />
        <View style={s.heroTextWrap}>
          <Text style={s.heroTitle}>Welcome back,</Text>
          <Text style={s.heroSubtitle}>Please login to enjoy full feature</Text>
        </View>
      </View>
      <ScrollView contentContainerStyle={s.content} keyboardShouldPersistTaps="handled" style={s.scroll}>
        <Field iconName="person-outline" placeholder="Username or Email" value={identifier} onChangeText={setIdentifier} keyboardType="email-address" />
        <Field iconName="lock-closed-outline" placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry />

        <Pressable style={s.forgotWrap} onPress={() => Alert.alert('Coming soon', 'Password reset coming soon.')}>
          <Text style={s.forgot}>Forgot password</Text>
        </Pressable>

        <Pressable style={s.cta} onPress={onLogin} disabled={loading}>
          {loading ? <ActivityIndicator color="#fff" /> : <Text style={s.ctaText}>Login</Text>}
        </Pressable>

        <Text style={s.divider}>Or login with</Text>
        <View style={s.socialRow}>
          <GoogleButton label="Google" />
        </View>
      </ScrollView>
      <View style={s.footer}>
        <Pressable onPress={() => navigation.navigate('Register')}>
          <Text style={s.bottomText}>
            Not have an account? <Text style={s.bottomAccent}>Register now</Text>
          </Text>
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
  const [confirm, setConfirm] = useState('');
  const [loading, setLoading] = useState(false);

  const onRegister = async () => {
    if (!username || !email || !password || !confirm) return Alert.alert('Error', 'Fill in all fields');
    if (password !== confirm) return Alert.alert('Error', 'Passwords do not match');
    if (password.length < 8) return Alert.alert('Error', 'Password must be at least 8 characters');
    setLoading(true);
    try {
      await register(username, email, password);
    } catch (e: any) {
      Alert.alert('Registration failed', e.message);
    }
    setLoading(false);
  };

  return (
    <SafeAreaView style={s.safe} edges={['top', 'bottom']}>
      <View style={s.heroWrap}>
        <Image source={require('../../assets/onboarding/slide5.jpg')} style={s.heroImage} resizeMode="cover" />
        <View style={s.heroOverlay} />
        <View style={s.heroTextWrap}>
          <Text style={s.heroTitle}>Welcome to Sportz44</Text>
          <Text style={s.heroSubtitle}>Create an account to explore amazing feature</Text>
        </View>
      </View>
      <ScrollView contentContainerStyle={s.content} keyboardShouldPersistTaps="handled" style={s.scroll}>
        <Field iconName="person-outline" placeholder="Username or Email" value={username} onChangeText={setUsername} />
        <Field iconName="mail-outline" placeholder="Email" value={email} onChangeText={setEmail} keyboardType="email-address" />
        <Field iconName="lock-closed-outline" placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry />
        <Field iconName="lock-closed-outline" placeholder="Confirm Password" value={confirm} onChangeText={setConfirm} secureTextEntry />

        <Pressable style={[s.cta, { marginTop: 20 }]} onPress={onRegister} disabled={loading}>
          {loading ? <ActivityIndicator color="#fff" /> : <Text style={s.ctaText}>Register</Text>}
        </Pressable>

        <Text style={s.divider}>Or register with</Text>
        <View style={s.socialRow}>
          <GoogleButton label="Google" />
        </View>
      </ScrollView>
      <View style={s.footer}>
        <Pressable onPress={() => navigation.navigate('Login')}>
          <Text style={s.bottomText}>
            Have an account? <Text style={s.bottomAccent}>Login</Text>
          </Text>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

const s = StyleSheet.create({
  safe: { flex: 1, backgroundColor: BG },
  heroWrap: { height: 180, overflow: 'hidden', position: 'relative' },
  heroImage: { width: '100%', height: '100%' },
  heroOverlay: { ...StyleSheet.absoluteFill, backgroundColor: 'rgba(0,0,0,0.55)' },
  heroTextWrap: { position: 'absolute', bottom: 20, left: 24, right: 24 },
  heroTitle: { color: '#fff', fontSize: 22, fontWeight: '800' },
  heroSubtitle: { color: 'rgba(255,255,255,0.7)', fontSize: 12, marginTop: 4 },
  scroll: { flex: 1 },
  content: { padding: 24, paddingTop: 20, paddingBottom: 16, flexGrow: 1 },
  title: { color: '#fff', fontSize: 20, fontWeight: '700' },
  subtitle: { color: MUTED, fontSize: 12, marginTop: 6, marginBottom: 24 },
  inputWrap: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: INPUT_BG,
    borderWidth: 1,
    borderColor: INPUT_BORDER,
    borderRadius: 12,
    paddingHorizontal: 12,
    paddingVertical: 2,
    marginBottom: 12,
  },
  inputIcon: { marginRight: 8 },
  input: { flex: 1, color: '#fff', fontSize: 13, paddingVertical: 12 },
  forgotWrap: { alignSelf: 'flex-end', marginTop: 2, marginBottom: 20 },
  forgot: { color: MUTED, fontSize: 11 },
  cta: {
    backgroundColor: TEAL,
    borderRadius: 24,
    paddingVertical: 14,
    alignItems: 'center',
  },
  ctaText: { color: '#fff', fontSize: 13, fontWeight: '700' },
  divider: { color: MUTED, fontSize: 11, textAlign: 'center', marginTop: 20, marginBottom: 14 },
  socialRow: { flexDirection: 'row', justifyContent: 'center', gap: 16 },
  googleBtn: { alignItems: 'center', gap: 6 },
  googleCircle: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
  googleImg: { width: 22, height: 22 }, // kept for future asset-based icon

  googleLabel: { color: MUTED, fontSize: 10 },
  footer: {
    paddingVertical: 16,
    paddingBottom: 8,
    alignItems: 'center',
    borderTopWidth: 1,
    borderTopColor: 'transparent',
  },
  bottomLink: { marginTop: 24, alignItems: 'center' },
  bottomText: { color: MUTED, fontSize: 12, textAlign: 'center' },
  bottomAccent: { color: TEAL, fontWeight: '600' },
});
