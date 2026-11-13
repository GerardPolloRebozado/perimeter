import { useEffect, useState } from 'react';
import { SafeAreaView, Text } from 'react-native';

const PING_URL = process.env.EXPO_PUBLIC_PING_URL;

export default function App() {
  const [body, setBody] = useState('loading');

  useEffect(() => {
    if (!PING_URL) {
      setBody('EXPO_PUBLIC_PING_URL is required');
      return;
    }

    fetch(PING_URL)
      .then((response) => response.text())
      .then((text) => setBody(text))
      .catch((error) => {
        setBody(error instanceof Error ? error.message : 'request failed');
      });
  }, []);

  return (
    <SafeAreaView>
      <Text>{body}</Text>
    </SafeAreaView>
  );
}
