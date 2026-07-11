export type DebouncedFunction<T> = {
    (): Promise<T>;
    cancel: () => void;
}

export function returnDebounceTest<T>(callback: () => T | Promise<T>, delay: number): DebouncedFunction<T> {
    let timeoutID: ReturnType<typeof setTimeout> | undefined;

    const debouncedFn = function () {
        return new Promise<T>((resolve, reject) => {
            if (timeoutID) clearTimeout(timeoutID);
            timeoutID = setTimeout(async () => {
                try{
                    const result = await callback();
                    resolve(result)
                } catch(error) {
                    reject(error);
                }
            }, delay);

        });
    };
    debouncedFn.cancel = () => {
        if (timeoutID) clearTimeout(timeoutID);
    };
    return debouncedFn;
}