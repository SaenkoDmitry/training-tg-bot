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

export const getProfile = () => api("/api/users/profile", {
    headers: {Authorization: `Bearer ${localStorage.getItem("token")}`},
})

export const updateProfile = (data: UpdateProfileRequest) => api("/api/users/profile", {
    method: "PATCH",
    headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
    body: JSON.stringify(data),
});
