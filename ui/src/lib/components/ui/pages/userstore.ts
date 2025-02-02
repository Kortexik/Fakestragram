// userStore.ts
import { writable } from "svelte/store";

export const currentUsername = writable("");
export const currentUserId = writable("");


export const fetchCurrentUser = async () => {
    try {
        const response = await fetch("http://localhost:8080/protected/me", {
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
