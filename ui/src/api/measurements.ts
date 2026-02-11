import {api} from "./client.ts";

export const deleteMeasurement = (id: number) =>
    api(`/api/measurements/${id}`, {method: "DELETE"});
