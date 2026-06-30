import React, { useContext, useState } from 'react';
import { AuthContext } from '../context/AuthContext';
import api from '../services/api';
import toast from 'react-hot-toast';
import { Trash2 } from 'lucide-react';

const Settings = () => {
  const { user, logout } = useContext(AuthContext);
  const [preferences, setPreferences] = useState({
    theme: user?.theme || 'dark',
    preferred_model: user?.preferred_model || 'default',
    preferred_voice: user?.preferred_voice || 'default',
  });

  const handleChange = (e) => setPreferences({ ...preferences, [e.target.name]: e.target.value });

  const handleSave = async () => {
    try {
      await api.put('/users/update', preferences);
      toast.success('Settings saved');
    } catch (error) {
      toast.error('Failed to save settings');
    }
  };

  const handleDeleteAccount = async () => {
    if (window.confirm("Are you sure you want to delete your account? This cannot be undone.")) {
      try {
        await api.delete('/users/delete');
        logout();
      } catch (error) {
        toast.error("Failed to delete account");
      }
    }
  };

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Settings</h1>
      
      <div className="space-y-6">
        <div className="bg-gray-900/60 border border-gray-800 rounded-2xl p-6">
          <h2 className="text-xl font-semibold mb-6">Preferences</h2>
          
          <div className="space-y-4 max-w-md">
            <div>
              <label className="block text-sm text-gray-400 mb-2">Theme</label>
              <select name="theme" value={preferences.theme} onChange={handleChange} className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-blue-500 focus:outline-none">
                <option value="dark">Dark Mode</option>
                <option value="light">Light Mode (Coming Soon)</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm text-gray-400 mb-2">AI Model</label>
              <select name="preferred_model" value={preferences.preferred_model} onChange={handleChange} className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-blue-500 focus:outline-none">
                <option value="default">Jarvis Default</option>
                <option value="llama3">Llama 3</option>
                <option value="mistral">Mistral</option>
              </select>
            </div>

            <div>
              <label className="block text-sm text-gray-400 mb-2">Voice Synthesis</label>
              <select name="preferred_voice" value={preferences.preferred_voice} onChange={handleChange} className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-blue-500 focus:outline-none">
                <option value="default">Default Male</option>
                <option value="female">Female Voice 1</option>
              </select>
            </div>

            <button onClick={handleSave} className="px-6 py-3 bg-blue-600 hover:bg-blue-500 rounded-xl font-medium mt-4">
              Save Preferences
            </button>
          </div>
        </div>

        <div className="bg-red-950/20 border border-red-900/50 rounded-2xl p-6">
          <h2 className="text-xl font-semibold text-red-500 mb-4 flex items-center gap-2"><Trash2 size={20} /> Danger Zone</h2>
          <p className="text-gray-400 mb-4 text-sm">Once you delete your account, there is no going back. Please be certain.</p>
          <button onClick={handleDeleteAccount} className="px-4 py-2 bg-red-600/20 hover:bg-red-600/40 text-red-500 rounded-lg font-medium border border-red-500/30 transition-colors">
            Delete Account
          </button>
        </div>
      </div>
    </div>
  );
};

export default Settings;
