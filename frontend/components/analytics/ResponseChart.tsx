'use client';

import { AreaChart, Area, ResponsiveContainer } from 'recharts';

const mockData = [
  { value: 65 },
  { value: 95 },
  { value: 85 },
  { value: 75 },
  { value: 90 },
  { value: 70 },
  { value: 85 },
  { value: 95 },
  { value: 60 },
  { value: 100 },
  { value: 80 },
  { value: 90 },
  { value: 70 },
  { value: 85 },
  { value: 95 },
];

export default function ResponseChart() {
  return (
    <div className="mt-6">
      <h3 className="text-[#0B2B4A]/80 text-base font-medium mb-4">
        Динамика ответов
      </h3>
      <ResponsiveContainer width="100%" height={120}>
        <AreaChart data={mockData}>
          <defs>
            <linearGradient id="colorGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#00B39F" stopOpacity={0.3}/>
              <stop offset="100%" stopColor="#00B39F" stopOpacity={0}/>
            </linearGradient>
          </defs>
          <Area 
            type="monotone" 
            dataKey="value" 
            stroke="#00B39F" 
            strokeWidth={3}
            fill="url(#colorGradient)"
            strokeLinecap="round"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}