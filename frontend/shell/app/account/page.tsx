'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function AccountPage() {
    const router = useRouter();

    useEffect(() => {
        router.replace('/account/profile');
    }, [router]);

    return (
        <div className="flex items-center justify-center h-40">
            <div className="loading-spinner" />
        </div>
    );
}
