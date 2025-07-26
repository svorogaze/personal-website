import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import React from "react";
import { Suspense } from "react";
import Link from 'next/link';
import ClientSideScrollRestorer from '../src/ClientSideScrollRestorer';
import {HeroUIProvider} from "@heroui/system";
import "./globals.css";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "svorogaze",
    description: "Personal portfolio of svorogaze",
};


export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" className="dark hydrated" data-arp="">
            <body
                className={`${geistSans.variable} ${geistMono.variable} antialiased`}
                cz-shortcut-listen="true"
            >
                <HeroUIProvider>
                <div className="navigation-bar">
                    <div className='flex items-center justify-center px-20 font-semibold'>
                        <div className='flex gap-2 xl:gap-10 md:gap-8'>
                            <Link href="/" className="navigation-bar-element">About me</Link>
                            <Link href="/blog" className="navigation-bar-element">Blog</Link>
                        </div>
                    </div>
                </div>
                {children}
                </HeroUIProvider>
            </body>
            <Suspense>
                <ClientSideScrollRestorer />
            </Suspense>
        </html>
    );
}
