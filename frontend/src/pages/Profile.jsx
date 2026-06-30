import React, { useContext, useState } from 'react';
import { AuthContext } from '../context/AuthContext';
import api from '../services/api';
import toast from 'react-hot-toast';

const Profile = () => {
  const { user } = useContext(AuthContext);
  const [formData, setFormData] = useState({
    full_name: user?.full_name || '',
    username: user?.username || '',
  });

  const handleChange = (e) => setFormData({ ...formData, [e.target.name]: e.target.value });

  const handleSave = async () => {
    try {
      await api.put('/users/update', formData);
      toast.success('Profile updated successfully');
    } catch (error) {
      toast.error('Failed to update profile');
    }
  };

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Profile</h1>
      
      <div className="bg-gray-900/60 border border-gray-800 rounded-2xl p-8 backdrop-blur-md">
        <div className="flex items-center gap-6 mb-8">
          <div className="w-24 h-24 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-3xl font-bold border-4 border-gray-800">
            {user?.username?.[0]?.toUpperCase() || 'U'}
          </div>
          <div>
            <h2 className="text-2xl font-bold">{user?.full_name || user?.username}</h2>
            <p className="text-gray-400">{user?.email}</p>
            <div className="mt-2 inline-block px-3 py-1 bg-blue-500/20 text-blue-400 rounded-full text-xs font-semibold">
              {user?.role.toUpperCase()}
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">Full Name</label>
              <input 
                type="text" 
                name="full_name"
                value={formData.full_name}
                onChange={handleChange}
                className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-blue-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">Username</label>
              <input 
                type="text" 
                name="username"
                value={formData.username}
                onChange={handleChange}
                className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-blue-500 focus:outline-none"
              />
            </div>
          </div>
          
          <button 
            onClick={handleSave}
            className="px-6 py-3 bg-blue-600 hover:bg-blue-500 rounded-xl font-medium transition-colors"
          >
            Save Changes
          </button>
        </div>
      </div>
    </div>
  );
};

export default Profile;
