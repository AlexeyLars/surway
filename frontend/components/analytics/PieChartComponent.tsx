'use client';

import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';

interface PieDataItem {
  name: string;
  value: number;
  percentage: number;
  // Индексная сигнатура для Recharts
  [key: string]: string | number;
}

interface PieChartComponentProps {
  data: PieDataItem[];
}

// Разрешаем number | string, чтобы удовлетворить кривые типы Recharts
interface CustomizedLabelProps {
  cx?: number | string;
  cy?: number | string;
  midAngle?: number | string;
  innerRadius?: number | string;
  outerRadius?: number | string;
  percent?: number;
}

const COLORS = [
  '#00B39F',
  '#3b82f6',
  '#84cc16',
  '#f59e0b',
  '#6366f1',
  '#6b7280',
  '#8b5cf6',
  '#0ea5e9',
];

const RADIAN = Math.PI / 180;

const renderCustomizedLabel = (props: CustomizedLabelProps) => {
  // Безопасное извлечение и преобразование в числа
  // Recharts иногда определяет типы как string, хотя шлет числа, поэтому используем Number()
  const cx = Number(props.cx) || 0;
  const cy = Number(props.cy) || 0;
  const midAngle = Number(props.midAngle) || 0;
  const innerRadius = Number(props.innerRadius) || 0;
  const outerRadius = Number(props.outerRadius) || 0;
  const percent = Number(props.percent) || 0;

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
          data={data}
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
}//