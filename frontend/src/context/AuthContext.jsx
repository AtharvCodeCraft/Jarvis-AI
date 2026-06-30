import React, { createContext, useState, useEffect } from 'react';
import api from '../services/api';
import toast from 'react-hot-toast';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [token, setToken] = useState(localStorage.getItem('token') || null);

  useEffect(() => {
    const fetchUser = async () => {
      if (token) {
        try {
          const res = await api.get('/users/me');
          setUser(res.data);
        } catch (error) {
          console.error("Failed to fetch user", error);
          logout();
        }
      }
      setLoading(false);
    };
    fetchUser();
  }, [token]);

  const login = async (username, password) => {
    try {
      const formData = new URLSearchParams();
      formData.append('username', username);
      formData.append('password', password);
      
      const res = await api.post('/auth/login', formData, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });
      setToken(res.data.access_token);
      localStorage.setItem('token', res.data.access_token);
      setUser(res.data.user);
      toast.success('Logged in successfully');
      return true;
    } catch (error) {
      let errMsg = 'Login failed';
      const detail = error.response?.data?.detail;
      if (typeof detail === 'string') {
        errMsg = detail;
      } else if (Array.isArray(detail) && detail.length > 0) {
        errMsg = detail[0].msg || errMsg;
      }
      toast.error(errMsg);
      return false;
    }
  };

  const register = async (userData) => {
    try {
      await api.post('/auth/register', userData);
      toast.success('Registration successful. Please login.');
      return true;
    } catch (error) {
      let errMsg = 'Registration failed';
      const detail = error.response?.data?.detail;
      if (typeof detail === 'string') {
        errMsg = detail;
      } else if (Array.isArray(detail) && detail.length > 0) {
        errMsg = detail[0].msg || errMsg;
      }
      toast.error(errMsg);
      return false;
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    toast.success('Logged out');
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
