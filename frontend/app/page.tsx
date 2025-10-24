import AnalyticsCard from "../components/analytics/AnalyticsCard";
import TypewriterText from "../components/ui/TypeWriter";

export default function Home() {
  const strings = [
    "Опрос за 30 секунд.",
    "Мнения без регистрации.",
    "Надёжная аналитика в реальном времени.",
  ];

  return (
    <main className="min-h-screen bg-[#F4F7FB] py-12">
      <div className="px-4 md:px-10 lg:px-20 xl:px-40">
        <div className="max-w-[1280px] mx-auto">
          <div className="flex flex-col lg:flex-row gap-12 lg:gap-16 items-center">
            <div className="flex flex-col gap-6 lg:w-1/2">
              <div className="flex flex-col gap-4">
                <h1 className="text-5xl md:text-6xl font-bold text-[#0B2B4A] leading-tighter tracking-tighter">
                  Опросы, которые дают результат.
                </h1>
                <div className="h-14">
                  <TypewriterText strings={strings} />
                </div>
              </div>

              <div className="flex flex-col sm:flex-row gap-4">
                <button className="px-6 py-3 bg-[#00B39F] text-white font-bold rounded-lg shadow-lg shadow-[#00B39F]/30 hover:shadow-xl hover:shadow-[#00B39F]/40 hover:bg-[#00B39F]/90 transition-all tracking-wide">
                  Создать опрос
                </button>
                <button className="px-6 py-3 bg-white text-[#0B2B4A] font-bold rounded-lg hover:bg-gray-50 transition-all tracking-wide">
                  Посмотреть шаблоны
                </button>
              </div>
            </div>

            <div className="lg:w-1/2 flex justify-center">
              <AnalyticsCard responses={1244} time={58} percentage={84} />
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
