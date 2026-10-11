import React, { useEffect, useState, useRef, useMemo, useCallback} from 'react';
import './App.css'
import { useDebouncedSuggestions } from './useDebouncedSuggestions';

const mockRequestList: Request[] = [
    {
        id: 1,
        createdAt: "2026-08-22",
        updatedAt: "2026-08-22",
        requestText: "mock text",
        requestStatus: "open"
    },
    {
        id: 2,
        createdAt: "2026-08-22",
        updatedAt: "2026-08-22",
        requestText: "mock text",
        requestStatus: "open"
    },
    {
        id: 3,
        createdAt: "2026-08-22",
        updatedAt: "2026-08-22",
        requestText: "mock text",
        requestStatus: "open"
    },
    {
        id: 4,
        createdAt: "2026-08-22",
        updatedAt: "2026-08-22",
        requestText: "mock text",
        requestStatus: "open"
    },
    {
        id: 5,
        createdAt: "2026-08-22",
        updatedAt: "2026-08-22",
        requestText: "mock text",
        requestStatus: "open"
    },
]

type requestStatus = "open" | "in_progress" | "fulfilled" | "cancelled" | ""

interface Request {
    id: number;
    createdAt: string;
    updatedAt: string;
    requestText: string;
    requestStatus: requestStatus
}

interface RequestPayload {
    id: number;
    created_at: string;
    updated_at: string;
    request_text: string;
    status_: requestStatus
}

interface NewRequest {
    body: string
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
    selectedRequestID: number
}

interface tagSearchBoxProps {
    onParentChange?: (value: string) => void;
    name: string
}

interface formProps {
    requestList: RequestJson
    setRequestList: React.Dispatch<React.SetStateAction<RequestJson>>
}

interface searchButtonProps {
    setRequestList: React.Dispatch<React.SetStateAction<RequestJson>>
}

export default function App() {
    const [requestList, setRequestList] = useState<RequestJson>({data: [], pageNumber: 0, nextLimit: false, prevLimit: false})
    const [selectedRequest, setSelectedRequest] = useState<Request>({id: 0, createdAt: "", updatedAt: "", requestText: "", requestStatus: ""})
    const [initialRun, setInitialRun] = useState<Boolean>(false)
    useEffect(() => {
        if (initialRun == false){
            setInitialRun(true)
            fetchRequestList()
            .then((data) => { setRequestList(data);}) 
            .catch((err) => {console.error(err);});
        }
    }, [])
    function handleSetRequest(request: Request){
        setSelectedRequest(request);
    }

    return (
        <div className='min-h-screen bg-zinc-950/90 bg-[url(./assets/grit.png)] bg-blend-multiply bg-fixed text-zinc-200'>
            <div className="mx-auto grid max-w-6xl gap-6 p-6 lg:grid-cols-[20rem_1fr]">
                <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
                <Panel className='flex flex-col m-5 items-center'><RequestSearchForm requestList={requestList} setRequestList={setRequestList}/></Panel>
                <NewRequestForm />
                </aside>
                <main className='space-y-6'>
                    <Panel className='flex-col'><RequestLister requestList={requestList} onRequestClick={handleSetRequest} selectedRequestID={selectedRequest.id} /></Panel>
                    <Panel className='flex-col'><ViewBox selectedRequest={selectedRequest}  /></Panel>
                </main>
            </div>
        </div>
    );
}

function Panel({ className = "", children }: { className?: string; children: React.ReactNode }) {
  return (
    <div className={`flex rounded-xl bg-zinc-850/10 backdrop-blur-md ring-1 ring-white/10 shadow-lg shadow-black/50 inset-shadow-xs inset-shadow-white/10 ${className}`}>
      {children}
    </div>
  );
}

// should add a way to change the color, but I don't know if I should add another color for buttons?
function Button({ className = "", onClickFunc, buttonText }: { className?: string; onClickFunc: () => void; buttonText: string }) {
  return (
    <button className={`${className} px-3 m-5 bg-indigo-900/70 rounded-xs border-slate-600 shadow-2xl shadow-[inset_0_1px_0_0_rgba(255,255,255,0.15),_0_4px_12px_rgba(0,0,0,0.6)]`} onClick={() => onClickFunc()}>
        {buttonText}
    </button>
  );
}

