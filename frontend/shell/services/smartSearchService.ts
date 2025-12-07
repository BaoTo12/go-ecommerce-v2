// Smart Search Service - Fuzzy matching with typo correction

import { Product, productService } from './productService';

export interface SmartSearchResult {
    products: Product[];
    suggestions: string[];
    correctedQuery?: string;
    didYouMean?: string;
}

// Levenshtein distance for fuzzy matching
const levenshteinDistance = (str1: string, str2: string): number => {
    const m = str1.length;
    const n = str2.length;
    const dp: number[][] = Array(m + 1).fill(null).map(() => Array(n + 1).fill(0));

    for (let i = 0; i <= m; i++) dp[i][0] = i;
    for (let j = 0; j <= n; j++) dp[0][j] = j;

    for (let i = 1; i <= m; i++) {
        for (let j = 1; j <= n; j++) {
            if (str1[i - 1] === str2[j - 1]) {
                dp[i][j] = dp[i - 1][j - 1];
            } else {
                dp[i][j] = Math.min(
                    dp[i - 1][j - 1] + 1, // substitution
                    dp[i - 1][j] + 1,     // deletion
                    dp[i][j - 1] + 1      // insertion
                );
            }
        }
    }

    return dp[m][n];
};

// Vietnamese character normalization
const normalizeVietnamese = (str: string): string => {
    return str
        .toLowerCase()
        .normalize('NFD')
        .replace(/[\u0300-\u036f]/g, '')  // Remove diacritics
        .replace(/đ/g, 'd')
        .replace(/[^a-z0-9\s]/g, '')
        .trim();
};

// Common typo corrections
const TYPO_MAP: Record<string, string> = {
    'iphone': 'iphone',
    'ipone': 'iphone',
    'iphoen': 'iphone',
    'ifone': 'iphone',
    'samsugn': 'samsung',
    'samsnug': 'samsung',
    'samung': 'samsung',
    'nokie': 'nike',
    'nkie': 'nike',
    'addidas': 'adidas',
    'addias': 'adidas',
    'lapotp': 'laptop',
    'laptpo': 'laptop',
    'dior': 'dior',
    'doir': 'dior',
};

// Popular search terms for suggestions
const POPULAR_TERMS = [
    'iphone 15', 'samsung galaxy', 'macbook', 'laptop gaming',
    'áo hoodie', 'quần jean', 'giày nike', 'giày adidas',
    'son dior', 'serum vitamin c', 'kem chống nắng',
    'tai nghe bluetooth', 'airpods', 'đồng hồ casio',
    'túi xách', 'balo laptop', 'nồi chiên không dầu',
];

export const smartSearchService = {
    // Smart search with fuzzy matching
    search: async (query: string): Promise<SmartSearchResult> => {
        await new Promise(resolve => setTimeout(resolve, 200));

        const normalizedQuery = normalizeVietnamese(query);
        const words = normalizedQuery.split(/\s+/);

        // Check for typo corrections
        let correctedQuery: string | undefined;
        const correctedWords = words.map(word => {
            if (TYPO_MAP[word]) {
                correctedQuery = query; // Mark that correction happened
                return TYPO_MAP[word];
            }
            return word;
        });

        // Find similar terms using Levenshtein distance
        let didYouMean: string | undefined;
        if (!correctedQuery) {
            for (const term of POPULAR_TERMS) {
                const normalizedTerm = normalizeVietnamese(term);
                const distance = levenshteinDistance(normalizedQuery, normalizedTerm);
                if (distance <= 2 && distance > 0) {
                    didYouMean = term;
                    break;
                }
            }
        }

        // Search with original and corrected terms
        const searchTerms = [query, correctedWords.join(' ')].filter(Boolean);
        const allProducts = await productService.getProducts({ limit: 100 });
        const productList = Array.isArray(allProducts) ? allProducts : allProducts.products || [];

        // Score products based on relevance
        const scored = productList.map(product => {
            const productName = normalizeVietnamese(product.name);
            const productDesc = normalizeVietnamese(product.description);
            let score = 0;

            for (const term of searchTerms) {
                const normalizedTerm = normalizeVietnamese(term);

                // Exact match in name
                if (productName.includes(normalizedTerm)) {
                    score += 10;
                }

                // Exact match in description
                if (productDesc.includes(normalizedTerm)) {
                    score += 5;
                }

                // Word-by-word matching
                const termWords = normalizedTerm.split(/\s+/);
                for (const word of termWords) {
                    if (word.length >= 3) {
                        if (productName.includes(word)) score += 3;
                        if (productDesc.includes(word)) score += 1;
                    }
                }

                // Fuzzy matching for each word
                const productWords = productName.split(/\s+/);
                for (const pWord of productWords) {
                    for (const tWord of termWords) {
                        if (pWord.length >= 3 && tWord.length >= 3) {
                            const distance = levenshteinDistance(pWord, tWord);
                            if (distance <= 1) score += 2;
                            else if (distance <= 2) score += 1;
                        }
                    }
                }
            }

            return { product, score };
        })
            .filter(s => s.score > 0)
            .sort((a, b) => b.score - a.score);

        // Generate suggestions based on query
        const suggestions = POPULAR_TERMS
            .filter(term => {
                const normalizedTerm = normalizeVietnamese(term);
                return normalizedTerm.includes(normalizedQuery) ||
                    normalizedQuery.includes(normalizedTerm.split(' ')[0]);
            })
            .slice(0, 5);

        return {
            products: scored.slice(0, 20).map(s => s.product),
            suggestions,
            correctedQuery,
            didYouMean,
        };
    },

    // Get autocomplete suggestions
    getAutocomplete: async (query: string): Promise<string[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        const normalizedQuery = normalizeVietnamese(query);

        return POPULAR_TERMS
            .filter(term => normalizeVietnamese(term).includes(normalizedQuery))
            .slice(0, 8);
    },

    // Get trending searches
    getTrending: (): string[] => {
        return [
            '🔥 iPhone 15 Pro Max',
            '🔥 Áo hoodie form rộng',
            '🔥 Giày Nike Air Force 1',
            '🔥 Son Dior',
            '🔥 Laptop gaming',
        ];
    },
};

export default smartSearchService;
