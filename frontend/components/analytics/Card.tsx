'use client';

interface CardProps {
    title: string,
    value: string
}

export default function StatCard({ title, value }: CardProps) {
    return (
        <div className="flex flex-col gap-2 rounded-lg p-4 bg-slate-50 bg-opacity-60">
            <p className="text-slate-800 text-sm md:text-base font-medium leading-normal">
                {title}
            </p>
            <p className="text-slate-800 tracking-light text-xl md:text-2xl font-bold leading-tight">
                {value}
            </p>
        </div>
    )
}