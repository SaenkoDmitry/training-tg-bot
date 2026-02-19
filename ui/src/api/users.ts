import {api} from "./client.ts";

export const getUserIcon = () =>
    api<IconDTO>(`/api/users/icon`, {method: "GET"});

export const changeUserIcon = (name: string) =>
    api(`/api/users/change-icon`, {
        method: "POST",
        body: JSON.stringify({
            name: name,
        }),
    });
