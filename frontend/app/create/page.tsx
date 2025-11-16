'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { PlusCircle, XCircle } from 'lucide-react';
import { createPoll } from '@/app/services/api';

export default function CreatePollPage() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [options, setOptions] = useState(['', '']);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const addOption = () => {
    setOptions([...options, '']);
  };

  const removeOption = (index: number) => {
    if (options.length <= 2) return;
    setOptions(options.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const filteredOptions = options.map(opt => opt.trim()).filter(Boolean);
    const trimmedTitle = title.trim();

    if (trimmedTitle.length < 5 || filteredOptions.length < 2) {
      setError('Вопрос должен быть длиннее 4 символов и иметь минимум два варианта ответа.');
      return;
    }

    setIsLoading(true);

    try {
      const data = await createPoll(trimmedTitle, filteredOptions);
      router.push(`/poll/${data.poll_id}/success`);
    } catch (err) {
      setError('Произошла ошибка при создании опроса. Попробуйте снова.');
      setIsLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-[#F4F7FB] py-12 md:py-20">
      <div className="px-4">
        <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-lg p-6 sm:p-8 md:p-10">
          <div className="mb-8 text-center">
            <h1 className="text-3xl md:text-4xl font-bold text-[#0B2B4A]">
              Создание опроса
            </h1>
            <p className="text-gray-500 mt-2">Все опросы поддерживают множественный выбор</p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-6">
              <div>
                <label htmlFor="poll-title" className="block text-lg font-bold text-gray-700 mb-2">
                  Ваш вопрос
                </label>
                <input
                  type="text"
                  id="poll-title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Например: Какие технологии вы используете?"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00B39F] transition-shadow"
                  required
                />
              </div>
              <div>
                <label className="block text-lg font-bold text-gray-700 mb-2">
                  Варианты ответа
                </label>
                <div className="space-y-3">
                  {options.map((option, index) => (
                    <div key={index} className="flex items-center gap-2">
                      <input
                        type="text"
                        value={option}
                        onChange={(e) => handleOptionChange(index, e.target.value)}
                        placeholder={`Вариант ${index + 1}`}
                        className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#00B39F] transition-shadow"
                        required
                      />
                      {options.length > 2 && (
                        <button
                          type="button"
                          onClick={() => removeOption(index)}
                          className="text-gray-400 hover:text-red-500 transition-colors p-2"
                          aria-label="Удалить вариант"
                        >
                          <XCircle size={24} />
                        </button>
                      )}
                    </div>
                  ))}
                </div>
                <button
                  type="button"
                  onClick={addOption}
                  className="mt-4 flex items-center gap-2 text-[#00B39F] font-semibold hover:text-[#008A7B] transition-colors"
                >
                  <PlusCircle size={20} />
                  Добавить вариант
                </button>
              </div>
            </div>
            {error && <p className="mt-6 text-center text-red-500">{error}</p>}
            <div className="mt-10 flex flex-col sm:flex-row gap-4">
              <button
                type="submit"
                disabled={isLoading}
                className="w-full px-6 py-3 bg-[#00B39F] text-white font-bold rounded-lg shadow-lg shadow-[#00B39F]/30 hover:shadow-xl hover:shadow-[#00B39F]/40 hover:bg-[#00B39F]/90 transition-all tracking-wide disabled:bg-gray-400 disabled:shadow-none"
              >
                {isLoading ? 'Создание...' : 'Создать опрос'}
              </button>
              <Link href="/" className="w-full">
                <button
                  type="button"
                  className="w-full px-6 py-3 bg-gray-100 text-[#0B2B4A] font-bold rounded-lg hover:bg-gray-200 transition-all tracking-wide text-center"
                >
                  Отмена
                </button>
              </Link>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
}