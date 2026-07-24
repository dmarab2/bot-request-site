import { useEffect, useState, useRef, useMemo, useCallback} from 'react';
import './App.css'
import { useDebouncedSuggestions } from './useDebouncedSuggestions';

const MOCK_TAGS_HERE: string[] = [
  "1girl", "1boy", "solo", "long_hair", "short_hair", "blonde_hair", 
  "blue_eyes", "brown_eyes", "holding_hands", "smile", "blush", 
  "background", "scenery", "highres", "masterpiece", "absurdres"
];

type requestStatus = "open" | "in_progress" | "fulfilled" | "cancelled" | ""
type elementVisibility = "none" | "flex"

interface Request {
    id: number;
    createdAt: string;
    updatedAt: string;
    requestText: string;
    requestStatus: requestStatus

}

interface RequestJson {
    data: Request[];
    pageNumber: number;
    nextLimit: boolean;
    prevLimit: boolean
}

interface RequestSetterProp {
    requestList: RequestJson;
    onRequestClick: (request: Request) => void
}

interface searchBoxProps {
    suggestedTags: string[];
    onParentChange?: (value: string) => void;
}

export default function App() {
    const [requestList, setRequestList] = useState<RequestJson>({data: [], pageNumber: 0, nextLimit: false, prevLimit: false})
    const [selectedRequest, setSelectedRequest] = useState<Request>({id: 0, createdAt: "", updatedAt: "", requestText: "", requestStatus: ""})
    useEffect(() => {
        fetchRequestList()
        .then((data) => { setRequestList(data);}) 
        .catch((err) => {console.error(err);});
    }, [])
    function handleSetRequest(request: Request){
        setSelectedRequest(request);
    }

    return (
        <div>
            <div>
                <ul><RequestLister requestList={requestList} onRequestClick={handleSetRequest} /></ul>
                <ViewBox selectedRequest={selectedRequest}  />
                <RequestTextSearch />
                <RequestTagSearch suggestedTags={MOCK_TAGS_HERE} />
            </div>
            <aside>
            </aside>
        </div>
    );
}

function RequestLister( { requestList, onRequestClick }: RequestSetterProp){
    const listItems = requestList.data.map(request => 
        <li key={request.id} onClick={() => onRequestClick(request)}>
            {request.requestText}
        </li>
    );

    return listItems;

}

function ViewBox( { selectedRequest }: {selectedRequest: Request} ) {
    return (
        <div>
            <p>Request: {selectedRequest.requestText}</p>
            <p>Status: {selectedRequest.requestStatus}</p>
            <p>Created on: {selectedRequest.createdAt}</p>
        </div>
    )
}

function RequestTextSearch() {
    const [value, setValue] = useState<string>("");
    const inputReference = useRef<HTMLInputElement>(null);

    return (
        <div className="relative w-full">
            <input
            ref={inputReference}
            value={value}
            placeholder='Enter request text here.'
            ></input>
        </div>
    )
}


function RequestTagSearch({ suggestedTags, onParentChange }: searchBoxProps) {
    const [value, setValue] = useState<string>("");
    const [activeIndex, setActiveIndex] = useState<number>(0);
    const [showDropdown, setShowDropdown] = useState<boolean>(false);
    const inputReference = useRef<HTMLInputElement>(null);

    const currentWord = useMemo(() => {
        const cursor = inputReference.current?.selectionStart ?? value.length;
        const fromCursor = value.slice(0, cursor);
        const lastWord = fromCursor.match(/\S+$/);
        return lastWord ? lastWord[0] : "";
    },[value])

    const { suggestions, loading, forceCancel } = useDebouncedSuggestions(currentWord, 400)

    const applySuggestion = useCallback((tag: string) => { 
        if (!tag) return;
        forceCancel();
        const cursor = inputReference.current?.selectionStart ?? value.length;
        const before = value.slice(0, cursor).replace(/\S+$/, tag + " ");
        const after = value.slice(cursor);
        const newValue = before + after;
        setValue(newValue);
        onParentChange?.(newValue);
        setShowDropdown(false);

        requestAnimationFrame(() => {
            const pos = before.length;
            inputReference.current?.setSelectionRange(pos, pos);
            inputReference.current?.focus();
      });
    }, [value, onParentChange])

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (!showDropdown || suggestions.length === 0) return;
        if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveIndex((i) => (i + 1) % suggestions.length);
        } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveIndex((i) => (i - 1 + suggestions.length) % suggestions.length);
        } else if (e.key === "Enter" || e.key === "Tab") {
        e.preventDefault();
        applySuggestion(suggestions[activeIndex]);
        } else if (e.key === "Escape") {
        setShowDropdown(false);
        }
    };

    return (
        <div className="relative w-full">
            <input
            ref={inputReference}
            value={value}
            onChange={(e) => {
                setValue(e.target.value);
                onParentChange?.(e.target.value);
                setShowDropdown(true)
                setActiveIndex(0);
            }}
            onBlur={() => {setTimeout(() => setShowDropdown(false), 100); }}
            onFocus={() => setShowDropdown(true)}
            onKeyDown={handleKeyDown}
            className="w-full border rounded px-3 py-2"
            placeholder="Enter tags here."
            />
            {showDropdown && (suggestions.length > 0 || loading) && (
                <ul className="absolute z-10 mt-1 w-full bg-white border rounded shadow-md max-h-48 overflow-y-auto">
                    {suggestions.map((tag: string, index: number) => (
                        <li
                        key={tag}
                        onMouseDown={() => applySuggestion(tag)}
                        className={`px-3 py-1 cursor-pointer ${
                        index === activeIndex ? "bg-blue-100" : ""
                        }`}
                        >
                            {tag}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}

async function fetchRequestList(): Promise<RequestJson>{
    console.log(import.meta.env.VITE_BACKEND_URL)
    try {
        const response = await fetch(import.meta.env.VITE_BACKEND_URL);
        if (!response.ok) {
            throw new Error(`There was an HTTP Error, Status: ${response.status}`);
        }
        const data = await response.json();
        console.log(data)
        const transformedData: RequestJson = {
            data: data.data.map((request: any) => ({
                id: request.id,
                createdAt: request.created_at,
                updatedAt: request.updated_at,
                requestText: request.request_text,
                requestStatus: request.status_,
            })),
            pageNumber: data.page_number,
            nextLimit: data.next_limit,
            prevLimit: data.prev_limit,
        };
        return transformedData;
    } catch(error) {
        console.error("There was an error:", error);
        throw error;
    }
}