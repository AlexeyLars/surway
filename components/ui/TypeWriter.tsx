'use client';

import Typewriter from 'typewriter-effect';


interface TypewriterTextProps {
    strings: string[]
}

export default function TypewriterText({strings} : TypewriterTextProps) {
  return (
    <h1 className="text-xl md:text-2xl text-[#0B2B4A]/80">
      <Typewriter
        onInit={(typewriter) => {
          strings.forEach((text) => {
            typewriter
              .typeString(text)
              .pauseFor(2000)     
              .deleteAll(20)         
          });
          typewriter.start();
        }}
        options={{
          loop: true,
          delay: 55,
        }}
      />
    </h1>
  );
}