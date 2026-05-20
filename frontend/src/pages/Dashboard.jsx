import React, { useEffect, useState } from 'react';
import JarvisFace from '../components/JarvisFace';
import CommandConsole from '../components/CommandConsole';
import SystemStats from '../components/SystemStats';
import MicIndicator from '../components/MicIndicator';
import { useJarvisSocket } from '../hooks/useJarvisSocket';
import { useVoice } from '../hooks/useVoice';

const LANGUAGES = [
  { code: 'en-IN', label: 'EN' },
  { code: 'hi-IN', label: 'हिंदी' },
  { code: 'mr-IN', label: 'मराठी' },
];

export default function Dashboard() {
  const [language, setLanguage] = useState('en-IN');

  const backendWsUrl =
    import.meta.env.VITE_BACKEND_WS ||
    `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.hostname}:8080/ws`;

  const { messages, status, coreState, sendCommand } = useJarvisSocket(backendWsUrl);
  const { isListening, transcript, startListening, setTranscript, voiceError } = useVoice(language);

  // Send voice transcript as command when recognition finishes
  useEffect(() => {
    if (transcript) {
      sendCommand(transcript, language);
      setTranscript(''); // reset
    }
  }, [transcript, sendCommand, setTranscript, language]);

  // Logging voice errors for debugging
  useEffect(() => {
    if (voiceError) {
      console.warn('Mic issue:', voiceError);
    }
  }, [voiceError]);

  const handleManualCommand = (e) => {
    e.preventDefault();
    const fd = new FormData(e.target);
    const cmd = fd.get('command');
    if (cmd.trim()) {
      sendCommand(cmd, language);
      e.target.reset();
    }
  };

  const placeholders = {
    'en-IN': 'Input manual override command...',
    'hi-IN': 'कमांड टाइप करें...',
    'mr-IN': 'कमांड टाइप करा...',
  };

  return (
    <div className="w-full max-w-7xl h-full p-6 grid grid-cols-1 md:grid-cols-3 gap-6 relative z-10 py-10">
      
      {/* Left Panel */}
      <div className="flex flex-col gap-6">
        <SystemStats />
      </div>

      {/* Center Panel: Face & Mic */}
      <div className="flex flex-col items-center justify-between relative py-12">
        <div className="text-center">
          <h1 className="text-5xl font-bold tracking-[0.2em] text-jarvis-core drop-shadow-[0_0_15px_rgba(0,210,255,0.8)]">
            RUDRA
          </h1>
          <p className="text-xs text-jarvis-core/60 mt-1 uppercase tracking-widest">Local AI Core</p>

          {/* ── Language Selector ── */}
          <div className="mt-4 flex items-center justify-center gap-1 bg-black/40 rounded-full px-2 py-1 border border-jarvis-core/20 w-fit mx-auto">
            {LANGUAGES.map(({ code, label }) => (
              <button
                key={code}
                onClick={() => setLanguage(code)}
                className={`px-3 py-1 rounded-full text-xs font-bold tracking-wide transition-all duration-200 ${
                  language === code
                    ? 'bg-jarvis-core text-black shadow-[0_0_10px_rgba(0,210,255,0.6)]'
                    : 'text-jarvis-core/60 hover:text-jarvis-core hover:bg-jarvis-core/10'
                }`}
              >
                {label}
              </button>
            ))}
          </div>
        </div>
        
        <div className="flex-1 flex items-center justify-center -mt-10">
          <JarvisFace coreState={isListening ? 'Listening' : coreState} />
        </div>
        
        <div className="mt-8">
          <MicIndicator isListening={isListening} onStart={startListening} />
          {voiceError && (
            <p className="mt-2 text-xs text-red-400 font-semibold text-center">
              {voiceError}
            </p>
          )}
        </div>
      </div>

      {/* Right Panel: Console */}
      <div className="flex flex-col h-[650px] pl-0 md:pl-6">
        <CommandConsole messages={messages} wsStatus={status} />
        
        <form onSubmit={handleManualCommand} className="mt-4 flex gap-2 w-full shadow-[0_4px_30px_rgba(0,0,0,0.5)] bg-jarvis-bg/90 rounded backdrop-blur-md p-1 border border-jarvis-core/20">
          <input 
            name="command" 
            autoComplete="off"
            placeholder={placeholders[language]} 
            className="flex-1 bg-transparent px-4 py-2 text-sm text-jarvis-core placeholder-jarvis-core/30 focus:outline-none placeholder:uppercase"
          />
          <button type="submit" className="bg-jarvis-core/20 hover:bg-jarvis-core/40 border border-jarvis-core/50 px-4 py-2 rounded text-jarvis-core transition-all uppercase text-xs tracking-wider font-bold">
            Execute
          </button>
        </form>
      </div>

      {/* Bottom Left Watermark */}
      <div className="absolute bottom-4 left-8 text-[10px] text-jarvis-core/40 tracking-widest uppercase pointer-events-none">
        Developed by Atharv
      </div>

      {/* Background glow override */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(circle_at_center,rgba(0,210,255,0.05)_0%,transparent_60%)] z-[-1]" />
    </div>
  );
}
