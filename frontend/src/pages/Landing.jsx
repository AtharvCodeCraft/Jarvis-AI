import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Zap, Shield, Cpu, Code } from 'lucide-react';

const Landing = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-8 text-center pt-20">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
        className="max-w-4xl"
      >
        <div className="inline-block px-4 py-1.5 rounded-full border border-blue-500/30 bg-blue-500/10 text-blue-400 font-medium text-sm mb-6">
          Jarvis AI V2.0 is here
        </div>
        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-br from-white via-blue-100 to-blue-400">
          Your Personal <br /> Artificial Intelligence
        </h1>
        <p className="text-xl text-gray-400 mb-10 max-w-2xl mx-auto">
          Automate your digital life, control your system, and have natural conversations with the most advanced AI assistant built for power users.
        </p>
        
        <div className="flex flex-wrap items-center justify-center gap-4 mb-20">
          <Link to="/register" className="px-8 py-4 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-semibold transition-all shadow-[0_0_20px_rgba(37,99,235,0.5)]">
            Get Started Now
          </Link>
          <Link to="/login" className="px-8 py-4 rounded-xl bg-gray-800 hover:bg-gray-700 border border-gray-700 text-white font-semibold transition-all">
            Sign In
          </Link>
        </div>
      </motion.div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl w-full">
        {[
          { icon: <Zap className="text-yellow-400" size={24} />, title: "Lightning Fast", desc: "Local execution for zero latency commands." },
          { icon: <Shield className="text-green-400" size={24} />, title: "Secure & Private", desc: "Your data stays isolated. JWT protected." },
          { icon: <Cpu className="text-blue-400" size={24} />, title: "System Control", desc: "Automate windows, apps, and browsers." },
          { icon: <Code className="text-purple-400" size={24} />, title: "Developer Ready", desc: "Extensible plugin architecture." },
        ].map((feature, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: i * 0.1 }}
            className="bg-gray-900/40 border border-gray-800 rounded-2xl p-6 hover:bg-gray-800/50 transition-colors"
          >
            <div className="w-12 h-12 bg-gray-800 rounded-xl flex items-center justify-center mb-4 border border-gray-700 shadow-inner">
              {feature.icon}
            </div>
            <h3 className="text-xl font-bold mb-2">{feature.title}</h3>
            <p className="text-gray-400 text-sm">{feature.desc}</p>
          </motion.div>
        ))}
      </div>
    </div>
  );
};

export default Landing;
