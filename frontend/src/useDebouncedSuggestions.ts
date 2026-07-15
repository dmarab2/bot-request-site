import { useState, useEffect, useRef } from 'react';
import { returnDebounceTest, type DebouncedFunction } from './debounceUtil';


interface Tag {
    CreatedAt: string,
    ID: number,
    Name: string,
    PostCount: number,
    UpdatedAt: string,
}


export function useDebouncedSuggestions(currentWord: string, delay: number = 400) {
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [loading, setLoading] = useState<boolean>(false);

    async function returnTagSuggestions(controller: AbortController): Promise<string[]> {
        const response = await fetch(`http://localhost:8080/api/tags/${currentWord}`, {signal: controller.signal})
        if (!response.ok) {
            throw new Error(`There was an HTTP Error, Status: ${response.status}`);
        }
        const data: Tag[] = await response.json()
        console.log(data);
        let tagNames: string[] = [];
        data.forEach(element => {
           tagNames.push(element.Name) 
        });
        return tagNames;
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

        const debouncedFetch = returnDebounceTest(returnTagSuggestions, 250)

        //debounceRef.current = debouncedFetch;

        debouncedFetch()
            .then((tagList) => {
                setSuggestions(tagList);
                setLoading(false);
            })
            .catch((err) => {
                if(err instanceof DOMException && err.name=="AbortError"){
                    console.log(`Network request for ${currentWord} was successfully aborted.`);
                    return;
                }
                console.error("Error on the backend:", err);
            });
            /*
            .finally(() => {
                setLoading(false);
            });
            */

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