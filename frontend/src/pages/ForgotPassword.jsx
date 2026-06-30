import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';
import toast from 'react-hot-toast';

const ForgotPassword = () => {
  const [email, setEmail] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await api.post('/auth/forgot-password', { email });
      toast.success('If this email is registered, a password reset link will be sent.');
    } catch (error) {
      toast.error('Failed to process request');
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <div className="w-full max-w-md bg-gray-900/50 backdrop-blur-xl border border-gray-800 rounded-2xl shadow-2xl overflow-hidden p-8">
        <h1 className="text-2xl font-bold mb-2">Reset Password</h1>
        <p className="text-gray-400 text-sm mb-6">Enter your email address and we'll send you a link to reset your password.</p>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-gray-800/50 border border-gray-700 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-blue-500"
              required
            />
          </div>
          <button type="submit" className="w-full bg-blue-600 hover:bg-blue-500 text-white font-medium py-2.5 rounded-xl transition-colors">
            Send Reset Link
          </button>
        </form>

        <p className="text-center text-gray-400 mt-6 text-sm">
          Remember your password? <Link to="/login" className="text-blue-400 hover:text-blue-300">Sign in</Link>
        </p>
      </div>
    </div>
  );
};

export default ForgotPassword;
