import { writable } from "svelte/store";

export const showPopup = writable(false);
export const isAuthenticated = writable(false);

export const usernameExists = writable(false);
export const emailExists = writable(false);

export const notifications = writable([]);

export const editing = writable(false);

export function setEditing(value) {
  editing.set(value);
}

export function setUsernameExists(value) {
  usernameExists.set(value);
}

export function setEmailExists(value) {
  emailExists.set(value);
}

export const togglePopup = () => {
  showPopup.update((visible) => !visible);
};

export function setAuthenticated(value) {
  isAuthenticated.set(value);
}

export async function checkAuthentication() {
  const token = localStorage.getItem("token");

  if (token) {
    try {
      const response = await fetch("http://localhost:8080/protected/home", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.ok) {
        setAuthenticated(true);
      } else {
        setAuthenticated(false);
      }
    } catch (error) {
      console.error("Error verifying authentication:", error);
      setAuthenticated(false);
    }
  } else {
    setAuthenticated(false);
  }
}



export async function fetchNotifications() {
  const token = localStorage.getItem("token");

  if (token) {
    try {
      const response = await fetch("http://localhost:8080/protected/notifications", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.ok) {
        const data = await response.json();
        notifications.set(data.notifications);
      } else {
        console.error("Failed to fetch notifications");
      }
    } catch (error) {
      console.error("Error fetching notifications:", error);
    }
  }
}

export async function markNotificationsAsSeen() {
  const token = localStorage.getItem("token");

  if (token) {
    try {
      const response = await fetch("http://localhost:8080/protected/notifications/mark-seen", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.ok) {
        notifications.update((n) => n.map((notif) => ({ ...notif, seen: true })));
      } else {
        console.error("Failed to mark notifications as seen");
      }
    } catch (error) {
      console.error("Error marking notifications as seen:", error);
    }
  }
}