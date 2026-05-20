import React from 'react';

export default function MicIndicator({ isListening, onStart }) {
  return (
    <div className="flex flex-col items-center">
      <button 
        onClick={onStart}
        className={`w-16 h-16 rounded-full flex items-center justify-center transition-all duration-300 border-2 ${isListening ? 'bg-red-500/20 border-red-500 shadow-[0_0_30px_rgba(255,0,0,0.6)] animate-pulse' : 'bg-jarvis-core/10 border-jarvis-core/50 hover:bg-jarvis-core/30 hover:shadow-[0_0_20px_rgba(0,210,255,0.4)]'}`}
      >
        <svg xmlns="http://www.w3.org/2000/svg" className={`h-8 w-8 ${isListening ? 'text-red-500' : 'text-jarvis-core'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
        </svg>
      </button>
      <span className={`mt-3 text-xs tracking-[0.2em] uppercase font-bold ${isListening ? 'text-red-500' : 'text-jarvis-core/70'}`}>
        {isListening ? 'Voice Link Active' : 'Tap to Speak'}
      </span>
    </div>
  );
}
