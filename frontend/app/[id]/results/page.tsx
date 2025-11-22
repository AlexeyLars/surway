import { getPoll } from '@/app/services/api';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import PieChartComponent, { PieDataItem } from '@/components/analytics/PieChartComponent';

interface Result {
  text: string;
  votes: number;
  percentage: number;
}

const ResultBar = ({ percentage, text, votes }: Result) => (
  <div>
    <div className="flex justify-between items-center mb-1">
      <p className="font-semibold text-gray-800 break-all">{text}</p>
      <p className="text-sm font-bold text-gray-600 shrink-0 ml-4">{votes} голосов ({percentage.toFixed(1)}%)</p>
    </div>
    <div className="w-full bg-gray-200 rounded-full h-2.5">
      <div 
        className="bg-[#00B39F] h-2.5 rounded-full transition-all duration-500 ease-out" 
        style={{ width: `${percentage}%` }}
      ></div>
    </div>
  </div>
);

export default async function PollResultsPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  
  let pollData;
  try {
    pollData = await getPoll(id);
  } catch {
    notFound();
  }

  const { poll, votes, total } = pollData;

  const results: Result[] = poll.options.map((optionText: string) => {
    const voteCount = votes[optionText] || 0;
    const percentage = total > 0 ? (voteCount / total) * 100 : 0;
    return { text: optionText, votes: voteCount, percentage };
  });

  const pieData: PieDataItem[] = results.map((result) => ({
    name: result.text,
    value: result.votes,
    percentage: result.percentage
  }));

  const topOption = results.reduce((max, current) => 
    current.votes > max.votes ? current : max
  , results[0] || { text: '-', votes: 0, percentage: 0 });

  const avgPercentage = results.length > 0 ? (100 / results.length).toFixed(1) : '0';

  return (
    <main className="min-h-screen bg-[#F4F7FB] py-12 md:py-20">
      <div className="px-4 max-w-6xl mx-auto">
        <div className="mb-8">
          <p className="text-lg font-semibold text-[#00B39F]">Результаты опроса</p>
          <h1 className="text-3xl md:text-4xl font-bold text-[#0B2B4A] mt-1 break-words">
            {poll.title}
          </h1>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
     
          <div className="lg:col-span-2 space-y-8">
            <div className="bg-white rounded-xl shadow-lg p-6 sm:p-8">
              <h2 className="text-2xl font-bold text-[#0B2B4A] mb-6">Распределение голосов</h2>
              <div className="space-y-6">
                {results.map((result, index) => (
                  <ResultBar 
                    key={index} 
                    text={result.text} 
                    votes={result.votes} 
                    percentage={result.percentage} 
                  />
                ))}
              </div>
            </div>

            {total > 0 && (
              <div className="bg-white rounded-xl shadow-lg p-6 sm:p-8">
                <h2 className="text-2xl font-bold text-[#0B2B4A] mb-6">Визуализация</h2>
                <PieChartComponent data={pieData} />
              </div>
            )}
          </div>

          <div className="space-y-6">
            <div className="bg-white rounded-xl shadow-lg p-6">
              <h3 className="text-sm font-semibold text-gray-500 mb-2">Лидирующий вариант</h3>
              <p className="text-2xl font-bold text-[#0B2B4A] mb-1 break-words">{topOption.text}</p>
              <p className="text-3xl font-bold text-[#00B39F]">{topOption.votes} голосов</p>
              <p className="text-sm text-gray-600 mt-1">({topOption.percentage.toFixed(1)}% от общего числа)</p>
            </div>

            <div className="bg-white rounded-xl shadow-lg p-6">
              <div className="space-y-4">
                <div>
                  <p className="text-sm font-semibold text-gray-500 mb-1">Всего голосов</p>
                  <p className="text-3xl font-bold text-[#0B2B4A]">{total}</p>
                </div>
                <div className="border-t pt-4">
                  <p className="text-sm font-semibold text-gray-500 mb-1">Вариантов ответа</p>
                  <p className="text-3xl font-bold text-[#0B2B4A]">{results.length}</p>
                </div>
                <div className="border-t pt-4">
                  <p className="text-sm font-semibold text-gray-500 mb-1">Средний процент</p>
                  <p className="text-3xl font-bold text-[#0B2B4A]">{avgPercentage}%</p>
                  <p className="text-xs text-gray-500 mt-1">на один вариант</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-xl shadow-lg p-6">
              <p className="text-sm font-semibold text-gray-500 mb-1">Создан</p>
              <p className="text-lg font-bold text-[#0B2B4A]">
                {new Date(poll.created_at).toLocaleDateString('ru-RU', {
                  day: '2-digit',
                  month: '2-digit',
                  year: 'numeric'
                })}
              </p>
              <p className="text-sm text-gray-600 mt-1">
                {new Date(poll.created_at).toLocaleTimeString('ru-RU', {
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </p>
            </div>
          </div>
        </div>

        <div className="mt-12 text-center">
            <Link href="/create">
                <button
                    type="button"
                    className="px-8 py-3 bg-[#00B39F] text-white font-bold rounded-lg shadow-lg shadow-[#00B39F]/30 hover:shadow-xl hover:bg-[#00B39F]/90 transition-all tracking-wide"
                >
                    Создать новый опрос
                </button>
            </Link>
        </div>
      </div>
    </main>
  );
}