'use client';

interface ProgressRingProps {
  percentage: number;
}

export default function ProgressRing({ percentage }: ProgressRingProps) {
  const radius = 105;
  const strokeWidth = 22;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (percentage / 100) * circumference;
  
  const size = 280;
  const center = size / 2; 

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
          className="transition-all duration-500 ease-in-out"
        />
      </svg>
     
      <div className="absolute text-center">
        <div className="text-4xl font-bold text-slate-800">
          {percentage}%
        </div>
        <div className="text-slate-500 mt-2">
          Completed
        </div>
      </div>
    </div>
  );
}