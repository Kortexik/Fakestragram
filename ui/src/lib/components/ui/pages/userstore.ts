// userStore.ts
import { writable } from "svelte/store";
import { API_URL } from "../../../../main";

export const currentUsername = writable("");
export const currentUserId = writable("");


export const fetchCurrentUser = async () => {
    try {
        const response = await fetch(`${API_URL}/protected/me`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
        });
        if (response.ok) {
            const data = await response.json();
            currentUsername.set(data.username);
            currentUserId.set(data.userID)
        } else {
            currentUsername.set("");
        }
    } catch {
        currentUsername.set("");
    }
};
