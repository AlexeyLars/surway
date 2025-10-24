'use client';

import { useState, useEffect } from 'react';

interface ProgressRingProps {
  percentage: number;
}

export default function ProgressRing({ percentage }: ProgressRingProps) {
  const [animatedPercentage, setAnimatedPercentage] = useState(0);
  
  const radius = 105;
  const strokeWidth = 22;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (animatedPercentage / 100) * circumference;
  
  const size = 280;
  const center = size / 2;

  useEffect(() => {
    const timer = setTimeout(() => {
      setAnimatedPercentage(percentage);
    }, 100); 
    
    return () => clearTimeout(timer);
  }, [percentage]);

  return (
    <div className="flex items-center justify-center relative">
      <svg 
        width={size} 
        height={size} 
        style={{ transform: 'rotate(-90deg)' }}
      >

        <circle
          cx={center}
          cy={center}
          r={radius}
          fill="none"
          stroke="#e5e7eb"
          strokeWidth={strokeWidth}
        />
       
        <circle
          cx={center}
          cy={center}
          r={radius}
          fill="none"
          stroke="#00B39F"
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          className="transition-all duration-[1500ms] ease-out"
        />
      </svg>
     

      <div className="absolute text-center">
        <div className="text-5xl font-bold text-[#0B2B4A]">
          {animatedPercentage}%
        </div>
        <div className="text-[#0B2B4A]/70 mt-2">
          Завершено
        </div>
      </div>
    </div>
  );
}