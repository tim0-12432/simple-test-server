import {Progress as ShadProgress} from "@/components/ui/progress";
import { useEffect, useState } from "react";

type ProgressProps = {
    active: boolean;
    className?: string;
    value?: number;
}

const Progress = (props: ProgressProps) => {
    const [progress, setProgress] = useState(0);
    const [timer, setTimer] = useState<NodeJS.Timeout | null>(null);
    
    useEffect(() => {
        if (props.active && !props.value) {
            setProgress(13);
            const timer = setInterval(() => setProgress(prev => prev > 80 ? prev : prev + 18), 100);
            setTimer(timer);
        } else {
            setProgress(props.value || 0);
            if (timer) {
                clearTimeout(timer);
                setTimer(null);
            }
        }
        return () => {
            if (timer) {
                clearTimeout(timer);
            }
            setTimer(null);
        };
    }, [props.active, props.value]);

    if (!props.active) {
        return <div className={props.className}></div>;
    }
    return (
        <ShadProgress value={props.value || progress} className={props.className + " [&>div]:bg-gradient-to-r [&>div]:from-cyan-400 [&>div]:via-sky-500 [&>div]:to-indigo-500 [&>div]:rounded-l-full"}></ShadProgress>
    );
}

export default Progress;
