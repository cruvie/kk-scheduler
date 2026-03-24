import {createGrpcWebTransport} from "@connectrpc/connect-web";
import {authInterceptor} from "./interceptor";

export const transport = createGrpcWebTransport({
    baseUrl: "http://localhost:8667",
    interceptors: [authInterceptor],
});