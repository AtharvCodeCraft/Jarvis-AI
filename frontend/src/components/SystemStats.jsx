import React from 'react';

export default function SystemStats() {
  return (
    <div className="glass-panel p-4 flex flex-col gap-4 max-w-sm">
      <h2 className="text-xl font-semibold tracking-wider text-jarvis-core border-b border-jarvis-core/30 pb-2 uppercase">System Telemetry</h2>
      
      <div className="space-y-4">
        <div>
          <div className="flex justify-between text-xs text-gray-400 mb-1 uppercase tracking-widest">
            <span>CPU Load</span>
            <span>12%</span>
          </div>
          <div className="w-full bg-black/40 h-2 rounded">
            <div className="bg-jarvis-core h-full rounded shadow-[0_0_10px_rgba(0,210,255,0.8)]" style={{ width: '12%' }}></div>
          </div>
        </div>
        
        <div>
          <div className="flex justify-between text-xs text-gray-400 mb-1 uppercase tracking-widest">
            <span>Memory Usage</span>
            <span>48%</span>
          </div>
          <div className="w-full bg-black/40 h-2 rounded">
            <div className="bg-jarvis-glow h-full rounded shadow-[0_0_10px_rgba(58,123,213,0.8)]" style={{ width: '48%' }}></div>
          </div>
        </div>
        
        <div>
          <div className="flex justify-between text-xs text-gray-400 mb-1 uppercase tracking-widest">
            <span>Network Link</span>
            <span className="text-green-400">Stable</span>
          </div>
          <div className="w-full bg-black/40 h-2 rounded flex overflow-hidden gap-1">
             {[...Array(20)].map((_, i) => (
                <div key={i} className={`h-full flex-1 ${i < 15 ? 'bg-green-500/50' : 'bg-gray-800'}`}></div>
             ))}
          </div>
        </div>
      </div>
      
      <div className="mt-auto border-t border-jarvis-core/20 pt-4">
        <ul className="text-xs text-gray-400 space-y-2 uppercase tracking-wide">
          <li className="flex justify-between"><span>Core Temp:</span><span className="text-jarvis-core">42°C</span></li>
          <li className="flex justify-between"><span>Active Modules:</span><span className="text-jarvis-core">10</span></li>
          <li className="flex justify-between"><span>Local AI:</span><span className="text-green-500">Standby</span></li>
        </ul>
      </div>
    </div>
  );
}
