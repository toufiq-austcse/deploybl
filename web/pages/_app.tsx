import "@/styles/globals.css";
import type {AppProps} from "next/app";
import Navbar from "@/components/ui/navbar";
import * as React from "react";
import {useEffect, useState} from "react";

function MyApp({Component, pageProps}: AppProps) {
    const [client, setClient] = useState<any>(null);

    useEffect(() => {
        // const apolloClient = createApolloClient();
        // setClient(apolloClient);
    }, []);

    return (
        <div className={"min-h-screen flex flex-col"}>
            <Navbar/>
            <div className="m-4 max-w-full px-5 sm:px-6 lg:px-20">
                <Component {...pageProps} />
            </div>
        </div>
    );
}

export default MyApp;
