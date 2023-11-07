export function GetDataFromEvent<T>(message: Buffer): T {
    const event = JSON.parse(message?.toString()) as Event;
    const header = event.header;
    console.log("header", JSON.stringify(header));
    const data = event?.body;
    if (typeof data == "string") {
        return JSON.parse(data) as T;
    } else {
        return data as T;
    }
}

interface Event {
    header: Header;
    body: any;
}

interface Header {
    version?: string;
    timestamp?: string;
    orgService?: string;
    from?: string;
    channel?: string;
    broker?: string;
    session?: string;
    transaction?: string;
    communication?: string;
    groupTags?: any[];
    identity?: Identity;
    baseApiVersion?: string;
    schemaVersion?: string;
    instanceData?: string;
}

interface Identity {
    device: number;
}
