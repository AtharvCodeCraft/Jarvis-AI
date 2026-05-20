import { useState, useCallback } from 'react';

export function useVoice(language = 'en-IN') {
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState('');
  const [voiceError, setVoiceError] = useState('');

  const startListening = useCallback(() => {
    setIsListening(true);
    setVoiceError('');

    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (SpeechRecognition) {
      const recognition = new SpeechRecognition();
      recognition.continuous = false;
      recognition.interimResults = false;
      recognition.lang = language;

      recognition.onstart = () => {
        setIsListening(true);
        setVoiceError('');
      };

      recognition.onresult = (event) => {
        const result = event.results[event.results.length - 1];
        const text = result[0].transcript.trim();
        if (text) {
          setTranscript(text);
        } else {
          setVoiceError('Could not detect voice command. Please speak clearly.');
        }
        setIsListening(false);
      };

      recognition.onerror = (e) => {
        console.warn('Speech recognition error:', e && e.error ? e.error : e);
        setVoiceError(`Speech recognition error: ${e.error || 'unknown'}`);
        setIsListening(false);
      };

      recognition.onend = () => {
        setIsListening(false);
        if (!transcript) {
          setVoiceError('Listening ended without a transcript. Try again.');
        }
      };

      recognition.start();
    } else {
      console.warn('Speech recognition not supported.');
      setVoiceError('Speech recognition not supported in this browser. Use Chrome or Edge.');
      setIsListening(false);
    }
  }, [language, transcript]);

  return { isListening, transcript, startListening, setTranscript, voiceError };
}
