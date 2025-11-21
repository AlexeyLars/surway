'use client';

import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';

interface PieDataItem {
  name: string;
  value: number;
  percentage: number;
}

interface PieChartComponentProps {
  data: PieDataItem[];
}

// Типизация для пропсов кастомного лейбла Recharts
interface CustomizedLabelProps {
  cx: number;
  cy: number;
  midAngle: number;
  innerRadius: number;
  outerRadius: number;
  percent: number;
}

const COLORS = [
  '#00B39F', // Основной
  '#3b82f6', // Синий
  '#84cc16', // Лаймовый
  '#f59e0b', // Мягкий янтарный
  '#6366f1', // Индиго
  '#6b7280', // Серый
  '#8b5cf6', // Сдержанный фиолетовый
  '#0ea5e9', // Небесно-голубой
];

const RADIAN = Math.PI / 180;

// Убрали any, добавили типизацию пропсов
const renderCustomizedLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent }: CustomizedLabelProps) => {
  const radius = innerRadius + (outerRadius - innerRadius) * 0.7;
  const x = cx + radius * Math.cos(-midAngle * RADIAN);
  const y = cy + radius * Math.sin(-midAngle * RADIAN);

  if (percent * 100 < 5) {
    return null;
  }

  return (
    <text 
      x={x} 
      y={y} 
      fill="white" 
      textAnchor="middle" 
      dominantBaseline="central" 
      className="font-semibold text-sm"
    >
      {`${(percent * 100).toFixed(0)}%`}
    </text>
  );
};

export default function PieChartComponent({ data }: PieChartComponentProps) {
  return (
    <ResponsiveContainer width="100%" height={350}>
      <PieChart>
        <Pie
          data={data} // Убрали "as any", Recharts должен принять типизированный массив
          cx="50%"
          cy="50%"
          labelLine={false}
          label={renderCustomizedLabel}
          outerRadius={120}
          fill="#8884d8"
          dataKey="value"
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip 
          formatter={(value: number, name: string) => [`${value} голосов`, name]}
          contentStyle={{
            backgroundColor: 'white',
            border: '1px solid #e2e8f0',
            borderRadius: '0.5rem',
          }}
        />
        <Legend 
          iconType="circle"
          wrapperStyle={{ fontSize: '14px', paddingTop: '20px' }}
          formatter={(value: string) => {
            const maxLength = 25;
            return value.length > maxLength ? value.substring(0, maxLength) + '...' : value;
          }}
        />
      </PieChart>
    </ResponsiveContainer>
  );
}