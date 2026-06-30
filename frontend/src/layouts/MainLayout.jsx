import React from 'react';
import { Outlet } from 'react-router-dom';
import Navbar from '../components/Navbar';
import { Toaster } from 'react-hot-toast';

const MainLayout = () => {
  return (
    <div className="flex min-h-screen bg-black text-white font-sans overflow-hidden">
      <Navbar />
      <div className="flex-1 ml-64 relative w-full h-screen overflow-y-auto overflow-x-hidden bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-gray-900 via-black to-black">
        <Toaster position="top-right" toastOptions={{ style: { background: '#1f2937', color: '#fff' } }} />
        <Outlet />
      </div>
    </div>
  );
};

export default MainLayout;
