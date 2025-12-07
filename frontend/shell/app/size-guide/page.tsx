'use client';

import React, { useState } from 'react';

interface SizeChart {
    category: string;
    sizes: string[];
    measurements: { label: string; values: Record<string, string> }[];
}

const SIZE_CHARTS: SizeChart[] = [
    {
        category: 'Áo',
        sizes: ['S', 'M', 'L', 'XL', 'XXL'],
        measurements: [
            { label: 'Chiều dài (cm)', values: { S: '65', M: '68', L: '71', XL: '74', XXL: '77' } },
            { label: 'Vai (cm)', values: { S: '42', M: '44', L: '46', XL: '48', XXL: '50' } },
            { label: 'Ngực (cm)', values: { S: '96', M: '100', L: '104', XL: '108', XXL: '112' } },
            { label: 'Cân nặng (kg)', values: { S: '45-55', M: '55-65', L: '65-75', XL: '75-85', XXL: '85-95' } },
        ],
    },
    {
        category: 'Quần',
        sizes: ['S', 'M', 'L', 'XL', 'XXL'],
        measurements: [
            { label: 'Vòng eo (cm)', values: { S: '68-72', M: '72-76', L: '76-80', XL: '80-84', XXL: '84-88' } },
            { label: 'Vòng mông (cm)', values: { S: '88-92', M: '92-96', L: '96-100', XL: '100-104', XXL: '104-108' } },
            { label: 'Chiều dài (cm)', values: { S: '98', M: '100', L: '102', XL: '104', XXL: '106' } },
        ],
    },
    {
        category: 'Giày',
        sizes: ['38', '39', '40', '41', '42', '43', '44'],
        measurements: [
            { label: 'Chiều dài chân (cm)', values: { '38': '24', '39': '24.5', '40': '25', '41': '25.5', '42': '26', '43': '26.5', '44': '27' } },
            { label: 'EU', values: { '38': '38', '39': '39', '40': '40', '41': '41', '42': '42', '43': '43', '44': '44' } },
            { label: 'US Nam', values: { '38': '6', '39': '6.5', '40': '7', '41': '8', '42': '8.5', '43': '9.5', '44': '10' } },
            { label: 'UK', values: { '38': '5', '39': '5.5', '40': '6', '41': '7', '42': '7.5', '43': '8.5', '44': '9' } },
        ],
    },
];

export default function SizeGuidePage() {
    const [activeCategory, setActiveCategory] = useState('Áo');
    const [footLength, setFootLength] = useState('');
    const [recommendedSize, setRecommendedSize] = useState<string | null>(null);

    const activeChart = SIZE_CHARTS.find(c => c.category === activeCategory);

    const calculateShoeSize = () => {
        const length = parseFloat(footLength);
        if (!length) return;

        if (length <= 24) setRecommendedSize('38');
        else if (length <= 24.5) setRecommendedSize('39');
        else if (length <= 25) setRecommendedSize('40');
        else if (length <= 25.5) setRecommendedSize('41');
        else if (length <= 26) setRecommendedSize('42');
        else if (length <= 26.5) setRecommendedSize('43');
        else setRecommendedSize('44');
    };

    return (
        <div className="min-h-screen bg-[#f5f5f5] dark:bg-gray-900">
            <div className="container mx-auto px-4 py-6">
                <h1 className="text-xl font-bold mb-6 dark:text-white">📏 Hướng Dẫn Chọn Size</h1>

                {/* Category Tabs */}
                <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm mb-6">
                    <div className="flex">
                        {SIZE_CHARTS.map(chart => (
                            <button
                                key={chart.category}
                                onClick={() => setActiveCategory(chart.category)}
                                className={`flex-1 py-3 text-sm font-medium border-b-2 transition-colors ${activeCategory === chart.category
                                        ? 'text-[#ee4d2d] border-[#ee4d2d]'
                                        : 'text-gray-500 dark:text-gray-400 border-transparent'
                                    }`}
                            >
                                {chart.category}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Size Chart */}
                {activeChart && (
                    <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm overflow-hidden mb-6 animate-fade-in">
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead>
                                    <tr className="bg-[#fef6f5] dark:bg-gray-700">
                                        <th className="px-4 py-3 text-left text-sm font-medium dark:text-white">Thông số</th>
                                        {activeChart.sizes.map(size => (
                                            <th key={size} className="px-4 py-3 text-center text-sm font-medium dark:text-white">
                                                {size}
                                            </th>
                                        ))}
                                    </tr>
                                </thead>
                                <tbody className="divide-y dark:divide-gray-700">
                                    {activeChart.measurements.map((m, i) => (
                                        <tr key={m.label} className={i % 2 === 0 ? 'bg-gray-50 dark:bg-gray-750' : ''}>
                                            <td className="px-4 py-3 text-sm font-medium dark:text-white">{m.label}</td>
                                            {activeChart.sizes.map(size => (
                                                <td key={size} className="px-4 py-3 text-center text-sm dark:text-gray-300">
                                                    {m.values[size]}
                                                </td>
                                            ))}
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                )}

                {/* Size Calculator (for shoes) */}
                {activeCategory === 'Giày' && (
                    <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-6 mb-6 animate-fade-in-up">
                        <h3 className="font-medium mb-4 dark:text-white">🧮 Tính Size Giày</h3>
                        <div className="flex items-end gap-4">
                            <div className="flex-1">
                                <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                                    Chiều dài chân của bạn (cm)
                                </label>
                                <input
                                    type="number"
                                    value={footLength}
                                    onChange={(e) => setFootLength(e.target.value)}
                                    placeholder="VD: 25.5"
                                    className="w-full border px-4 py-3 rounded-sm outline-none focus:border-[#ee4d2d] dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    step="0.1"
                                />
                            </div>
                            <button
                                onClick={calculateShoeSize}
                                className="px-6 py-3 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                            >
                                Tính size
                            </button>
                        </div>

                        {recommendedSize && (
                            <div className="mt-4 p-4 bg-green-50 dark:bg-green-900 rounded-sm animate-fade-in">
                                <p className="text-green-700 dark:text-green-300">
                                    ✓ Size đề xuất cho bạn: <strong className="text-xl">{recommendedSize}</strong>
                                </p>
                            </div>
                        )}
                    </div>
                )}

                {/* How to Measure */}
                <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-6">
                    <h3 className="font-medium mb-4 dark:text-white">📐 Cách đo</h3>
                    <div className="grid md:grid-cols-3 gap-6">
                        <div className="text-center">
                            <div className="w-20 h-20 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-3xl mb-3">
                                👕
                            </div>
                            <h4 className="font-medium mb-2 dark:text-white">Áo</h4>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Đo từ đường may vai xuống gấu áo (chiều dài), và đo ngang ngực dưới nách (vòng ngực).
                            </p>
                        </div>
                        <div className="text-center">
                            <div className="w-20 h-20 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-3xl mb-3">
                                👖
                            </div>
                            <h4 className="font-medium mb-2 dark:text-white">Quần</h4>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Đo vòng eo ở vị trí nhỏ nhất, đo vòng mông ở vị trí lớn nhất.
                            </p>
                        </div>
                        <div className="text-center">
                            <div className="w-20 h-20 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-3xl mb-3">
                                👟
                            </div>
                            <h4 className="font-medium mb-2 dark:text-white">Giày</h4>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Đo từ gót chân đến đầu ngón chân dài nhất. Nên đo vào buổi chiều khi chân nở ra.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
