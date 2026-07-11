import { useState, useEffect, useRef } from 'react';
import { returnDebounceTest, type DebouncedFunction } from './debounceUtil';

const MOCK_TAGS: string[] = [
    "1girl", "1boy", "solo", "long_hair", "short_hair", "blonde_hair",
    "blue_eyes", "brown_eyes", "holding_hands", "smile", "blush",
    "background", "scenery", "highres", "masterpiece", "absurdres"
];


export function useDebouncedSuggestions(currentWord: string, delay: number = 400) {
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [loading, setLoading] = useState<boolean>(false);

    function returnMockTags(): string[] {
        const lower = currentWord.toLowerCase();
        return MOCK_TAGS
            .filter((t) => t.toLowerCase().startsWith(lower))
            .slice(0, 10); // using 10 here since 10 suggestions is a standard on sites with tags
    }

    // a mutable reference to track the active debounce controller instance
    const debounceRef = useRef<DebouncedFunction<string[]> | null>(null);

    useEffect(() => {
        // if there is no word being typed, empty out suggestions and cancel any active timers
        if (!currentWord.trim()) {
            setSuggestions([]);
            setLoading(false);
            debounceRef.current?.cancel();
            return;
        }

        setLoading(true);

        const debouncedFetch = returnDebounceTest(returnMockTags, 400)

        debounceRef.current = debouncedFetch;

        debouncedFetch()
            .then((remoteData) => {
                setSuggestions(remoteData);
            })
            .catch((err) => {
                console.error("Error on the backend:", err);
            })
            .finally(() => {
                setLoading(false);
            });

        // cancels the timer automatically on the next keystroke or component unmount
        return () => {
            debouncedFetch.cancel();
        };
    }, [currentWord, delay]);

    // expose an explicit cancel command to the component body
    const forceCancel = () => {
        debounceRef.current?.cancel();
        setLoading(false);
    };

    return { suggestions, loading, forceCancel };
}