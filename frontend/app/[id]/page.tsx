'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { getPoll, voteOnPoll } from '@/app/services/api';
import { Loader2 } from 'lucide-react';

interface Poll {
  id: string;
  title: string;
  options: string[];
}

interface PollData {
  poll: Poll;
}

export default function PollVotePage() {
  const router = useRouter();
  const params = useParams();
  const pollId = params.id as string;

  const [pollData, setPollData] = useState<PollData | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<number[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isVoting, setIsVoting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!pollId) return;
    const fetchPoll = async () => {
      try {
        const data = await getPoll(pollId);
        setPollData(data);
      } catch (err) {
        setError('Опрос не найден или произошла ошибка при его загрузке.');
      } finally {
        setIsLoading(false);
      }
    };
    fetchPoll();
  }, [pollId]);

  const handleOptionToggle = (index: number) => {
    setSelectedOptions((prev) =>
      prev.includes(index)
        ? prev.filter((item) => item !== index)
        : [...prev, index]
    );
  };

  const handleVote = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedOptions.length === 0) {
      setError('Пожалуйста, выберите хотя бы один вариант.');
      return;
    }

    setIsVoting(true);
    setError(null);

    try {
      await voteOnPoll(pollId, selectedOptions);
      router.push(`/${pollId}/results`);
    } catch (err) {
      setError('Не удалось проголосовать. Пожалуйста, попробуйте снова.');
      setIsVoting(false);
    }
  };

  if (isLoading) {
    return (
      <main className="min-h-screen bg-[#F4F7FB] flex items-center justify-center">
        <Loader2 className="animate-spin text-[#00B39F]" size={48} />
      </main>
    );
  }

  if (error && !pollData) {
    return (
      <main className="min-h-screen bg-[#F4F7FB] flex items-center justify-center text-center px-4">
        <div>
          <h1 className="text-2xl font-bold text-red-500">Ошибка</h1>
          <p className="text-gray-600 mt-2">{error}</p>
        </div>
      </main>
    );
  }

  if (!pollData) {
    return null;
  }

  return (
    <main className="min-h-screen bg-[#F4F7FB] py-12 md:py-20">
      <div className="px-4">
        <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-lg p-6 sm:p-8 md:p-10">
          <form onSubmit={handleVote}>
            <div className="mb-8">
              <p className="text-sm font-semibold text-gray-500">Голосование (можно выбрать несколько)</p>
              <h1 className="text-2xl md:text-3xl font-bold text-[#0B2B4A] mt-1">
                {pollData.poll.title}
              </h1>
            </div>
            <div className="space-y-4">
              {pollData.poll.options.map((option, index) => (
                <label
                  key={index}
                  className={`flex items-center p-4 border rounded-lg cursor-pointer transition-all ${
                    selectedOptions.includes(index)
                      ? 'border-[#00B39F] bg-teal-50 ring-2 ring-[#00B39F]'
                      : 'border-gray-300 bg-white hover:bg-gray-50'
                  }`}
                >
                  <input
                    type="checkbox"
                    checked={selectedOptions.includes(index)}
                    onChange={() => handleOptionToggle(index)}
                    className="h-5 w-5 text-[#00B39F] focus:ring-[#00B39F] border-gray-300 rounded"
                  />
                  <span className="ml-4 text-lg text-gray-800 font-medium">{option}</span>
                </label>
              ))}
            </div>
            {error && <p className="mt-6 text-center text-red-500">{error}</p>}
            <div className="mt-10">
              <button
                type="submit"
                disabled={isVoting || selectedOptions.length === 0}
                className="w-full flex items-center justify-center gap-2 px-6 py-4 bg-[#00B39F] text-white text-lg font-bold rounded-lg shadow-lg shadow-[#00B39F]/30 hover:shadow-xl hover:bg-[#00B39F]/90 transition-all tracking-wide disabled:bg-gray-400 disabled:shadow-none disabled:cursor-not-allowed"
              >
                {isVoting ? (
                  <>
                    <Loader2 className="animate-spin" size={24} />
                    Голосование...
                  </>
                ) : (
                  'Проголосовать'
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
}