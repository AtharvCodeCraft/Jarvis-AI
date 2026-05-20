import React, { useEffect, useRef, useState } from 'react';

export default function CommandConsole({ messages, wsStatus }) {
  const endRef = useRef(null);
  const [searchQuery, setSearchQuery] = useState("");

  const filteredMessages = messages.filter(msg => 
    msg.text.toLowerCase().includes(searchQuery.toLowerCase())
  );

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <div className="flex-1 glass-panel p-4 flex flex-col h-[500px]">
      <div className="flex justify-between items-center mb-4 border-b border-jarvis-core/30 pb-2">
        <h2 className="text-xl font-semibold tracking-wider text-jarvis-core">COMMAND CONSOLE</h2>
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${wsStatus === 'Connected' ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></div>
          <span className="text-xs uppercase text-gray-400">{wsStatus}</span>
        </div>
      </div>

      <div className="mb-3 px-1">
        <input 
          type="text" 
          placeholder="SEARCH CONSOLE HISTORY..." 
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full bg-black/40 border-b border-jarvis-core/20 px-3 py-1.5 text-xs text-jarvis-core placeholder-jarvis-core/30 focus:outline-none focus:border-jarvis-core/60 transition-colors uppercase tracking-widest"
        />
      </div>
      
      <div className="flex-1 overflow-y-auto space-y-3 font-mono text-sm pr-2">
        <div className="text-jarvis-glow tracking-widest text-xs">J.A.R.V.I.S Initialization Sequence Completed.</div>
        
        {filteredMessages.length === 0 && searchQuery && (
          <div className="text-jarvis-core/50 text-xs text-center mt-4 uppercase tracking-widest">No matching commands found.</div>
        )}

        {filteredMessages.map((msg, idx) => (
          <div key={idx} className={`flex flex-col ${msg.sender === 'User' ? 'items-end' : 'items-start'}`}>
            <span className="text-[10px] uppercase text-gray-500 mb-1">{msg.sender}</span>
            <div className={`px-4 py-2 rounded max-w-[80%] ${msg.sender === 'User' ? 'bg-jarvis-core/20 text-white' : 'bg-black/50 border border-jarvis-glow/40 text-jarvis-core'}`}>
              {msg.text}
            </div>
          </div>
        ))}
        <div ref={endRef} />
      </div>
    </div>
  );
}
