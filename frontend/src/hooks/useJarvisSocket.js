import { useState, useEffect, useRef } from 'react';

// Map BCP-47 lang codes to preferred TTS voice name fragments
const VOICE_PREFERENCES = {
  'en-IN': ['Microsoft Heera', 'Microsoft Ravi', 'Google UK English', 'Microsoft Mark'],
  'hi-IN': ['Microsoft Hemant', 'Microsoft Kalpana', 'Google हिन्दी', 'hi-IN'],
  'mr-IN': ['Microsoft Swara', 'mr-IN', 'hi-IN'],
};

function pickVoice(lang) {
  const voices = window.speechSynthesis.getVoices();
  const prefs = VOICE_PREFERENCES[lang] || VOICE_PREFERENCES['en-IN'];
  for (const pref of prefs) {
    const found = voices.find(v =>
      v.name.includes(pref) || v.lang === pref
    );
    if (found) return found;
  }
  // Fall back to any voice matching the lang prefix (e.g. "hi")
  const prefix = lang.split('-')[0];
  return voices.find(v => v.lang.startsWith(prefix)) || null;
}

export function useJarvisSocket(url) {
  const [messages, setMessages] = useState([]);
  const [status, setStatus] = useState('Disconnected');
  const [coreState, setCoreState] = useState('Idle');
  const socketRef = useRef(null);

  useEffect(() => {
    let reconnectAttempts = 0;
    let reconnectTimer = null;
    let mounted = true;

    const getWebSocketUrl = () => {
      if (url && url.trim().length > 0) return url;
      const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
      return `${protocol}://${window.location.hostname}:8080/ws`;
    };

    const connect = () => {
      if (!mounted) return;
      const wsUrl = getWebSocketUrl();
      setStatus('Connecting...');
      console.debug('Connecting WebSocket to', wsUrl);

      const socket = new WebSocket(wsUrl);
      socketRef.current = socket;

      socket.onopen = () => {
        reconnectAttempts = 0;
        setStatus('Connected');
        console.debug('WebSocket connected');
      };

      socket.onclose = () => {
        setStatus('Disconnected');
        console.warn('WebSocket disconnected');
        if (!mounted) return;

        const delay = Math.min(30000, 1000 * 2 ** reconnectAttempts);
        reconnectAttempts += 1;
        reconnectTimer = setTimeout(connect, delay);
      };

      socket.onerror = (err) => {
        console.error('WebSocket error', err);
        setStatus('Disconnected');
        socket.close();
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'command_result' || data.type === 'speech_output') {
            const lang = data.language || 'en-IN';
            setMessages(prev => [...prev, { sender: 'Jarvis', text: data.text, lang }]);

            if ('speechSynthesis' in window) {
              const utterance = new SpeechSynthesisUtterance(data.text);
              utterance.lang = lang;
              utterance.rate = 1.0;
              utterance.pitch = 1.0;

              const trySpeak = () => {
                const voice = pickVoice(lang);
                if (voice) utterance.voice = voice;
                window.speechSynthesis.speak(utterance);
              };

              if (window.speechSynthesis.getVoices().length > 0) {
                trySpeak();
              } else {
                window.speechSynthesis.addEventListener('voiceschanged', trySpeak, { once: true });
              }
            }

            setCoreState('Speaking');
            setTimeout(() => setCoreState('Idle'), 3000);
          }
        } catch (err) {
          console.error('Failed to parse WS message', err);
        }
      };
    };

    connect();

    return () => {
      mounted = false;
      if (reconnectTimer) clearTimeout(reconnectTimer);
      socketRef.current?.close();
    };
  }, [url]);

  const sendCommand = (text, language = 'en-IN') => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      setMessages(prev => [...prev, { sender: 'User', text, lang: language }]);
      setCoreState('Thinking');
      socketRef.current.send(JSON.stringify({ type: 'manual_command', text, language }));
    } else {
      setMessages(prev => [...prev, { sender: 'System', text: 'Error: WebSocket not connected. Start Go backend.', lang: 'en-IN' }]);
    }
  };

  return { messages, status, coreState, setCoreState, sendCommand };
}
