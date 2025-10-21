import ProgressRing from "../components/analytics/ProgressRing";
import Card from "../components/analytics/Card";
import AnalyticsCard from "../components/analytics/AnalyticsCard";


export default function Home() {
  return (
    <main className="py-12">
      <AnalyticsCard responses={1244} time={58} percentage={84} />
    </main>
    
  );
}