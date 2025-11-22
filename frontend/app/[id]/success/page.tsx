import Link from 'next/link';
import { ArrowRight, BarChart2, MessageSquare, Check } from 'lucide-react';
import { getPoll } from '@/app/services/api';
import CopyLinkButton from '@/components/ui/CopyLinkButton';


function getBaseUrl() {
  if (process.env.VERCEL_URL) return `https://${process.env.VERCEL_URL}`;
  return `http://localhost:3000`; 
}

export default async function PollSuccessPage({ params }: { params: { id: string } }) {
  const pollId = params.id;

  let pollData;
  try {
    pollData = await getPoll(pollId);
  } catch (error) {
    return (
      <main className="min-h-screen bg-[#F4F7FB] flex items-center justify-center p-4">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-500">Опрос не найден</h1>
          <p className="text-gray-600">Не удалось загрузить данные для этого опроса.</p>
          <Link href="/create" className="mt-4 inline-block text-[#00B39F] hover:underline">
            Создать новый опрос
          </Link>
        </div>
      </main>
    );
  }

  const baseUrl = getBaseUrl();
  const voteUrl = `${baseUrl}/${pollId}`;
  const resultsUrl = `${baseUrl}/${pollId}/results`;

  return (
    <main className="min-h-screen bg-[#F4F7FB] py-12 md:py-20">
      <div className="px-4">
        <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-lg p-6 sm:p-8 md:p-10">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check size={32} className="text-green-600" />
            </div>
            <h1 className="text-3xl md:text-4xl font-bold text-[#0B2B4A]">
              Опрос успешно создан!
            </h1>
            <p className="text-gray-500 mt-2 text-lg">
                «{pollData.poll.title}»
            </p>
          </div>

          <div className="space-y-6">
            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <label className="block text-lg font-bold text-gray-700 mb-2">
                Ссылка для голосования
              </label>
              <p className="text-sm text-gray-500 mb-3">
                Отправьте эту ссылку участникам, чтобы они могли отдать свой голос.
              </p>
              <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3 bg-white p-3 rounded-md border">
                <a 
                  href={voteUrl} 
                  target="_blank" 
                  rel="noopener noreferrer" 
                  className="font-mono text-sm text-[#00B39F] break-all hover:underline flex-grow"
                >
                  {voteUrl}
                </a>
                <CopyLinkButton textToCopy={voteUrl} />
              </div>
            </div>


            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <label className="block text-lg font-bold text-gray-700 mb-2">
                Ссылка на результаты
              </label>
              <p className="text-sm text-gray-500 mb-3">
                Используйте эту ссылку, чтобы отслеживать статистику в реальном времени.
              </p>
              <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3 bg-white p-3 rounded-md border">
                <a 
                  href={resultsUrl} 
                  target="_blank" 
                  rel="noopener noreferrer" 
                  className="font-mono text-sm text-[#00B39F] break-all hover:underline flex-grow"
                >
                  {resultsUrl}
                </a>
                <CopyLinkButton textToCopy={resultsUrl} />
              </div>
            </div>
          </div>
          
          <div className="mt-10 border-t pt-6 flex flex-col sm:flex-row gap-4 justify-between">
            <Link href="/create" className="w-full sm:w-auto">
              <button
                  type="button"
                  className="w-full px-6 py-3 bg-gray-100 text-[#0B2B4A] font-bold rounded-lg hover:bg-gray-200 transition-all tracking-wide text-center"
                >
                  Создать ещё один опрос
              </button>
            </Link>
            <Link href={resultsUrl} className="w-full sm:w-auto">
              <button
                  type="button"
                  className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-[#00B39F] text-white font-bold rounded-lg shadow-lg shadow-[#00B39F]/30 hover:shadow-xl hover:bg-[#00B39F]/90 transition-all tracking-wide"
                >
                  Посмотреть результаты
                  <ArrowRight size={20} />
              </button>
            </Link>
          </div>
        </div>
      </div>
    </main>
  );
}