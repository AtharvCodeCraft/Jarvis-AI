import React, { useContext } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { Home, LayoutDashboard, MessageSquare, Settings, User, LogOut, Cpu, LogIn, UserPlus } from 'lucide-react';
import { motion } from 'framer-motion';

const Navbar = () => {
  const { user, logout } = useContext(AuthContext);
  const location = useLocation();

  const loggedOutLinks = [
    { name: 'Home', path: '/', icon: <Home size={18} /> },
    { name: 'Login', path: '/login', icon: <LogIn size={18} /> },
    { name: 'Sign Up', path: '/register', icon: <UserPlus size={18} /> },
  ];

  const loggedInLinks = [
    { name: 'Dashboard', path: '/dashboard', icon: <LayoutDashboard size={18} /> },
    { name: 'Chat', path: '/chat', icon: <MessageSquare size={18} /> },
    { name: 'Automation', path: '/automation', icon: <Cpu size={18} /> },
    { name: 'Profile', path: '/profile', icon: <User size={18} /> },
    { name: 'Settings', path: '/settings', icon: <Settings size={18} /> },
  ];

  const links = user ? loggedInLinks : loggedOutLinks;

  return (
    <nav className="fixed top-0 left-0 h-full w-64 bg-gray-900/80 backdrop-blur-md border-r border-gray-800 p-4 flex flex-col justify-between z-50">
      <div>
        <div className="flex items-center gap-3 mb-10 mt-4 px-2">
          <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center font-bold text-white shadow-[0_0_15px_rgba(59,130,246,0.5)]">J</div>
          <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-cyan-300">Jarvis AI</span>
        </div>
        
        <ul className="space-y-2">
          {links.map((link) => {
            const isActive = location.pathname === link.path;
            return (
              <li key={link.name}>
                <Link to={link.path}>
                  <div className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-300 ${isActive ? 'bg-blue-600/20 text-blue-400 border border-blue-500/30' : 'text-gray-400 hover:bg-gray-800 hover:text-gray-200'}`}>
                    {link.icon}
                    <span className="font-medium">{link.name}</span>
                    {isActive && (
                      <motion.div
                        layoutId="active-nav"
                        className="absolute left-0 w-1 h-8 bg-blue-500 rounded-r-full"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ duration: 0.3 }}
                      />
                    )}
                  </div>
                </Link>
              </li>
            );
          })}
        </ul>
      </div>

      {user && (
        <div className="border-t border-gray-800 pt-4">
          <button onClick={logout} className="w-full flex items-center gap-3 px-4 py-3 text-red-400 hover:bg-red-500/10 rounded-xl transition-colors">
            <LogOut size={18} />
            <span className="font-medium">Logout</span>
          </button>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
