import { useEffect, useState } from 'react';
import './App.css'

type requestStatus = "open" | "in_progress" | "fulfilled" | "cancelled"

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

export default function App() {
    const [requestList, setRequestList] = useState<RequestJson>({data: [], pageNumber: 0, nextLimit: false, prevLimit: false})
    const [selectedRequest, setSelectedRequest] = useState<Request>({id: 0, createdAt: "", updatedAt: "", requestText: "", requestStatus: "open"})
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
            <ul><RequestLister requestList={requestList} onRequestClick={handleSetRequest} /></ul>
            <ViewBox selectedRequest={selectedRequest}  />
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
        </div>
    )
}

async function fetchRequestList(): Promise<RequestJson>{
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