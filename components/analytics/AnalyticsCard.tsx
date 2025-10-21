'use client';

import ProgressRing from "./ProgressRing";
import Card from "./Card";

interface AnalyticsCardProps {
    responses: number,
    time: number,
    percentage: number
}

export default function AnalyticsCard({ time, responses, percentage }: AnalyticsCardProps) {
    return (
        <div className="bg-white rounded-2xl shadow-2xl p-6 md:p-8 w-full max-w-md mx-auto">

            <div className="flex justify-center mb-6 md:mb-8">
                <ProgressRing percentage={percentage} />
            </div>
            

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <Card title="Average completion time" value={time.toString() + "s"}/>
                <Card title="Total responses" value={responses.toString()} />
            </div>
        </div>
    )
}