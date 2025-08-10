import React from "react";
import Image from "next/image";
import Link from "next/link";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <div className="min-h-screen bg-white text-gray-900">
      {/* Static sidebar for desktop */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:flex lg:w-20 lg:flex-col lg:overflow-y-auto lg:bg-gray-900 lg:pb-4">
        <div className="flex h-16 shrink-0 items-center justify-center">
          <Image
            src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500"
            alt="Your Company"
            width={32}
            height={32}
            className="h-8 w-auto"
          />
        </div>
        <nav className="mt-8 flex flex-col items-center space-y-1">
          <Link
            href="/books"
            className="group rounded-md p-3 bg-gray-800 text-white"
          >
            <span className="sr-only">Dashboard</span>
            📚
          </Link>
        </nav>
      </div>

      {/* Main wrapper */}
      <div className="lg:pl-20">
        {/* Topbar */}
        <header className="sticky top-0 z-40 flex h-16 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow sm:px-6 lg:px-8">
          <div className="text-lg font-semibold">Library</div>
        </header>

        {/* Main content */}
        <main className="p-4 sm:p-6 lg:p-8">{children}</main>
      </div>
    </div>
  );
}
