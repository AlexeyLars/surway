'use client';

import { useState } from 'react';
import { Clipboard, Check } from 'lucide-react';

export default function CopyLinkButton({ textToCopy }: { textToCopy: string }) {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(textToCopy).then(() => {
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    });
  };

  return (
    <button
      onClick={handleCopy}
      className="flex items-center gap-2 px-4 py-2 text-sm font-semibold bg-gray-100 text-[#0B2B4A] rounded-lg hover:bg-gray-200 transition-all"
    >
      {isCopied ? (
        <>
          <Check size={16} className="text-green-500" />
          Скопировано
        </>
      ) : (
        <>
          <Clipboard size={16} />
          Копировать
        </>
      )}
    </button>
  );
}