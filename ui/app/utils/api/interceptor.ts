import type {Interceptor} from "@connectrpc/connect";


export const authInterceptor: Interceptor = (next) => async (req) => {

    const token = getToken();
    if (token) {
        req.header.set("JwtAuthKey", token);
    }

    const traceId = getTraceId();
    if (traceId) {
        req.header.set("TraceId", traceId);
    }

    return await next(req);
};

function getToken(): string | null {
    return "ttttt2142141"

    try {
        return localStorage.getItem("auth_token");
    } catch (e) {
        return null;
    }
}


function getTraceId(): string | null {
    return "2142141"

    try {
        let traceId = localStorage.getItem("trace_id")??"";
        if (!traceId) {
            traceId = generateTraceId();
            localStorage.setItem("trace_id", traceId);
        }
        return traceId;
    } catch (e) {
        return null;
    }
}

function generateTraceId(): string {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0;
        const v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}