function RequestLister( { requestList, onRequestClick, selectedRequestID }: RequestSetterProp){
    
    const badge: Record<requestStatus, string> = {
        open: "bg-emerald-500/15 text-emerald-300",
        in_progress: "bg-amber-500/15 text-amber-300",
        fulfilled: "bg-indigo-500/15 text-indigo-300",
        cancelled: "bg-zinc-500/15 text-zinc-400",
        "": "",
    };

    const listItems = requestList.data.map(request => 
        <li key={request.id}>
            <button
                onClick={() => onRequestClick(request)}
                className={`flex w-full items-center justify-between gap-4 px-4 py-3 text-left transition-colors hover:bg-white/5 ${
                request.id === selectedRequestID ? "bg-indigo-500/10 border-l-2 border-indigo-400" : "border-l-2 border-transparent"
                }`}
            >
                <span className='truncate'>{request.requestText}</span>
                <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs ${badge[request.requestStatus]}`}>{request.requestStatus}</span>
            </button>
        </li>
    );

    return (
        <ul className='flex justify-center items-center flex-col'>
            {listItems}
        </ul>
    )

}

function ViewBox( { selectedRequest }: {selectedRequest: Request} ) {
    const listParas = [
        <p>Request: {selectedRequest.requestText}</p>,
        <p>Status: {selectedRequest.requestStatus}</p>,
        <p>Created on: {selectedRequest.createdAt}</p>
    ]
    return (
        <div>
            {listParas}
        </div>
    )
}

function RequestSearchForm({ requestList, setRequestList }: formProps){
    return (
        <div>
            <RequestTextSearch />
            <RequestTagSearch name='tagSearch' />
            <RequestSearchButton setRequestList={setRequestList} />
        </div>
    )
}

function RequestSearchButton({ setRequestList }: searchButtonProps) {
    return (
        <>
            <button className="m-5 bg-indigo-900 rounded-xs border-slate-600 shadow-2xl shadow-[inset_0_1px_0_0_rgba(255,255,255,0.15),_0_4px_12px_rgba(0,0,0,0.6)]" onClick={() => setRequestList({data: mockRequestList, pageNumber: 0, nextLimit: false, prevLimit: false})}>Search</button>
        </>
    )
}

function RequestTextSearch() {
    const [value, setValue] = useState<string>("");
    const inputReference = useRef<HTMLInputElement>(null);

    return (
        <div className="relative">
            <input
            ref={inputReference}
            value={value}
            placeholder='Enter request text here.'
            className="w-full border rounded px-3 py-2 bg-black my-2"
            onChange={(e) => setValue(e.target.value)}
            ></input>
        </div>
    )
}


function RequestTagSearch({ onParentChange, name }: tagSearchBoxProps) {
    const [value, setValue] = useState<string>("");
    const [activeIndex, setActiveIndex] = useState<number>(0);
    const [showDropdown, setShowDropdown] = useState<boolean>(false);
    const inputReference = useRef<HTMLInputElement>(null);
    const listReference = useRef<HTMLLIElement[]>([]);
    const targetReference = useRef<HTMLLIElement>(null);

    const currentWord = useMemo(() => {
        const cursor = inputReference.current?.selectionStart ?? value.length;
        const fromCursor = value.slice(0, cursor);
        const lastWord = fromCursor.match(/\S+$/);
        return lastWord ? lastWord[0] : "";
    },[value])

    // debouncing is used to ensure that multiple requests aren't sent while typing is done
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

    // having a way to navigate suggestions with the keys is useful for typing many tags
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

    useEffect(() => {
        if (!showDropdown) return;
        const el = listReference.current[activeIndex]
        el?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        })
        }, [activeIndex, listReference])

    return (
        <div className='relative'>
            <input
            ref={inputReference}
            value={value}
            onChange={(e) => {
                setValue(e.target.value);
                onParentChange?.(e.target.value);
                setShowDropdown(true)
                setActiveIndex(0);
            }}
            name={name}
            onBlur={() => {setTimeout(() => setShowDropdown(false), 100); }}
            onFocus={() => setShowDropdown(true)}
            onKeyDown={handleKeyDown}
            className="w-full border rounded px-3 py-2 bg-black my-2"
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
                        ref={(el) => {listReference.current[index] = el!}}
                        >
                            {tag}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}

// I will leave these values here for now, in case I use them for something else later
function NewRequestForm() {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const newRequestRef = useRef<HTMLDialogElement>(null);
    const [textValue, setTextValue] = useState<string>("");
    const [tagValue, setTagValue] = useState<string>("");
    const [status, setStatus] = useState('idle');
    const [errorMessage, setErrorMessage] = useState<string>('');
    const [responseData, setResponseData] = useState<Request>();
    
    const openWindow = () => {
        newRequestRef.current?.showModal();
    }

    const closeWindow = () => {
        newRequestRef.current?.close();
    }
    
    useEffect(() => {
        const modal = newRequestRef.current
        if (!modal) return;
        
        if (isOpen) {
            openWindow();
            return;
        } else {
            closeWindow();
            return;
        }

    })

    async function formAction(event: React.SubmitEvent<HTMLFormElement>){
        event.preventDefault();
        setStatus("submitting")
        setErrorMessage("")

        const formData = new FormData(event.currentTarget)
        const bodyText = formData.get("requestText")?.toString()
        const tagText = formData.get("requestTags")?.toString().trim().split(" ")
        console.log(tagText)
        if (bodyText === undefined) throw new Error("Request text cannot be empty!")
        const textRequest: NewRequest = {
            body: bodyText
        }
        try{
            console.log("Now sending request...")
            const responseRequest = submitNewRequest(textRequest);
            responseRequest
            .then((data) => {
                setResponseData(data);
                console.log("Successfully submitted new request!");
                setStatus("Idle");
                event.target.reset();
            }) 
            .catch((err) => {throw new Error(err);});
        } catch(error){
            if (error instanceof Error){
                console.error("There was an error: ", error.message);
                setErrorMessage(error.message)
            }
            setStatus("error")
        }

    }
    
    return (
        <>
            <Button className='' onClickFunc={() => setIsOpen(true)} buttonText='Make a New Request'></Button>

            <dialog ref={newRequestRef} className={`p-5 rounded-xl bg-zinc-850/10 backdrop-blur-md ring-1 ring-white/10 shadow-lg shadow-black/50 inset-shadow-xs inset-shadow-white/10 m-auto`}>
                <form onSubmit={formAction}>
                    <div className='flex flex-col min-w-100 min-h-50'>
                        <h2>Enter your request</h2>
                        <textarea name='requestText' placeholder='Enter your request here.' className='min-h-1/2 px-3 py-2 my-2 bg-black' onChange={(e) => setTextValue(e.target.value)}></textarea>
                        <h2>(Optional): Add tags to your request.</h2>
                        <RequestTagSearch name='requestTags' onParentChange={setTagValue} />
                        <button type="submit" className="m-5 bg-indigo-900 rounded-xs border-slate-600 shadow-2xl shadow-[inset_0_1px_0_0_rgba(255,255,255,0.15),_0_4px_12px_rgba(0,0,0,0.6)]">{status === "submitting" ? "Submitting..." : "Submit Request"}</button>
                        <button type="button" onClick={() => setIsOpen(false)} className="m-5 bg-indigo-900 rounded-xs border-slate-600 shadow-2xl shadow-[inset_0_1px_0_0_rgba(255,255,255,0.15),_0_4px_12px_rgba(0,0,0,0.6)]">Close this Window</button>

                    </div>
                </form>
            </dialog>
        </>
    )
}


async function fetchRequestList(): Promise<RequestJson>{
    console.log(`${import.meta.env.VITE_BACKEND_ROOT}${import.meta.env.VITE_BACKEND_GET_OPEN_REQUESTS}`)
    try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_ROOT}${import.meta.env.VITE_BACKEND_GET_OPEN_REQUESTS}`);
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
        console.error("There was an error: ", error);
        throw error;
    }
}

async function submitNewRequest(newRequest: NewRequest): Promise<Request>{
    const postURL = import.meta.env.VITE_BACKEND_ROOT + import.meta.env.VITE_BACKEND_REQUESTS;
    try{
        const response = await fetch(postURL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newRequest)
        });
        if (!response.ok) {
            throw new Error(`Error! Status is ${response.status}`)
        }
        const data: RequestPayload = await response.json();
        const finishedRequest: Request = convertRequestPayloadToRequest(data);
        return finishedRequest;

    } catch(error) {
        console.error("There was an error: ", error)
        throw error;
    }
}

function convertRequestPayloadToRequest(payload: RequestPayload): Request{
    const newRequest: Request = {
        id: payload.id,
        createdAt: payload.created_at,
        updatedAt: payload.updated_at,
        requestText: payload.request_text,
        requestStatus: payload.status_
    };
    return newRequest;
}

function convertRequestToRequestPayload(request: Request): RequestPayload{
    const newPayload: RequestPayload = {
        id: request.id,
        created_at: request.createdAt,
        updated_at: request.updatedAt,
        request_text: request.requestText,
        status_: request.requestStatus
    };
    return newPayload;
}