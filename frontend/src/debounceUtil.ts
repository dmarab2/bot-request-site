export type DebouncedFunction<T> = {
    (): Promise<T>;
    cancel: () => void;
}

// general debounce function to use for various parts of the application
export function returnDebounceTest<T>(callback: (controller: AbortController) => T | Promise<T>, delay: number): DebouncedFunction<T> {
    let timeoutID: ReturnType<typeof setTimeout> | undefined;
    let currentController: AbortController | undefined;

    const debouncedFn = function () {
        return new Promise<T>((resolve, reject) => {
            if (timeoutID) clearTimeout(timeoutID);
            if (currentController) currentController.abort();

            currentController = new AbortController();

            timeoutID = setTimeout(async () => {
                try{
                    const result = await callback(currentController!);
                    resolve(result)
                } catch(error) {
                    reject(error);
                }
            }, delay);

        });
    };
    debouncedFn.cancel = () => {
        if (timeoutID) clearTimeout(timeoutID);
        if (currentController) currentController.abort();
    };
    return debouncedFn;
}