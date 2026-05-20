import React, { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Sphere, MeshDistortMaterial } from '@react-three/drei';

function Rings({ state }) {
  const ref1 = useRef();
  const ref2 = useRef();

  useFrame(() => {
    const speedMultiplier = state === 'Thinking' ? 3 : state === 'Speaking' ? 2 : 1;
    if (ref1.current) ref1.current.rotation.z += 0.01 * speedMultiplier;
    if (ref2.current) {
      ref2.current.rotation.z -= 0.015 * speedMultiplier;
      ref2.current.rotation.x += 0.005 * speedMultiplier;
    }
  });

  return (
    <group>
      <mesh ref={ref1}>
        <torusGeometry args={[2.5, 0.05, 16, 100]} />
        <meshBasicMaterial color="#00d2ff" wireframe opacity={0.5} transparent />
      </mesh>
      <mesh ref={ref2} rotation={[Math.PI / 4, 0, 0]}>
        <torusGeometry args={[3, 0.02, 16, 100]} />
        <meshBasicMaterial color="#3a7bd5" opacity={0.6} transparent />
      </mesh>
    </group>
  );
}

export default function JarvisFace({ coreState }) {
  // Compute dynamic scale based on state to simulate the 3D breathing/pulsing effect on the 2D image
  let scale = 1;
  let filter = 'drop-shadow(0 0 15px rgba(0,210,255,0.5))';

  if (coreState === 'Listening') {
    scale = 1.1;
    filter = 'drop-shadow(0 0 25px rgba(0,255,204,0.8)) hue-rotate(-30deg)';
  } else if (coreState === 'Speaking') {
    scale = 1.05;
    filter = 'drop-shadow(0 0 30px rgba(255,0,60,0.8)) hue-rotate(150deg)';
  } else if (coreState === 'Thinking') {
    scale = 1.15;
    filter = 'drop-shadow(0 0 20px rgba(0,210,255,0.8)) contrast(1.2)';
  }

  return (
    <div className="w-80 h-80 rounded-full flex items-center justify-center relative shadow-[0_0_80px_rgba(0,210,255,0.2)]">
      
      {/* 3D Orbiting Rings */}
      <div className="absolute inset-0 z-10 pointer-events-none">
        <Canvas camera={{ position: [0, 0, 5] }}>
          <ambientLight intensity={1} />
          <directionalLight position={[10, 10, 5]} intensity={2} />
          <Rings state={coreState} />
        </Canvas>
      </div>

      {/* 2D Ironman Mask */}
      <div className="absolute inset-0 z-20 flex items-center justify-center pointer-events-none mix-blend-screen">
        <img 
          src="/ironman.png" 
          alt="Jarvis Mask Core" 
          className="w-56 h-56 object-cover rounded-full transition-all duration-300 ease-in-out"
          style={{ 
            transform: `scale(${scale})`,
            filter: filter
          }}
        />
      </div>
      
    </div>
  );
}
